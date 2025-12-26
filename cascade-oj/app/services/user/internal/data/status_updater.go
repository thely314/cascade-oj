package data

import (
	"context"
	"time"

	"cascade-oj/ent/problemset"

	"github.com/go-kratos/kratos/v2/log"
)

type contestStatusUpdate struct {
	contestID int64
	status    string
}

// flush updates to db when reaches these limits
const flushIntervalSeconds = 30
const maxStatusUpdateQueueSize int = 2

var statusUpdateChan chan *contestStatusUpdate

func InitStatusUpdater(ctx context.Context, data *Data, logger log.Logger) {
	if statusUpdateChan != nil {
		return
	}
	statusUpdateChan = make(chan *contestStatusUpdate, 100)

	logHelper := log.NewHelper(logger)

	go func() {
		// flush when timeout
		ticker := time.NewTicker(flushIntervalSeconds * time.Second)
		defer ticker.Stop()

		updates := map[int64]string{}

		// flush updates
		flush := func() {
			if len(updates) == 0 {
				return
			}
			tx, err := data.db.Tx(ctx)
			if err != nil {
				logHelper.Errorf("failed to start transaction for status updates: %v", err)
			}
			for contestID, status := range updates {
				err := data.db.ProblemSet.Update().
					Where(problemset.IDEQ(contestID)).
					SetStatus(problemset.Status(status)).
					Exec(ctx)
				if err != nil {
					logHelper.Errorf("failed to update contest %d status to %s: %v", contestID, status, err)
					tx.Rollback()
					return
				}
			}
			tx.Commit()
			logHelper.Infof("flushed %d contest status updates", len(updates))
			updates = map[int64]string{}
		}

		for {
			select {
			case <-ctx.Done():
				flush()
				logHelper.Info("status updater stopped")
				return
			case updateInfo := <-statusUpdateChan:
				updates[updateInfo.contestID] = updateInfo.status
				if len(updates) >= maxStatusUpdateQueueSize {
					flush()
				}
			case <-ticker.C:
				flush()
			default:
				// ignore
			}
		}
	}()
}

// enqueue status update
func EnqueueStatusUpdate(contestID int64, status string) {
	if statusUpdateChan == nil {
		return
	}
	statusUpdateChan <- &contestStatusUpdate{
		contestID: contestID,
		status:    status,
	}
}

func CalculateContestStatus(now, start, end time.Time) string {
	if now.Before(start) {
		return problemset.StatusUpcoming.String()
	} else if now.After(end) {
		return problemset.StatusEnded.String()
	}
	return problemset.StatusOngoing.String()
}
