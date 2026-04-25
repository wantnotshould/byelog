// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package service

import "github.com/wantnotshould/byelog/internal/repository"

type VisitLogService struct {
	repo *repository.VisitLogRepository
}

func NewVisitLogService(repo *repository.VisitLogRepository) *VisitLogService {
	return &VisitLogService{repo}
}
