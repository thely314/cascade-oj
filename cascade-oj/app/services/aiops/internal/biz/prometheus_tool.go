package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cascade-oj/app/services/aiops/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// PrometheusTool provides read access to Prometheus metrics via HTTP API.
type PrometheusTool struct {
	addr   string
	client *http.Client
	log    *log.Helper
}

// promQueryResponse is the JSON structure for Prometheus instant query.
type promQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// NewPrometheusTool creates a new PrometheusTool.
func NewPrometheusTool(c *conf.Prometheus, logger log.Logger) *PrometheusTool {
	timeout := 10 * time.Second
	if c.Timeout != nil {
		timeout = c.Timeout.AsDuration()
	}
	return &PrometheusTool{
		addr: c.Addr,
		client: &http.Client{
			Timeout: timeout,
		},
		log: log.NewHelper(logger),
	}
}

// query sends an instant query to the Prometheus HTTP API.
func (p *PrometheusTool) query(ctx context.Context, promQL string) (*promQueryResponse, error) {
	apiURL := fmt.Sprintf("%s/api/v1/query", strings.TrimRight(p.addr, "/"))
	params := url.Values{}
	params.Set("query", promQL)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("prometheus query request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prometheus query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading prometheus response: %w", err)
	}

	var result promQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing prometheus response: %w", err)
	}
	if result.Status != "success" {
		return nil, fmt.Errorf("prometheus query error: %s", result.Error)
	}
	return &result, nil
}

// MetricsSnapshot holds a summary of current service metrics.
type MetricsSnapshot struct {
	Timestamp   time.Time        `json:"timestamp"`
	Services    []ServiceMetrics `json:"services"`
	QueryErrors []string         `json:"query_errors"`
}

// ServiceMetrics holds per-service metrics summary.
type ServiceMetrics struct {
	Service    string  `json:"service"`
	ErrorRate  float64 `json:"error_rate_5xx"`
	P99Latency float64 `json:"p99_latency_seconds"`
	QPS        float64 `json:"qps"`
	Status     string  `json:"status"` // "healthy", "degraded", "down"
}

// defaultServices lists all services monitored by Prometheus.
var defaultServices = []string{"cascade_oj_public", "cascade_oj_user", "cascade_oj_admin", "cascade_oj_judge"}

// GetMetricsSnapshot gathers a snapshot of key metrics from Prometheus.
func (p *PrometheusTool) GetMetricsSnapshot(ctx context.Context, services []string) (*MetricsSnapshot, error) {
	if len(services) == 0 {
		services = defaultServices
	}

	snapshot := &MetricsSnapshot{
		Timestamp: time.Now(),
		Services:  make([]ServiceMetrics, 0, len(services)),
	}

	for _, svc := range services {
		sm := ServiceMetrics{Service: svc, Status: "healthy"}

		// Query 5xx error rate: rate(http_server_requests_seconds_count{job="...",status=~"5.."}[5m])
		errorQL := fmt.Sprintf(
			`sum(rate(http_server_requests_seconds_count{job="cascade-oj-services",instance=~"%s.*",status=~"5.."}[5m])) by (instance)`,
			svc,
		)
		errorResp, err := p.query(ctx, errorQL)
		if err != nil {
			p.log.Warnf("error query for %s: %v", svc, err)
			snapshot.QueryErrors = append(snapshot.QueryErrors, err.Error())
		} else if len(errorResp.Data.Result) > 0 {
			if v, ok := errorResp.Data.Result[0].Value[1].(string); ok {
				fmt.Sscanf(v, "%f", &sm.ErrorRate)
			}
		}

		// Query P99 latency: histogram_quantile(0.99, rate(http_server_requests_seconds_bucket{...}[5m]))
		latencyQL := fmt.Sprintf(
			`histogram_quantile(0.99, sum(rate(http_server_requests_seconds_bucket{job="cascade-oj-services",instance=~"%s.*"}[5m])) by (le, instance))`,
			svc,
		)
		latencyResp, err := p.query(ctx, latencyQL)
		if err != nil {
			p.log.Warnf("latency query for %s: %v", svc, err)
			snapshot.QueryErrors = append(snapshot.QueryErrors, err.Error())
		} else if len(latencyResp.Data.Result) > 0 {
			if v, ok := latencyResp.Data.Result[0].Value[1].(string); ok {
				fmt.Sscanf(v, "%f", &sm.P99Latency)
			}
		}

		// Query QPS: rate(http_server_requests_seconds_count{...}[1m])
		qpsQL := fmt.Sprintf(
			`sum(rate(http_server_requests_seconds_count{job="cascade-oj-services",instance=~"%s.*"}[1m])) by (instance)`,
			svc,
		)
		qpsResp, err := p.query(ctx, qpsQL)
		if err != nil {
			p.log.Warnf("qps query for %s: %v", svc, err)
			snapshot.QueryErrors = append(snapshot.QueryErrors, err.Error())
		} else if len(qpsResp.Data.Result) > 0 {
			if v, ok := qpsResp.Data.Result[0].Value[1].(string); ok {
				fmt.Sscanf(v, "%f", &sm.QPS)
			}
		}

		// Determine status based on thresholds
		if sm.ErrorRate > 1.0 || sm.P99Latency > 5.0 {
			sm.Status = "degraded"
		}
		if sm.ErrorRate > 10.0 {
			sm.Status = "down"
		}

		snapshot.Services = append(snapshot.Services, sm)
	}

	return snapshot, nil
}

// GetActiveAlerts queries Prometheus for currently firing alerts.
func (p *PrometheusTool) GetActiveAlerts(ctx context.Context) ([]AlertItem, error) {
	apiURL := fmt.Sprintf("%s/api/v1/alerts", strings.TrimRight(p.addr, "/"))
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alerts query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Alerts []struct {
				Labels      map[string]string `json:"labels"`
				Annotations map[string]string `json:"annotations"`
				State       string            `json:"state"`
				ActiveAt    string            `json:"activeAt"`
			} `json:"alerts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing alerts response: %w", err)
	}

	var alerts []AlertItem
	for _, a := range result.Data.Alerts {
		if a.State == "firing" {
			alerts = append(alerts, AlertItem{
				Source:   a.Labels["alertname"],
				Severity: a.Labels["severity"],
				Reason:   a.Annotations["description"],
			})
		}
	}
	return alerts, nil
}
