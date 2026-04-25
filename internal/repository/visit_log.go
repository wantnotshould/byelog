// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import (
	"context"

	"github.com/wantnotshould/byelog/internal/ent/db"
)

type VisitLogRepository struct {
	*baseRepository
}

func NewVisitLogRepository(client *db.Client) *VisitLogRepository {
	return &VisitLogRepository{newBaseRepository(client)}
}

func (r *VisitLogRepository) CreateBulk(ctx context.Context, logs []*db.VisitLog) error {
	if len(logs) == 0 {
		return nil
	}

	builders := make([]*db.VisitLogCreate, len(logs))
	for i, l := range logs {
		builders[i] = r.client.VisitLog.Create().
			SetAppID(l.AppID).
			SetIP(l.IP).
			SetMethod(l.Method).
			SetPath(l.Path).
			SetQuery(l.Query).
			SetTitle(l.Title).
			SetReferer(l.Referer).
			SetOs(l.Os).
			SetBrowser(l.Browser).
			SetBrowserVersion(l.BrowserVersion).
			SetDeviceType(l.DeviceType).
			SetDeviceModel(l.DeviceModel).
			SetEngine(l.Engine).
			SetIsBot(l.IsBot).
			SetCreatedAt(l.CreatedAt)
	}

	return r.client.VisitLog.CreateBulk(builders...).Exec(ctx)
}
