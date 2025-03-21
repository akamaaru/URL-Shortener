package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akamaaru/url-shortener/internal/lib/api/response"
	"github.com/akamaaru/url-shortener/internal/lib/logger/sl"
	"github.com/akamaaru/url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {}

type Response struct {
	response.Response
	Alias string
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) (error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		log.Info("alias decoded", slog.Any("alias", alias))

		// TODO validate on non-existing alias

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, response.Error("url doesn't exist"))
			return
		} else if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: 	response.OK(),
		Alias:		alias,
	})
}