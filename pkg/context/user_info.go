package context

import (
	"github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/types"
	"github.com/labstack/echo/v4"
	"net/http"
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
