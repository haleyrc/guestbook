package guest

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Guest struct {
	ID        int64  `db:"id"`
	Message   string `db:"message"`
	CreatedAt string `db:"created_at"`
}

func NewGuest(message string) (Guest, error) {
	guest := Guest{
		Message:   message,
		CreatedAt: time.Now().UTC().Format(time.DateTime),
	}
	return guest, nil
}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	repo := &Repo{
		db: db,
	}

	return repo
}

func (r *Repo) FindAll(ctx context.Context, count int) ([]Guest, error) {
	guests := []Guest{}
	query := `SELECT id, message, created_at FROM guest ORDER BY created_at DESC LIMIT ?`
	if err := r.db.SelectContext(ctx, &guests, query, count); err != nil {
		return nil, fmt.Errorf("guest: repo: find all: failed to query database: %w", err)
	}
	return guests, nil
}

func (r *Repo) Insert(ctx context.Context, guest Guest) error {
	query := `INSERT INTO guest (message, created_at) VALUES (?, ?)`
	if _, err := r.db.ExecContext(ctx, query, guest.Message, guest.CreatedAt); err != nil {
		return fmt.Errorf("guest: repo: insert: failed to insert guest: %w", err)
	}
	return nil
}
