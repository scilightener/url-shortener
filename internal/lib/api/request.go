package api

import (
	"net/http"
	"url-shortener/internal/lib/consts"
)

func GetRequestId(r *http.Request) string {
	if id := r.Context().Value(consts.RequestIdKey); id != nil {
		return id.(string)
	}

	return consts.EmptyString
}
