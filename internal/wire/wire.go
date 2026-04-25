//go:build wireinject
// +build wireinject

// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package wire

import (
	"github.com/google/wire"
	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/handler"
	"github.com/wantnotshould/byelog/internal/repository"
	"github.com/wantnotshould/byelog/internal/service"
)

type App struct {
	VisitLogHandler *handler.VisitLogHandler
}

var providerSet = wire.NewSet(
	handler.NewVisitLogHandler,
	service.NewVisitLogService,
	repository.NewVisitLogRepository,
	wire.Struct(new(App), "*"),
)

func InitApp(client *db.Client) *App {
	wire.Build(providerSet)
	return &App{}
}
