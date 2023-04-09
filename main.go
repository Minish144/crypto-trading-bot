package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/minish144/crypto-trading-bot/internal/di"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ndi, err := di.NewDI(ctx)
	if err != nil {
		log.Fatalf("di init error: %v", err)
	}

	z := zap.S().With("context", "main")

	z.Info("starting bot")

	ctx = ndi.Start(ctx)

	<-ctx.Done()

	z.Infow("stopping bot", "reason", ctx.Err())

	ndi.Stop(context.Background())

	z.Info("bot has stopped")
}
