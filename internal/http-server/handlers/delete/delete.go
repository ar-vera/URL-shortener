package delete

import (
	"URL-shortener/internal/lib/api/response"
	"URL-shortener/internal/lib/logger/sl"
	"URL-shortener/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name URLDeleter --dir ../../../../.. --output ./mocks --filename mock_url_deleter.go --with-expecter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Alias is empty")

			render.JSON(w, r, response.Error("Invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)

		if err != nil && errors.Is(err, storage.ErrURLNotFound) {
			log.Info("URL not found", slog.String("alias", alias))

			render.JSON(w, r, response.Error("URL not found"))

			return
		}

		if err != nil {
			log.Error("Failed to delete URL", sl.Err(err))

			render.JSON(w, r, response.Error("Failed to delete URL"))

			return
		}

		log.Info("URL deleted", slog.String("alias", alias))

		render.JSON(w, r, response.OK())
	}
}
