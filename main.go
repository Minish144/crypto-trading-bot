package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Minish144/crypto-trading-bot/di"
	"go.uber.org/zap"
)

func main() {
	ndi, err := di.NewDI()
	if err != nil {
		log.Fatalf("di init error: %v", err)
	}

	z := zap.S().With("context", "main")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	z.Info("starting bot")

	ctx = ndi.Start(ctx)

	<-ctx.Done()

	z.Infow("stopping bot", "reason", ctx.Err())

	ndi.Stop()

	z.Info("bot has stopped")
}
