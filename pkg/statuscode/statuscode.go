package statuscode

import (
	"net/http"

	errmsg "github.com/gocastsian/roham/pkg/err_msg"
)

const (
	IntCodeInvalidParam   = "Invalid request parameter"
	IntCodeNotAuthorize   = "You need to authorize first"
	IntCodeNotPermission  = "You don't have permission"
	IntCodeRecordNotFound = "Record not found"
	IntCodeUnExpected     = "Unexpected issue"
	IntCodeNotFound       = "Not found"
	IntCodeUserExistence  = "User already exists"
	IntCodeValidation     = "User validation error"
)

// MapToHTTPStatusCode maps internal error codes to HTTP status codes
func MapToHTTPStatusCode(err errmsg.ErrorResponse) int {
	switch err.InternalErrCode {
	case IntCodeInvalidParam:
		return http.StatusBadRequest
	case IntCodeNotAuthorize:
		return http.StatusUnauthorized
	case IntCodeNotPermission:
		return http.StatusForbidden
	case IntCodeRecordNotFound:
		return http.StatusNotFound
	case IntCodeValidation:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
