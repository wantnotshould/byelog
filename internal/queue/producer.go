// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package queue

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/wantnotshould/byelog/conf"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(cfg conf.Kafka) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    cfg.BatchSize,
			MaxAttempts:  cfg.MaxAttempts,
			BatchTimeout: 1 * time.Second,
			Async:        cfg.Async,
		},
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, key []byte, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
