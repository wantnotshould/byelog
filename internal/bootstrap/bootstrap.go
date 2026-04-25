// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package bootstrap

import (
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/logger"
)

func Run() {
	conf.Init()
	logger.Init(conf.Get().Logger)
}
