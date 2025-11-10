package redirect

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

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(slog.String("Operation", op),
			slog.String("Request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Alias is empty")

			render.JSON(w, r, response.Error("Invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)

		if err != nil && errors.Is(err, storage.ErrURLNotFound) {
			log.Info("URL not found", slog.String("alias", alias))

			render.JSON(w, r, response.Error("URL not found"))

			return
		}

		if err != nil {
			log.Error("Failed to get URL", sl.Err(err))

			render.JSON(w, r, response.Error("Internal error"))

			return
		}

		log.Info("Got URL", slog.String("URL", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
