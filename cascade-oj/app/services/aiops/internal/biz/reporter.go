package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cascade-oj/ent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReportData holds the structured report data for persistence and API responses.
type ReportData struct {
	ID                   int64
	Title                string
	Summary              string
	RootCausesJSON       string
	SuggestedActionsJSON string
	NoiseCount           int
	RealRiskCount        int
	Status               string
	FullReportMarkdown   string
	CreatedAt            *timestamppb.Timestamp
	AcknowledgedAt       *timestamppb.Timestamp
}

// ReportSummary is a brief view for list endpoints.
type ReportSummaryData struct {
	ID            int64
	Title         string
	NoiseCount    int
	RealRiskCount int
	Status        string
	CreatedAt     *timestamppb.Timestamp
}

// Reporter generates and persists analysis reports.
type Reporter struct {
	alertEventRepo  AlertEventRepo
	alertReportRepo AlertReportRepo
}

type AlertEventRepo interface {
	Create(ctx context.Context, source, severity, metricName string, metricValue, threshold float64, rawData string) (*ent.AlertEvent, error)
	BatchCreate(ctx context.Context, events []*ent.AlertEventCreate) error
	UpdateAggregatedInto(ctx context.Context, alertIDs []int64, reportID int64) (int, error)
	ListUnaggregated(ctx context.Context, limit int) ([]*ent.AlertEvent, error)
}

type AlertReportRepo interface {
	Create(ctx context.Context, title, summary, rootCauses, suggestedActions, fullReportMD string, noiseCount, realRiskCount int, status string) (*ent.AlertReport, error)
	GetByID(ctx context.Context, id int64) (*ent.AlertReport, error)
	List(ctx context.Context, offset, limit int, statusFilter string) ([]*ent.AlertReport, int, error)
	Acknowledge(ctx context.Context, id int64) error
}

// NewReporter creates a new Reporter.
func NewReporter(
	alertEventRepo AlertEventRepo,
	alertReportRepo AlertReportRepo,
) *Reporter {
	return &Reporter{
		alertEventRepo:  alertEventRepo,
		alertReportRepo: alertReportRepo,
	}
}

// GenerateReport creates a new report from an analysis result and persists it.
func (r *Reporter) GenerateReport(ctx context.Context, result *AnalysisResult) (*ReportData, error) {
	rootCausesJSON, _ := json.Marshal(result.RootCauses)
	suggestedActionsJSON, _ := json.Marshal(result.SuggestedActions)

	title := fmt.Sprintf("Alert Analysis - %s (%d risks, %d noise)",
		time.Now().Format("2006-01-02 15:04"),
		len(result.RealRisks),
		len(result.NoiseAlerts),
	)

	// Generate markdown report
	markdown := r.buildMarkdown(title, result)

	// Guard against empty summary (LLM may return blank executive_summary)
	summary := result.LLMRawResponse
	if strings.TrimSpace(summary) == "" {
		summary = fmt.Sprintf("Analysis complete: %d real risks identified across %d services.",
			len(result.RealRisks), len(result.Services))
	}

	// Persist to database
	report, err := r.alertReportRepo.Create(ctx,
		title,
		summary,
		string(rootCausesJSON),
		string(suggestedActionsJSON),
		markdown,
		len(result.NoiseAlerts),
		len(result.RealRisks),
		"published",
	)
	if err != nil {
		return nil, fmt.Errorf("persisting report: %w", err)
	}

	return &ReportData{
		ID:                   report.ID,
		Title:                title,
		Summary:              summary,
		RootCausesJSON:       string(rootCausesJSON),
		SuggestedActionsJSON: string(suggestedActionsJSON),
		NoiseCount:           len(result.NoiseAlerts),
		RealRiskCount:        len(result.RealRisks),
		Status:               "published",
		FullReportMarkdown:   markdown,
		CreatedAt:            timestamppb.New(time.Now()),
	}, nil
}

// GetReport retrieves a single report by ID from the database.
func (r *Reporter) GetReport(ctx context.Context, id int64) (*ReportData, error) {
	report, err := r.alertReportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("report %d not found: %w", id, err)
	}
	return entToReportData(report), nil
}

// ListReports lists reports with pagination.
func (r *Reporter) ListReports(ctx context.Context, page, pageSize int, statusFilter string) ([]ReportSummaryData, int, error) {
	offset := (page - 1) * pageSize
	reports, total, err := r.alertReportRepo.List(ctx, offset, pageSize, statusFilter)
	if err != nil {
		return nil, 0, err
	}

	var summaries []ReportSummaryData
	for _, report := range reports {
		summaries = append(summaries, ReportSummaryData{
			ID:            report.ID,
			Title:         report.Title,
			NoiseCount:    report.NoiseCount,
			RealRiskCount: report.RealRiskCount,
			Status:        string(report.Status),
			CreatedAt:     timestamppb.New(report.CreatedAt),
		})
	}
	return summaries, total, nil
}

