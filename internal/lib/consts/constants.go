package consts

const (
	EmptyString = ""

	RequestIdKey = "request_id"

	LogErr                 = "error"
	LogErrDecodingRequest  = "failed to decode request"
	LogErrEncodingResponse = "failed to encode response"
	LogErrBodyEmpty        = "request body is empty"

	LogErrValidation = "validation failed"

	LogErrFailSaveUrl = "failed to save url"

	LogInfoUrlSaved = "url saved"

	ApiUnknownErr       = "unknown error"
	ApiInternalErr      = "internal error"
	ApiInvalidRequest   = "invalid request"
	ApiUrlNotFound      = "url not found"
	ApiUrlAlreadyExists = "url already exists"
)
