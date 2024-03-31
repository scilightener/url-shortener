package save

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"time"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/bl"
	"url-shortener/internal/lib/consts"
	"url-shortener/internal/lib/jsn"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
	Url   string `json:"url" validate:"required,url"`
}

type Response struct {
	api.Response
	Alias         string    `json:"alias,omitempty"`
	ValidUntilUTC time.Time `json:"valid_until_utc,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UrlRepo
type UrlRepo interface {
	SaveUrl(alias, url string, validUntil time.Time) (int64, error)
	GetUrl(alias string) (string, error)
}

const (
	aliasKey = "alias"
	idKey    = "id"
)

func New(log *slog.Logger, urlRepo UrlRepo) http.HandlerFunc {
	const eo = "handlers.url.save.New"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("eo", eo),
			slog.String(consts.RequestIdKey, api.GetRequestId(r)),
		)

		request := new(Request)

		if err := json.NewDecoder(r.Body).Decode(request); errors.Is(err, io.EOF) {
			log.Error(consts.LogErrBodyEmpty)
			jsn.EncodeResponse(w, http.StatusBadRequest,
				&Response{Response: api.ErrResponse(consts.ApiInvalidRequest)}, log)
			return
		} else if err != nil {
			log.Error(consts.LogErrDecodingRequest, sl.Err(err))
			jsn.EncodeResponse(w, http.StatusInternalServerError,
				&Response{Response: api.ErrResponse(consts.ApiInternalErr)}, log)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			var validErrs validator.ValidationErrors
			errors.As(err, &validErrs)
			log.Error(consts.LogErrValidation, sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest,
				&Response{Response: api.ValidationError(validErrs)}, log)
			return
		}

		alias := request.Alias
		if len(alias) == 0 {
			alias = bl.GenerateUniqueAlias(urlRepo)
		}

		validUntilUTC := bl.GetValidUntilUTC()

		id, err := urlRepo.SaveUrl(request.Alias, request.Url, validUntilUTC)
		if err != nil {
			log.Error(consts.LogErrFailSaveUrl, sl.Err(err))
			if errors.Is(err, storage.ResourceAlreadyExists) {
				jsn.EncodeResponse(w, http.StatusConflict,
					&Response{Response: api.ErrResponse(consts.ApiUrlAlreadyExists)}, log)
			} else {
				jsn.EncodeResponse(w, http.StatusInternalServerError,
					&Response{Response: api.ErrResponse(consts.ApiInternalErr)}, log)
			}
			return
		}

		log.Info(consts.LogInfoUrlSaved, slog.String(aliasKey, alias), slog.Int64(idKey, id))

		jsn.EncodeResponse(w, http.StatusCreated,
			&Response{api.OkResponse(), alias, validUntilUTC}, log)
	}
}
