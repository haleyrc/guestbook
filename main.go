package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"github.com/haleyrc/guestbook/internal/app"
)

func main() {
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	a := app.New(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := a.Start(ctx); err != nil {
		logger.Error("failed to start server", slog.Any("error", err))
	}
}
