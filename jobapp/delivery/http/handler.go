package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Test(c echo.Context) error {
	var name string
	if err := c.Bind(&name); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, name)
}
