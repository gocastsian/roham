package errmsg

import "errors"

var (
	ErrValidationFailed     = errors.New("input validation failed")
	ErrUnexpectedError      = errors.New("unexpected error occurred")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrGetUserInfo          = errors.New("get user info failed")
	ErrFailedDecodeBase64   = errors.New("decode data to base 64 failed")
	ErrFailedUnmarshalJson  = errors.New("unmarshal data to JSON failed")
)

// Define constant messages generally
const (
	MessageMissingXUserData  = "Missing X-User-Data header"
	MessageInvalidBase64     = "Invalid Base64 data"
	MessageInvalidJsonFormat = "Invalid JSON format"
	ServerError              = "Internal server error"
)
