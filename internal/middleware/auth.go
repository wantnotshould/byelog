// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/wantnotshould/byelog/internal/cache/redis"
	"github.com/wantnotshould/byelog/internal/logger"
)

const HeaderAppID = "Bye-Log-App-Id"

func CheckAppID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		appIDStr := r.Header.Get(HeaderAppID)

		if appIDStr == "" {
			http.Error(w, "missing header", http.StatusBadRequest)
			return
		}

		uid, err := uuid.Parse(appIDStr)
		if err != nil {
			http.Error(w, "invalid uuid format", http.StatusBadRequest)
			return
		}

		exist, err := redis.DB().Exists(r.Context(), uid.String())
		if err != nil {
			logger.Error("failed to check redis exists", err)
			http.Error(w, "invalid or unauthorized AppID", http.StatusUnauthorized)
			return
		}

		if !exist {
			http.Error(w, "invalid or unauthorized AppID", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
