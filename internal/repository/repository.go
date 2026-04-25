// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "github.com/wantnotshould/byelog/internal/ent/db"

type baseRepository struct {
	client *db.Client
}

func newBaseRepository(client *db.Client) *baseRepository {
	return &baseRepository{client: client}
}
