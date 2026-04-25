// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package conf

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wantnotshould/byelog/cmd/flags"
	"github.com/wantnotshould/byelog/pkg/utils"
)

type Scheme struct {
	Port string `json:"port"`
}

type Redis struct {
	Addr         string        `json:"addr"`
	Password     string        `json:"password"`
	DB           int           `json:"db"`
	Prefix       string        `json:"prefix"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	MaxRetries   int           `json:"max_retries"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout"`
}

type Database struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	DBName       string        `json:"db_name"`
	Timezone     string        `json:"timezone"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

func (d *Database) DSN() string {
	tz := d.Timezone
	if tz == "" {
		tz = "Asia/Shanghai"
	}

	params := url.Values{}
	params.Set("charset", "utf8mb4")
	params.Set("parseTime", "True")
	params.Set("loc", tz)

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		params.Encode(),
	)
}

type Logger struct {
	LogsDir    string `json:"logs_dir"`
	MaxSize    int    `json:"max_size"` // MB
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}

type Kafka struct {
	Brokers      []string      `json:"brokers"`
	Topic        string        `json:"topic"`
	GroupID      string        `json:"group_id"`
	MaxAttempts  int           `json:"max_attempts"`
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
	Async        bool          `json:"async"`
}

type Config struct {
	Scheme   Scheme   `json:"scheme"`
	Redis    Redis    `json:"redis"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
	Kafka    Kafka    `json:"kafka"`
}

func (c *Config) validate() error {
	if c.Scheme.Port == "" || !strings.HasPrefix(c.Scheme.Port, ":") {
		return utils.Err("scheme.Port is empty or format error")
	}

	if c.Redis.Addr == "" {
		return utils.Err("redis address can't be empty")
	}

	if c.Redis.Prefix == "" {
		return utils.Err("redis prefix can't be empty")
	}

	if c.Redis.PoolSize <= 0 {
		return utils.Err("redis pool_size must be > 0")
	}

	if c.Redis.MinIdleConns <= 0 {
		return utils.Err("redis min_idle_conns must be > 0")
	}

	if c.Redis.MaxRetries < 0 {
		return utils.Err("redis max_retries must be >= 0")
	}

	if c.Redis.DialTimeout <= 0 {
		return utils.Err("redis dial_timeout must be > 0")
	}

	if c.Redis.ReadTimeout <= 0 {
		return utils.Err("redis read_timeout must be > 0")
	}

	if c.Redis.WriteTimeout <= 0 {
		return utils.Err("redis write_timeout must be > 0")
	}

	if c.Redis.PoolTimeout <= 0 {
		return utils.Err("redis pool_timeout must be > 0")
	}

	if c.Database.Host == "" {
		return utils.Err("database host can't be empty")
	}

	if c.Database.Port == "" {
		return utils.Err("database port can't be empty")
	}

	if c.Database.User == "" {
		return utils.Err("database user can't be empty")
	}

	if c.Database.Password == "" {
		return utils.Err("database password can't be empty")
	}

	if c.Database.DBName == "" {
		return utils.Err("database name can't be empty")
	}

	if c.Logger.LogsDir == "" {
		return utils.Err("logger logs_dir can't be empty")
	}

	if c.Logger.MaxSize <= 0 {
		return utils.Err("logger max_size must be greater than 0")
	}

	if c.Logger.MaxBackups < 0 {
		return utils.Err("logger max_backups can't be negative")
	}

	if c.Logger.MaxAge < 0 {
		return utils.Err("logger max_age can't be negative")
	}

	if len(c.Kafka.Brokers) == 0 {
		return utils.Err("kafka brokers list can't be empty")
	}

	if slices.Contains(c.Kafka.Brokers, "") {
		return utils.Err("kafka broker address can't be empty")
	}

	for _, b := range c.Kafka.Brokers {
		if !strings.Contains(b, ":") {
			return utils.Err("kafka broker address format error, missing port (e.g., 127.0.0.1:9092)")
		}
	}

	if c.Kafka.Topic == "" {
		return utils.Err("kafka topic can't be empty")
	}

	if c.Kafka.GroupID == "" {
		return utils.Err("kafka group_id can't be empty")
	}

	if c.Kafka.MaxAttempts < 0 {
		return utils.Err("kafka max_attempts can't be negative")
	}

	if c.Kafka.BatchSize < 0 {
		return utils.Err("kafka batch_size can't be negative")
	}

	if c.Kafka.BatchTimeout < 0 {
		return utils.Err("kafka batch_timeout can't be negative")
	}

	if c.Kafka.BatchTimeout > 0 && c.Kafka.BatchTimeout < 10*time.Millisecond {
		return utils.Err("kafka batch_timeout is too small, minimum is 10ms")
	}

	return nil
}

var (
	cfgPtr   atomic.Pointer[Config]
	fullPath string
	once     sync.Once
)

func Get() *Config {
	return cfgPtr.Load()
}

func defaultConfig() *Config {
	logsDir := filepath.Join(flags.Data, "logs")

	return &Config{
		Scheme: Scheme{
			Port: ":26039",
		},
		Redis: Redis{
			Addr:         "127.0.0.1:6379",
			Password:     "",
			DB:           0,
			Prefix:       "byelog",
			PoolSize:     100,
			MinIdleConns: 10,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		},
		Database: Database{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Password:     "root",
			DBName:       "byelog",
			Timezone:     "Asia/Shanghai",
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			MaxLifetime:  time.Hour,
		},
		Logger: Logger{
			LogsDir:    logsDir,
			MaxSize:    50,
			MaxBackups: 10,
			MaxAge:     24,
		},
		Kafka: Kafka{
			Brokers:      []string{"127.0.0.1:9092"},
			Topic:        "visit_log_topic",
			GroupID:      "visit_log_group",
			MaxAttempts:  3,
			BatchSize:    100,
			BatchTimeout: time.Second,
			Async:        true,
		},
	}
}

func load() error {
	if fullPath == "" {
		return utils.Err("config path not initialized, call Init() first")
	}

	var newConfig Config
	if err := utils.ReadJSON(fullPath, &newConfig); err != nil {
		return utils.Wrap("failed to load config file", err)
	}

	if err := newConfig.validate(); err != nil {
		return utils.Wrap("config validation failed", err)
	}

	cfgPtr.Store(&newConfig)

	return nil
}

func Init() {
	once.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get working directory: %v\n", err)
		}

		fullPath = filepath.Join(wd, flags.Data, "config.json")

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			def := defaultConfig()
			if err := utils.WriteJSON(fullPath, &def); err != nil {
				log.Fatalf("failed to initialize config file: %v", err)
			}
			cfgPtr.Store(def)
		} else {
			if err := load(); err != nil {
				log.Fatalf("failed to load config file: %v", err)
			}
		}
	})
}
