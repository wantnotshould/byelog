// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package redis

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wantnotshould/byelog/conf"
)

type Redis struct {
	client *redis.Client
	prefix string
}

var rdb *Redis

func Init(cfg conf.Redis) {
	rdb = newRedis(cfg)
}

func DB() *Redis {
	return rdb
}

func newRedis(cfg conf.Redis) *Redis {
	opt := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	}

	c := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v\n", err)
	}

	return &Redis{
		client: c,
		prefix: cfg.Prefix,
	}
}

func (r *Redis) key(k string) string {
	if r.prefix == "" {
		return k
	}
	return r.prefix + ":" + k
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	s, err := r.client.Get(ctx, r.key(key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return s, err
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(key), value, ttl).Err()
}

func (r *Redis) SetForever(ctx context.Context, key string, value any) error {
	return r.client.Set(ctx, r.key(key), value, redis.KeepTTL).Err()
}

func (r *Redis) SetNX(ctx context.Context, key string, value any, exp time.Duration) (bool, error) {
	// return r.client.SetNX(ctx, r.key(key), value, exp).Result()
	res, err := r.client.SetArgs(ctx, r.key(key), value, redis.SetArgs{
		Mode: "NX",
		TTL:  exp,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return res == "OK", nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, r.key(key)).Result()
	return n > 0, err
}
