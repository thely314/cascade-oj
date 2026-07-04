package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "cascade-oj/api/cascade/aiops/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Analyzer orchestrates the alert analysis workflow.
type Analyzer struct {
	promTool  *PrometheusTool
	logTool   *LogReaderTool
	llmClient *LLMClient
	reporter  *Reporter
	log       *log.Helper
}

// NewAnalyzer creates a new Analyzer.
func NewAnalyzer(
	promTool *PrometheusTool,
	logTool *LogReaderTool,
	llmClient *LLMClient,
	reporter *Reporter,
	logger log.Logger,
) *Analyzer {
	return &Analyzer{
		promTool:  promTool,
		logTool:   logTool,
		llmClient: llmClient,
		reporter:  reporter,
		log:       log.NewHelper(logger),
	}
}

// AnalysisResult is the output of a single analysis run.
type AnalysisResult struct {
	TimeRangeMinutes int32
	Services         []string
	RawMetrics       string
	RawLogs          string
	LLMRawResponse   string
	RootCauses       []RootCause
	NoiseAlerts      []AlertItem
	RealRisks        []AlertItem
	SuggestedActions []SuggestedAction
	CreatedAt        time.Time
}

// RootCause represents an identified root cause for an alert cluster.
type RootCause struct {
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
	AffectedSvc string  `json:"affected_service"`
}

// AlertItem represents a single alert (noise or real risk).
type AlertItem struct {
	Source         string  `json:"source"`
	Severity       string  `json:"severity"`
	MetricName     string  `json:"metric_name"`
	MetricValue    float64 `json:"metric_value"`
	Reason         string  `json:"reason"`
	AggregatedFrom int     `json:"aggregated_from,omitempty"`
}

// SuggestedAction represents a remediation suggestion.
type SuggestedAction struct {
	Priority    string `json:"priority"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// Run executes the full analysis pipeline.
func (a *Analyzer) Run(ctx context.Context, timeRangeMinutes int32, services []string) (*AnalysisResult, error) {
	a.log.Infof("Starting analysis: timeRange=%dm services=%v", timeRangeMinutes, services)

	// Step 1: Gather metrics snapshot
	a.log.Info("Step 1: Querying Prometheus metrics...")
	metricsSnapshot, err := a.promTool.GetMetricsSnapshot(ctx, services)
	if err != nil {
		a.log.Errorf("Failed to get metrics: %v", err)
		return nil, fmt.Errorf("metrics query failed: %w", err)
	}

	// Step 2: Gather recent error logs
	a.log.Info("Step 2: Scanning error logs...")
	logSummaries, err := a.logTool.ScanRecent(ctx, timeRangeMinutes, services, 200)
	if err != nil {
		a.log.Errorf("Failed to scan logs: %v", err)
		// Non-fatal: continue with metrics only
		logSummaries = nil
	}

	// Step 3: Format data for LLM with truncation
	metricsJSON, _ := json.MarshalIndent(metricsSnapshot, "", "  ")
	metricsData := string(metricsJSON)
	if len(metricsData) > 4000 {
		metricsData = metricsData[:4000] + "\n... [truncated]"
	}

	logsData := "No error logs found in the specified time range."
	if len(logSummaries) > 0 {
		logsJSON, _ := json.MarshalIndent(logSummaries, "", "  ")
		logsData = string(logsJSON)
		if len(logsData) > 4000 {
			logsData = logsData[:4000] + "\n... [truncated]"
		}
	}

	// Step 4: Call LLM for analysis.
	a.log.Info("Step 3: Calling LLM for analysis...")
	llmOutput, err := a.llmClient.Analyze(ctx, metricsData, logsData)
	if err != nil {
		a.log.Errorf("LLM analysis failed: %v", err)
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Step 5: Assemble result
	result := &AnalysisResult{
		TimeRangeMinutes: timeRangeMinutes,
		Services:         services,
		RawMetrics:       metricsData,
		RawLogs:          logsData,
		LLMRawResponse:   llmOutput.ExecutiveSummary,
		RootCauses:       llmOutput.RootCauses,
		NoiseAlerts:      llmOutput.NoiseAlerts,
		RealRisks:        llmOutput.RealRisks,
		SuggestedActions: llmOutput.SuggestedActions,
		CreatedAt:        time.Now(),
	}

	a.log.Infof("Analysis complete: %d noise, %d real risks, %d root causes",
		len(result.NoiseAlerts), len(result.RealRisks), len(result.RootCauses))

	return result, nil
}

// GetMetricsSnapshot returns a debug view of current Prometheus metrics.
func (a *Analyzer) GetMetricsSnapshot(ctx context.Context, services []string) (*pb.GetMetricsSummaryReply, error) {
	snapshot, err := a.promTool.GetMetricsSnapshot(ctx, services)
	if err != nil {
		return &pb.GetMetricsSummaryReply{
			PrometheusStatus: fmt.Sprintf("error: %v", err),
			QueriedAt:        timestamppb.Now(),
		}, nil
	}

	jsonData, _ := json.MarshalIndent(snapshot, "", "  ")
	return &pb.GetMetricsSummaryReply{
		PrometheusStatus:    "ok",
		MetricsSnapshotJson: string(jsonData),
		QueriedAt:           timestamppb.New(snapshot.Timestamp),
	}, nil
}
