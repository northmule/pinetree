package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/northmule/pinetree/internal/config"
	"github.com/northmule/pinetree/internal/logger"
	"github.com/northmule/pinetree/internal/view"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	var err error

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	l, err := logger.NewLogger(cfg.Value().Log.FlePath, cfg.Value().Log.Level)
	if err != nil {
		return err
	}
	l.Info("Подготовка консольного клиента")
	viewPage := view.NewView(l, cfg)
	return viewPage.InitMain(ctx)
}