// AckReport acknowledges a report by updating its status to "acknowledged".
func (r *Reporter) AckReport(ctx context.Context, id int64) error {
	return r.alertReportRepo.Acknowledge(ctx, id)
}

// entToReportData converts an ent AlertReport to a ReportData.
func entToReportData(r *ent.AlertReport) *ReportData {
	data := &ReportData{
		ID:                   r.ID,
		Title:                r.Title,
		Summary:              r.Summary,
		RootCausesJSON:       r.RootCauses,
		SuggestedActionsJSON: r.SuggestedActions,
		NoiseCount:           r.NoiseCount,
		RealRiskCount:        r.RealRiskCount,
		Status:               string(r.Status),
		FullReportMarkdown:   r.FullReportMarkdown,
		CreatedAt:            timestamppb.New(r.CreatedAt),
	}
	if r.AcknowledgedAt != nil {
		data.AcknowledgedAt = timestamppb.New(*r.AcknowledgedAt)
	}
	return data
}

// buildMarkdown generates a human-readable markdown report.
func (r *Reporter) buildMarkdown(title string, result *AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("**Scan Period**: last %d minutes  \n", result.TimeRangeMinutes))
	sb.WriteString(fmt.Sprintf("**Services Analyzed**: %s  \n", strings.Join(result.Services, ", ")))
	sb.WriteString(fmt.Sprintf("**Generated At**: %s\n\n", result.CreatedAt.Format("2006-01-02 15:04:05 MST")))

	// Statistics
	sb.WriteString("---\n\n")
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Count |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Real Risks | **%d** |\n", len(result.RealRisks)))
	sb.WriteString(fmt.Sprintf("| Noise (filtered) | **%d** |\n", len(result.NoiseAlerts)))
	sb.WriteString(fmt.Sprintf("| Root Causes | **%d** |\n", len(result.RootCauses)))
	sb.WriteString(fmt.Sprintf("| Suggestions | **%d** |\n\n", len(result.SuggestedActions)))

	// Executive summary
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(result.LLMRawResponse)
	sb.WriteString("\n\n")

	// Real risks
	if len(result.RealRisks) > 0 {
		sb.WriteString("## Real Risks\n\n")
		for i, risk := range result.RealRisks {
			severity := strings.ToUpper(risk.Severity)
			sb.WriteString(fmt.Sprintf("### %d. %s (%s)\n\n", i+1, risk.Source, severity))
			sb.WriteString(fmt.Sprintf("- **Source**: %s\n", risk.Source))
			sb.WriteString(fmt.Sprintf("- **Reason**: %s\n", risk.Reason))
			if risk.MetricValue > 0 {
				sb.WriteString(fmt.Sprintf("- **Metric Value**: %.2f\n", risk.MetricValue))
			}
			sb.WriteString("\n")
		}
	}

	// Root causes
	if len(result.RootCauses) > 0 {
		sb.WriteString("## Root Cause Analysis\n\n")
		for i, rc := range result.RootCauses {
			sb.WriteString(fmt.Sprintf("### %d. %s (Confidence: %.0f%%)\n\n", i+1, rc.AffectedSvc, rc.Confidence*100))
			sb.WriteString(fmt.Sprintf("%s\n\n", rc.Description))
		}
	}

	// Suggested actions
	if len(result.SuggestedActions) > 0 {
		sb.WriteString("## Suggested Actions\n\n")
		sb.WriteString("| Priority | Action | Description |\n")
		sb.WriteString("|----------|--------|-------------|\n")
		for _, sa := range result.SuggestedActions {
			sb.WriteString(fmt.Sprintf("| **%s** | %s | %s |\n", strings.ToUpper(sa.Priority), sa.Action, sa.Description))
		}
		sb.WriteString("\n")
	}

	// Noise alerts (collapsed)
	if len(result.NoiseAlerts) > 0 {
		sb.WriteString("## Filtered Noise\n\n")
		sb.WriteString(
			fmt.Sprintf("<details>\n<summary>Click to expand - %d alerts classified as noise</summary>\n\n", len(result.NoiseAlerts)),
		)
		sb.WriteString("| Source | Reason |\n")
		sb.WriteString("|--------|--------|\n")
		for _, na := range result.NoiseAlerts {
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", na.Source, na.Reason))
		}
		sb.WriteString("\n</details>\n\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("*Report generated by Cascade OJ AIOps agent.*\n")

	return sb.String()
}
