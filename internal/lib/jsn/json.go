package jsn

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/consts"
	"url-shortener/internal/lib/logger/sl"
)

func EncodeResponse(w http.ResponseWriter, statusCode int, response any, log *slog.Logger) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(consts.LogErrEncodingResponse, sl.Err(err))
		http.Error(w, consts.ApiUnknownErr, http.StatusInternalServerError)
	}
}
