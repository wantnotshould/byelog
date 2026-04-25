// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package queue

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/segmentio/kafka-go"
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/logger"
	"github.com/wantnotshould/byelog/internal/service"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	srv    *service.VisitLogService
}

func NewConsumer(cfg conf.Kafka, srv *service.VisitLogService) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			GroupID:  cfg.GroupID,
			Topic:    cfg.Topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		srv: srv,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	defer func() {
		if err := c.reader.Close(); err != nil {
			logger.Error("failed to close kafka reader", "error", err)
		}
	}()

	logger.Info("kafka Consumer started...")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("kafka Consumer stopped gracefully (context canceled)")
				return
			}

			if errors.Is(err, io.EOF) {
				logger.Info("kafka Consumer connection closed")
				return
			}

			logger.Error("read message error", "error", err)
			return
		}

		var logEntry db.VisitLog
		if err := json.Unmarshal(m.Value, &logEntry); err != nil {
			logger.Error("unmarshal error", "error", err, "payload", string(m.Value))
			continue
		}

		c.srv.Push(&logEntry)
	}
}
