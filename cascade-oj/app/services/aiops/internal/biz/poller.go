package biz

import (
	"context"
	"time"

	"cascade-oj/app/services/aiops/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// Poller runs periodic health checks and triggers analysis when alert thresholds are exceeded.
type Poller struct {
	enabled   bool
	interval  time.Duration
	threshold int32
	analyzer  *Analyzer
	reporter  *Reporter
	promTool  *PrometheusTool
	log       *log.Helper
	stopCh    chan struct{}
}

// NewPoller creates a new Poller. Returns nil if polling is disabled in config.
func NewPoller(
	c *conf.Polling,
	analyzer *Analyzer,
	reporter *Reporter,
	promTool *PrometheusTool,
	logger log.Logger,
) *Poller {
	if !c.Enabled {
		return nil
	}
	interval := 5 * time.Minute
	if c.Interval != nil {
		interval = c.Interval.AsDuration()
	}
	threshold := int32(0)
	if c.AlertThreshold > 0 {
		threshold = c.AlertThreshold
	}
	return &Poller{
		enabled:   c.Enabled,
		interval:  interval,
		threshold: threshold,
		analyzer:  analyzer,
		reporter:  reporter,
		promTool:  promTool,
		log:       log.NewHelper(logger),
		stopCh:    make(chan struct{}),
	}
}

// Start begins the polling loop in a background goroutine.
func (p *Poller) Start() {
	if p == nil {
		return
	}
	p.log.Infof("AIOps poller started: interval=%v threshold=%d", p.interval, p.threshold)
	go p.loop()
}

// Stop signals the polling loop to stop.
func (p *Poller) Stop() {
	if p == nil {
		return
	}
	close(p.stopCh)
}

// loop is the main polling loop.
func (p *Poller) loop() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Debounce: track last analysis time to avoid rapid successive triggers
	var lastAnalysis time.Time
	minInterval := 2 * time.Minute

	for {
		select {
		case <-p.stopCh:
			p.log.Info("AIOps poller stopped")
			return
		case <-ticker.C:
			p.checkAndTrigger(&lastAnalysis, minInterval)
		}
	}
}

// checkAndTrigger polls Prometheus for active alerts and triggers analysis if threshold exceeded.
func (p *Poller) checkAndTrigger(lastAnalysis *time.Time, minInterval time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check for active alerts
	alerts, err := p.promTool.GetActiveAlerts(ctx)
	if err != nil {
		p.log.Warnf("Poller: failed to query alerts: %v", err)
		return
	}

	// Also check if any service is in degraded/down state
	snapshot, err := p.promTool.GetMetricsSnapshot(ctx, nil)
	if err != nil {
		p.log.Warnf("Poller: failed to query metrics: %v", err)
		return
	}

	degradedCount := 0
	for _, svc := range snapshot.Services {
		if svc.Status != "healthy" {
			degradedCount++
		}
	}

	totalAlertSignals := len(alerts) + degradedCount
	p.log.Debugf("Poller: %d alerts, %d degraded services (threshold: %d)", len(alerts), degradedCount, p.threshold)

	if int32(totalAlertSignals) < p.threshold {
		return
	}

	// Debounce check
	if time.Since(*lastAnalysis) < minInterval {
		p.log.Infof("Poller: debounced - last analysis was %v ago", time.Since(*lastAnalysis).Round(time.Second))
		return
	}

	// Trigger analysis
	p.log.Infof("Poller: threshold exceeded (%d >= %d), triggering analysis...", totalAlertSignals, p.threshold)
	*lastAnalysis = time.Now()

	// Detach from HTTP request context — the full pipeline (Prometheus →
	// log scan → LLM CoT → DB persist) is a long-running task that must
	// not be killed by the HTTP server timeout.
	// 避免 ctx 被 HTTP 请求超时取消，整个分析流程是一个长时间运行的任务
	bgCtx := context.WithoutCancel(ctx)

	result, err := p.analyzer.Run(bgCtx, 15, nil)
	if err != nil {
		p.log.Errorf("Poller: analysis failed: %v", err)
		return
	}

	_, err = p.reporter.GenerateReport(bgCtx, result)
	if err != nil {
		p.log.Errorf("Poller: report generation failed: %v", err)
		return
	}

	p.log.Infof("Poller: auto-analysis complete - %d risks, %d noise", len(result.RealRisks), len(result.NoiseAlerts))
}
