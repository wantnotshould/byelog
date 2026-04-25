// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package cmd

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wantnotshould/byelog/cmd/flags"
	"github.com/wantnotshould/byelog/conf"
	"github.com/wantnotshould/byelog/internal/bootstrap"
	"github.com/wantnotshould/byelog/internal/router"
	"golang.org/x/net/netutil"
)

func Execute() {
	flag.StringVar(&flags.Data, "data", "data", "Data directory")
	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug mode")
	flag.Parse()

	bootstrap.Run()
	defer bootstrap.Release()

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Interrupt signal (CTRL+C)
		syscall.SIGINT,  // SIGINT signal
		syscall.SIGTERM, // SIGTERM signal (Termination)
		syscall.SIGQUIT, // SIGQUIT signal (Quit)
	)
	defer stop()

	addr := conf.Get().Scheme.Port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on address %s: %s\n", addr, err)
	}
	limitedLn := netutil.LimitListener(ln, 1<<11)

	mux := http.NewServeMux()
	router.Init(mux)

	srv := &http.Server{
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	go func() {
		if err := srv.Serve(limitedLn); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Block until a signal is received
	<-ctx.Done()

	// Gracefully shut down the server with a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown failed, forcing exit: %v\n", err)
	}
}
