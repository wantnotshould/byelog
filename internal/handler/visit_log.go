// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"fmt"
	"net/http"

	"github.com/wantnotshould/byelog/internal/service"
)

type VisitLogHandler struct {
	srv *service.VisitLogService
}

func NewVisitLogHandler(srv *service.VisitLogService) *VisitLogHandler {
	return &VisitLogHandler{srv}
}

func (h *VisitLogHandler) Collect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
