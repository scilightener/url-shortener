package save

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/consts"
	"url-shortener/internal/lib/jsn"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	Alias string `jsn:"alias,omitempty"`
	Url   string `jsn:"url" validate:"required,url"`
}

type Response struct {
	api.Response
	Alias string `jsn:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UrlSaver
type UrlSaver interface {
	SaveUrl(alias, url string, validUntil int64) (int64, error)
}

const (
	aliasKey = "alias"
	idKey    = "id"
)

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	const eo = "handlers.url.save.New"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("eo", eo),
			slog.String(consts.RequestIdKey, r.Context().Value(consts.RequestIdKey).(string)),
		)

		request := new(Request)

		if err := json.NewDecoder(r.Body).Decode(request); errors.Is(err, io.EOF) {
			log.Error(consts.LogErrBodyEmpty)
			jsn.EncodeResponse(w, http.StatusBadRequest,
				&Response{api.ErrResponse(consts.ApiInvalidRequest), consts.EmptyString}, log)
			return
		} else if err != nil {
			log.Error(consts.LogErrDecodingRequest, sl.Err(err))
			jsn.EncodeResponse(w, http.StatusInternalServerError,
				&Response{api.ErrResponse(consts.ApiInternalErr), consts.EmptyString}, log)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			var validErrs validator.ValidationErrors
			errors.As(err, &validErrs)
			log.Error(consts.LogErrValidation, sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest,
				&Response{api.ValidationError(validErrs), consts.EmptyString}, log)
			return
		}

		alias := request.Alias
		if len(alias) == 0 {
			alias = uuid.New().String()
		}

		id, err := urlSaver.SaveUrl(request.Alias, request.Url, 0)
		if err != nil {
			log.Error(consts.LogErrFailSaveUrl, sl.Err(err))
			if errors.Is(err, storage.ResourceAlreadyExists) {
				jsn.EncodeResponse(w, http.StatusConflict,
					&Response{api.ErrResponse(consts.ApiUrlAlreadyExists), consts.EmptyString}, log)
			} else {
				jsn.EncodeResponse(w, http.StatusInternalServerError,
					&Response{api.ErrResponse(consts.ApiInternalErr), consts.EmptyString}, log)
			}
			return
		}

		log.Info(consts.LogInfoUrlSaved, slog.String(aliasKey, alias), slog.Int64(idKey, id))

		jsn.EncodeResponse(w, http.StatusCreated, &Response{api.OkResponse(), alias}, log)
	}
}
