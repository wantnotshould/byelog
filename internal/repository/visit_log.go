// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "github.com/wantnotshould/byelog/internal/ent/db"

type VisitLogRepository struct {
	*baseRepository
}

func NewVisitLogRepository(client *db.Client) *VisitLogRepository {
	return &VisitLogRepository{newBaseRepository(client)}
}
