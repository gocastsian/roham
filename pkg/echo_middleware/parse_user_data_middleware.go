package echomiddleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	errmsg "roham/pkg/err_msg"
	"roham/types"
)

func ParseUserDataMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// read header
		base64Data := c.Request().Header.Get("X-User-Data")
		if base64Data == "" {
			return c.JSON(http.StatusBadRequest,
				errmsg.ErrorResponse{
					Message: errmsg.ErrGetUserInfo.Error(),
					Errors: map[string]interface{}{
						"header_data_error": errmsg.MessageMissingXUserData,
					},
				},
			)

		}

		// Decode Base64
		jsonData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			// return echo.NewHTTPError(http.StatusBadRequest, "Invalid Base64 data")
			return c.JSON(http.StatusBadRequest,
				errmsg.ErrorResponse{
					Message: errmsg.ErrFailedDecodeBase64.Error(),
					Errors: map[string]interface{}{
						"decode_data_error": errmsg.MessageInvalidBase64,
					},
				},
			)
		}

		// Parse decoded JSON
		var userInfo types.UserInfo
		err = json.Unmarshal(jsonData, &userInfo)
		if err != nil {
			// return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
			return c.JSON(http.StatusBadRequest,
				errmsg.ErrorResponse{
					Message: errmsg.ErrFailedUnmarshalJson.Error(),
					Errors: map[string]interface{}{
						"decode_data_error": errmsg.MessageInvalidJsonFormat,
					},
				},
			)
		}

		// save in context
		c.Set("userInfo", &userInfo)

		// continue request
		return next(c)
	}
}
