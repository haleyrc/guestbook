package handler

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/haleyrc/guestbook/internal/guest"
	"github.com/haleyrc/guestbook/template"
)

type Guestbook struct {
	logger *slog.Logger
	repo   *guest.Repo
}

func New(logger *slog.Logger, db *sqlx.DB) *Guestbook {
	repo := guest.NewRepo(db)

	guestbook := &Guestbook{
		logger: logger,
		repo:   repo,
	}

	return guestbook
}

func (h *Guestbook) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guests, err := h.repo.FindAll(ctx, 10)
	if err != nil {
		h.logger.Error("failed to load guests", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	props := template.HomePageProps{
		Guests: guests,
	}
	if err := template.HomePage(props).Render(ctx, w); err != nil {
		h.logger.Error("failed to render page", slog.Any("error", err))
		return
	}
}

func (h *Guestbook) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	message := r.FormValue("message")

	guest, err := guest.NewGuest(message)
	if err != nil {
		h.logger.Error("failed to load guests", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.repo.Insert(ctx, guest); err != nil {
		h.logger.Error("failed to load guests", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
