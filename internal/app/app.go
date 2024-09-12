package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/haleyrc/guestbook/internal/handler"
)

type App struct {
	db     *sqlx.DB
	logger *slog.Logger
	router *http.ServeMux
}

func New(logger *slog.Logger) *App {
	router := http.NewServeMux()

	app := &App{
		logger: logger,
		router: router,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	db, err := sqlx.Open("sqlite3", "demo.db")
	if err != nil {
		return fmt.Errorf("app: start: failed to connect to database: %w", err)
	}

	a.db = db

	a.loadRoutes()

	server := http.Server{
		Addr:    ":8086",
		Handler: a.router,
	}

	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("failed to listen and serve", slog.Any("error", err))
		}
		close(done)
	}()

	a.logger.Info("Server listening", slog.String("addr", ":8086"))
	select {
	case <-done:
		break
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		server.Shutdown(ctx)
		cancel()
	}

	return nil
}

func (a *App) loadRoutes() {
	guestbook := handler.New(a.logger, a.db)
	a.router.HandleFunc("GET /{$}", guestbook.Home)
	a.router.HandleFunc("POST /{$}", guestbook.Create)
}
