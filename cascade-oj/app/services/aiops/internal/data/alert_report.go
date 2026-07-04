package data

import (
	"context"
	"time"

	"cascade-oj/app/services/aiops/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/alertreport"

	"github.com/go-kratos/kratos/v2/log"
)

// AlertReportRepo handles alert report persistence.
type AlertReportRepo struct {
	data *Data
	log  *log.Helper
}

// NewAlertReportRepo creates a new AlertReportRepo.
func NewAlertReportRepo(data *Data, logger log.Logger) biz.AlertReportRepo {
	return &AlertReportRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create inserts a new alert report.
func (r *AlertReportRepo) Create(ctx context.Context, title, summary, rootCauses, suggestedActions, fullReportMD string, noiseCount, realRiskCount int, status string) (*ent.AlertReport, error) {
	return r.data.DB().AlertReport.Create().
		SetTitle(title).
		SetSummary(summary).
		SetRootCauses(rootCauses).
		SetSuggestedActions(suggestedActions).
		SetNoiseCount(noiseCount).
		SetRealRiskCount(realRiskCount).
		SetStatus(alertreport.Status(status)).
		SetFullReportMarkdown(fullReportMD).
		Save(ctx)
}

// GetByID retrieves a report by its ID.
func (r *AlertReportRepo) GetByID(ctx context.Context, id int64) (*ent.AlertReport, error) {
	return r.data.DB().AlertReport.Get(ctx, id)
}

// List returns paginated reports with optional status filter.
func (r *AlertReportRepo) List(ctx context.Context, offset, limit int, statusFilter string) ([]*ent.AlertReport, int, error) {
	q := r.data.DB().AlertReport.Query()
	if statusFilter != "" {
		q = q.Where(alertreport.StatusEQ(alertreport.Status(statusFilter)))
	}
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	reports, err := q.
		Order(ent.Desc(alertreport.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	return reports, total, err
}

// Acknowledge marks a report as acknowledged by an OCE.
func (r *AlertReportRepo) Acknowledge(ctx context.Context, id int64) error {
	now := time.Now()
	return r.data.DB().AlertReport.UpdateOneID(id).
		SetStatus(alertreport.StatusAcknowledged).
		SetAcknowledgedAt(now).
		Exec(ctx)
}
