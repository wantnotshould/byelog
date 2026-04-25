// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mileusna/useragent"
	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/logger"
	"github.com/wantnotshould/byelog/internal/queue"
	"github.com/wantnotshould/byelog/pkg/utils"
)

type VisitLogHandler struct {
	producer *queue.KafkaProducer
}

func NewVisitLogHandler(p *queue.KafkaProducer) *VisitLogHandler {
	return &VisitLogHandler{p}
}

func (h *VisitLogHandler) Collect(w http.ResponseWriter, r *http.Request) {
	data := h.extractData(r)

	payload, _ := json.Marshal(data)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := h.producer.Publish(ctx, nil, payload); err != nil {
			logger.Error("kafka publish failed", "error", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func (h *VisitLogHandler) extractData(r *http.Request) *db.VisitLog {
	uaRaw := r.UserAgent()
	ua := useragent.Parse(uaRaw)
	ip := utils.IP(r)

	return &db.VisitLog{
		IP:             ip,
		Method:         r.Method,
		Path:           r.URL.Path,
		Query:          r.URL.RawQuery,
		Referer:        r.Referer(),
		Title:          r.URL.Query().Get("title"),
		Os:             ua.OS,
		Browser:        ua.Name,
		BrowserVersion: ua.Version,
		DeviceType:     getDeviceType(ua),
		DeviceModel:    ua.Device,
		Engine:         getEngine(ua.Name),
		IsBot:          ua.Bot,

		CreatedAt: time.Now(),
	}
}

func getDeviceType(ua useragent.UserAgent) string {
	if ua.Mobile {
		return "mobile"
	}
	if ua.Tablet {
		return "tablet"
	}
	if ua.Desktop {
		return "desktop"
	}
	if ua.Bot {
		return "bot"
	}
	return "unknown"
}

func getEngine(browserName string) string {
	switch browserName {
	case "Chrome":
		return "Blink"
	case "Safari":
		return "WebKit"
	case "Firefox":
		return "Gecko"
	case "Internet Explorer":
		return "Trident"
	case "Edge":
		return "EdgeHTML"
	default:
		return "Unknown"
	}
}
