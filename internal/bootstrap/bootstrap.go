// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package bootstrap

import (
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/cache/redis"
	"github.com/wantnotshould/byelog/internal/database"
	"github.com/wantnotshould/byelog/internal/logger"
)

func Run() {
	conf.Init()
	logger.Init(conf.Get().Logger)
	redis.Init(conf.Get().Redis)
	database.Init(conf.Get().Database)
	database.Migrate()
}

func Release() {
	database.Close()
	redis.DB().Close()
	logger.Close()
}
