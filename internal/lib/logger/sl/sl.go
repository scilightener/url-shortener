package sl

import (
	"log/slog"
	"url-shortener/internal/lib/consts"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   consts.LogErr,
		Value: slog.StringValue(err.Error()),
	}
}
