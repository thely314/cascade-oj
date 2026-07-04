package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cascade-oj/app/services/aiops/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// LLMClient wraps an OpenAI-compatible API for alert analysis.
type LLMClient struct {
	apiKey     string
	baseURL    string
	model      string
	client     *http.Client
	maxRetries int
	log        *log.Helper
}

// chatMessage represents a single message in a chat completion request.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest is the OpenAI-compatible chat completion request body.
type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
	Temperature    float64         `json:"temperature,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
}

// responseFormat enforces JSON output mode.
type responseFormat struct {
	Type string `json:"type"`
}

// chatResponse is the OpenAI-compatible chat completion response.
type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// NewLLMClient creates a new LLMClient.
func NewLLMClient(c *conf.LLM, logger log.Logger) *LLMClient {
	timeout := 60 * time.Second
	if c.Timeout != nil {
		timeout = c.Timeout.AsDuration()
	}
	maxRetries := int(c.MaxRetries)
	if maxRetries <= 0 {
		maxRetries = 2
	}
	return &LLMClient{
		apiKey:     c.ApiKey,
		baseURL:    c.BaseUrl,
		model:      c.Model,
		client:     &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
		log:        log.NewHelper(logger),
	}
}

// systemPrompt is the SRE expert system prompt for alert analysis.
const systemPrompt = `You are an expert SRE (Site Reliability Engineer) with deep experience in alert analysis and noise reduction.
Your task is to analyze service metrics and error logs, then produce a structured analysis that:

1. Identifies real risks that require immediate attention (with root cause analysis).
2. Aggregates duplicate or cascading alerts caused by the same root cause.
3. Filters out low-priority or informational noise.
4. Provides actionable remediation suggestions ordered by priority.

Classification rules:
- "CRITICAL": Service is down, error rate > 10%, or P99 latency > 10s — needs immediate action.
- "WARNING": Error rate > 1%, P99 latency > 1s, or trending degradation — schedule investigation.
- "INFO": Normal metrics but anomaly detected — watch only.

For each real risk, include:
- Root cause description with confidence level (0.0-1.0)
- Affected services
- Suggested actions with priority (urgent/high/medium/low)

Output strictly as valid JSON matching the specified schema.

Example output in JSON:
{
  "title": "Alert Analysis - 2024-06-01 12:00",
  "executive_summary": "Analysis complete: 1 real risks identified across 5 services.",
  "root_causes": [
	{
	  "description": "Service A is down",
	  "confidence": 0.95,
	  "affected_service": "Service A, B"
  	},
	{
	  "description": "Service B cannot reach upstream Service A",
	  "confidence": 0.9,
	  "affected_service": "Service B"
  	}
  ],
  "noise_alerts": [
	{
	  "source": "Service C",
	  "severity": "INFO",
	  "metric_name": "error_rate",
	  "metric_value": 0.5,
	  "reason": "Low error rate, no action needed"
  	}
  ],
  "real_risks": [
	{
	  "source": "Service B",
	  "severity": "CRITICAL",
	  "metric_name": "error_rate",
	  "metric_value": 12.5,
	  "reason": "High error rate, requires immediate attention"
  	}
  ],
  "suggested_actions": [
  	{
	  "priority": "urgent",
	  "action": "Restart Service A",
	  "description": "Service A is down, restart it to restore functionality."
  	}
  ]
}

`

// AnalysisOutput is the structured output expected from the LLM.
type AnalysisOutput struct {
	Title            string            `json:"title"`
	ExecutiveSummary string            `json:"executive_summary"`
	RootCauses       []RootCause       `json:"root_causes"`
	NoiseAlerts      []AlertItem       `json:"noise_alerts"`
	RealRisks        []AlertItem       `json:"real_risks"`
	SuggestedActions []SuggestedAction `json:"suggested_actions"`
}

// Analyze sends metrics and log data to the LLM for analysis.
func (c *LLMClient) Analyze(ctx context.Context, metricsData, logsData string) (*AnalysisOutput, error) {
	userPrompt := fmt.Sprintf(`Analyze the following service metrics and error logs. Perform alert noise reduction.

## Current Metrics Snapshot
%s

## Recent Error Logs
%s

Return a JSON object with the analysis results.`, metricsData, logsData)

	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := chatRequest{
		Model:          c.model,
		Messages:       messages,
		ResponseFormat: &responseFormat{Type: "json_object"},
		Temperature:    0.1,
		MaxTokens:      2000,
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			c.log.Warnf("LLM retry attempt %d/%d", attempt, c.maxRetries)
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}

		result, err := c.sendRequest(ctx, req)
		if err != nil {
			lastErr = err
			continue
		}
		return result, nil
	}
	return nil, fmt.Errorf("LLM request failed after %d retries: %w", c.maxRetries+1, lastErr)
}

// sendRequest sends a chat completion request and parses the response.
func (c *LLMClient) sendRequest(ctx context.Context, req chatRequest) (*AnalysisOutput, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/chat/completions", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if chatResp.Error != nil {
		return nil, fmt.Errorf("LLM API error: %s", chatResp.Error.Message)
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned no choices")
	}

	content := chatResp.Choices[0].Message.Content
	var output AnalysisOutput
	if err := json.Unmarshal([]byte(content), &output); err != nil {
		// If JSON parsing fails, return a basic structure with raw content
		log.Warnf("failed to unmarshal LLM answer, which is: '%s'", content)
		return &AnalysisOutput{
			Title:            "Alert Analysis",
			ExecutiveSummary: content,
			RootCauses:       []RootCause{},
			NoiseAlerts:      []AlertItem{},
			RealRisks:        []AlertItem{},
			SuggestedActions: []SuggestedAction{},
		}, nil
	}

	return &output, nil
}
