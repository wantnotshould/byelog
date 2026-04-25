// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package bootstrap

import (
	"context"
	"time"

	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/cache/redis"
	"github.com/wantnotshould/byelog/internal/database"
	"github.com/wantnotshould/byelog/internal/logger"
	"github.com/wantnotshould/byelog/internal/wire"
)

var (
	consumerCancel context.CancelFunc
)

func Run() {
	conf.Init()
	logger.Init(conf.Get().Logger)
	redis.Init(conf.Get().Redis)
	database.Init(conf.Get().Database)
	database.Migrate()

	client := database.GetDB()
	app := wire.InitApp(client)

	var ctx context.Context
	ctx, consumerCancel = context.WithCancel(context.Background())

	go func() {
		logger.Info("consumer starting...")
		app.VisitLogConsumer.Start(ctx)
	}()
}

func Release() {
	if consumerCancel != nil {
		logger.Info("stopping kafka consumer...")
		consumerCancel()
	}

	time.Sleep(2 * time.Second)

	database.Close()
	redis.DB().Close()
	logger.Close()
}
