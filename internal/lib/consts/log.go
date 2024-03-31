package consts

const (
	LogErr                 = "error"
	LogErrDecodingRequest  = "failed to decode request"
	LogErrEncodingResponse = "failed to encode response"
	LogErrBodyEmpty        = "request body is empty"

	LogErrValidation = "validation failed"

	LogErrFailSaveUrl = "failed to save url"
	LogErrFailGetUrl  = "failed to get url"

	LogInfoAliasEmpty  = "alias is empty"
	LogInfoUrlSaved    = "url saved"
	LogInfoUrlNotFound = "url not found"
	LogInfoGotUrl      = "got url"
)
