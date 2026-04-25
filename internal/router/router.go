// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package router

import (
	"net/http"

	"github.com/wantnotshould/byelog/internal/database"
	"github.com/wantnotshould/byelog/internal/middleware"
	"github.com/wantnotshould/byelog/internal/wire"
)

func Init(mux *http.ServeMux) {
	client := database.GetDB()
	app := wire.InitApp(client)
	mux.HandleFunc("POST /api/v1/collect", middleware.CheckAppID(app.VisitLogHandler.Collect))
}
