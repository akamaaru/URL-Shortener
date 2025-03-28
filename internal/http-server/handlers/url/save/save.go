package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/akamaaru/url-shortener/internal/lib/api/response"
	"github.com/akamaaru/url-shortener/internal/lib/logger/sl"
	"github.com/akamaaru/url-shortener/internal/lib/random"
	"github.com/akamaaru/url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL		string `json:"url" validate:"required,url"`
	Alias 	string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias 		string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 5

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateError))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		// TODO validate on already existing alias

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		} else if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("url added", slog.String("alias", alias))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: 	response.OK(),
		Alias:		alias,
	})
}