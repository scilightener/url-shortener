package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/consts"
	"url-shortener/internal/lib/jsn"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UrlGetter
type UrlGetter interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	const eo = "handlers.rest.redirect.New"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("eo", eo),
			slog.String(consts.RequestIdKey, api.GetRequestId(r)),
		)

		alias := r.PathValue(consts.AliasKey)
		if len(alias) == 0 {
			log.Info(consts.LogInfoAliasEmpty)
			jsn.EncodeResponse(w, http.StatusBadRequest, &api.Response{Error: consts.ApiInvalidRequest}, log)
			return
		}

		url, err := urlGetter.GetUrl(r.Context(), alias)
		if errors.Is(err, storage.ResourceNotFound) {
			log.Info(consts.LogInfoUrlNotFound, consts.AliasKey, alias)
			jsn.EncodeResponse(w, http.StatusNotFound, &api.Response{Error: consts.ApiUrlNotFound}, log)
			return
		}
		if err != nil {
			log.Error(consts.LogErrFailGetUrl, sl.Err(err))
			jsn.EncodeResponse(w, http.StatusInternalServerError, &api.Response{Error: consts.ApiInternalErr}, log)
			return
		}

		log.Info(consts.LogInfoGotUrl, slog.String(consts.UrlKey, url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
