//go:build wireinject
// +build wireinject

// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package wire

import (
	"github.com/google/wire"
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/handler"
	"github.com/wantnotshould/byelog/internal/queue"
	"github.com/wantnotshould/byelog/internal/repository"
	"github.com/wantnotshould/byelog/internal/service"
)

func provideKafkaProducer() *queue.KafkaProducer {
	cfg := conf.Get().Kafka
	return queue.NewProducer(cfg)
}

func provideKafkaConsumer(srv *service.VisitLogService) *queue.KafkaConsumer {
	cfg := conf.Get().Kafka
	return queue.NewConsumer(cfg, srv)
}

type App struct {
	VisitLogHandler  *handler.VisitLogHandler
	VisitLogConsumer *queue.KafkaConsumer
}

var providerSet = wire.NewSet(
	provideKafkaProducer,
	provideKafkaConsumer,
	handler.NewVisitLogHandler,
	service.NewVisitLogService,
	repository.NewVisitLogRepository,
	wire.Struct(new(App), "*"),
)

func InitApp(client *db.Client) *App {
	wire.Build(providerSet)
	return &App{}
}
