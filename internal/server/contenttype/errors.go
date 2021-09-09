package contenttype

import "errors"

var (
	ErrContentTypeNotSupported     = errors.New("content type is not supported")
	ErrRespContentTypeNotSupported = errors.New("no response content type supported")
)
