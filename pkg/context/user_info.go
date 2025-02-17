package context

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"roham/pkg/err_msg"
	"roham/types"
)

func ExtractUserInfo(c echo.Context) (*types.UserInfo, error) {
	userInfo, ok := c.Get("userInfo").(*types.UserInfo)
	if !ok {
		return nil, c.JSON(http.StatusInternalServerError,
			errmsg.ErrorResponse{
				Message: errmsg.ErrUnexpectedError.Error(),
				Errors: map[string]interface{}{
					"field": "user_info",
				},
			},
		)
	}
	return userInfo, nil
}
