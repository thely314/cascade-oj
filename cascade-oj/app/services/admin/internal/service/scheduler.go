package service

import (
	"context"
	"time"

	"cascade-oj/app/services/admin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestStatusScheduler struct {
	contestUsecase *biz.ContestUsecase
	log            *log.Helper
	interval       time.Duration
	stopCh         chan struct{}
}

func NewContestStatusScheduler(contestUsecase *biz.ContestUsecase, logger log.Logger, interval time.Duration) *ContestStatusScheduler {
	return &ContestStatusScheduler{
		contestUsecase: contestUsecase,
		log:            log.NewHelper(logger),
		interval:       interval,
		stopCh:         make(chan struct{}),
	}
}

func (s *ContestStatusScheduler) Start() {
	s.log.Info("Contest status scheduler started")
	ticker := time.NewTicker(s.interval)

	go func() {
		defer ticker.Stop()
		s.updateStatuses()
		for {
			select {
			case <-ticker.C:
				s.updateStatuses()
			case <-s.stopCh:
				s.log.Info("Contest status scheduler stopped")
				return
			}
		}
	}()
}

func (s *ContestStatusScheduler) Stop() {
	close(s.stopCh)
}

func (s *ContestStatusScheduler) updateStatuses() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.contestUsecase.UpdateContestStatuses(ctx)
	if err != nil {
		s.log.Errorf("Failed to update contest statuses: %v", err)
	} else {
		s.log.Info("Contest statuses updated successfully")
	}
}
