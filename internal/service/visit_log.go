// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package service

import (
	"context"
	"time"

	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/logger"
	"github.com/wantnotshould/byelog/internal/repository"
)

type VisitLogService struct {
	repo      *repository.VisitLogRepository
	logChan   chan *db.VisitLog
	batchSize int
	interval  time.Duration
}

func NewVisitLogService(repo *repository.VisitLogRepository) *VisitLogService {
	s := &VisitLogService{
		repo:      repo,
		logChan:   make(chan *db.VisitLog, 2000),
		batchSize: 500,
		interval:  5 * time.Second,
	}

	go s.startBatchWorker()
	return s
}

func (s *VisitLogService) startBatchWorker() {
	var buffer []*db.VisitLog
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case l, ok := <-s.logChan:
			if !ok {
				s.repo.CreateBulk(context.Background(), buffer)
				return
			}
			buffer = append(buffer, l)
			if len(buffer) >= s.batchSize {
				s.repo.CreateBulk(context.Background(), buffer)
				buffer = buffer[:0]
				ticker.Reset(s.interval)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				s.repo.CreateBulk(context.Background(), buffer)
				buffer = buffer[:0]
			}
		}
	}
}

func (s *VisitLogService) Push(logData *db.VisitLog) {
	select {
	case s.logChan <- logData:
	default:
		logger.Warn("VisitLog channel is full, dropping message")
	}
}
