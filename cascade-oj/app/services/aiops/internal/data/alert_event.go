package data

import (
	"context"
	"time"

	"cascade-oj/app/services/aiops/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/alertevent"

	"github.com/go-kratos/kratos/v2/log"
)

// AlertEventRepo handles alert event persistence.
type AlertEventRepo struct {
	data *Data
	log  *log.Helper
}

// NewAlertEventRepo creates a new AlertEventRepo.
func NewAlertEventRepo(data *Data, logger log.Logger) biz.AlertEventRepo {
	return &AlertEventRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Create inserts a new alert event record.
func (r *AlertEventRepo) Create(ctx context.Context, source, severity, metricName string, metricValue, threshold float64, rawData string) (*ent.AlertEvent, error) {
	return r.data.DB().AlertEvent.Create().
		SetSource(source).
		SetSeverity(alertevent.Severity(severity)).
		SetMetricName(metricName).
		SetMetricValue(metricValue).
		SetThreshold(threshold).
		SetRawData(rawData).
		SetReceivedAt(time.Now()).
		Save(ctx)
}

// BatchCreate inserts multiple alert events in one transaction.
func (r *AlertEventRepo) BatchCreate(ctx context.Context, events []*ent.AlertEventCreate) error {
	bulk := make([]*ent.AlertEventCreate, len(events))
	copy(bulk, events)
	_, err := r.data.DB().AlertEvent.CreateBulk(bulk...).Save(ctx)
	return err
}

// UpdateAggregatedInto marks alerts as aggregated into a report.
func (r *AlertEventRepo) UpdateAggregatedInto(ctx context.Context, alertIDs []int64, reportID int64) (int, error) {
	return r.data.DB().AlertEvent.Update().
		Where(alertevent.IDIn(alertIDs...)).
		SetAggregatedInto(reportID).
		Save(ctx)
}

// ListUnaggregated returns alerts not yet aggregated into any report.
func (r *AlertEventRepo) ListUnaggregated(ctx context.Context, limit int) ([]*ent.AlertEvent, error) {
	return r.data.DB().AlertEvent.Query().
		Where(alertevent.AggregatedInto(0)).
		Order(ent.Desc(alertevent.FieldReceivedAt)).
		Limit(limit).
		All(ctx)
}
