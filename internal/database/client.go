// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package database

import (
	"context"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wantnotshould/byelog/cmd/flags"
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/ent/db"
	"github.com/wantnotshould/byelog/internal/logger"
)

var client *db.Client

func Init(cfg conf.Database) {
	dsn := cfg.DSN()
	drv, err := entsql.Open(dialect.MySQL, dsn)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	dbConn := drv.DB()
	dbConn.SetMaxIdleConns(cfg.MaxIdleConns)
	dbConn.SetMaxOpenConns(cfg.MaxOpenConns)
	dbConn.SetConnMaxLifetime(cfg.MaxLifetime * time.Second)

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	newClient := db.NewClient(db.Driver(drv))

	if flags.Debug {
		newClient = newClient.Debug()
	}

	client = newClient
}

func GetDB() *db.Client {
	if client == nil {
		log.Fatal("database client is not initialized, call Init() first")
	}
	return client
}

func Migrate() {
	ctx := context.Background()
	if err := GetDB().Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func Close() {
	if client != nil {
		if err := client.Close(); err != nil {
			logger.Warn("error failed to close database", err)
		}
	}
}
