package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	errmsg "roham/pkg/err_msg"
	"roham/pkg/statuscode"
	"roham/pkg/validator"
	"roham/userapp/service/user"
)

type Handler struct {
	UserService user.Service
}

func NewHandler(userSrv user.Service) Handler {
	return Handler{
		UserService: userSrv,
	}
}

func (h Handler) GetAllUsers(c echo.Context) error {

	res, err := h.UserService.GetAllUsers(c.Request().Context())
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.JSON(http.StatusOK, res)
}
