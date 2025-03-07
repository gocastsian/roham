package http

import (
	"log/slog"
	"net/http"
	"roham/layerapp/service/layer"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	LayerService layer.Service
	Logger       *slog.Logger
}

func NewHandler(LayerService layer.Service, logger *slog.Logger) Handler {
	return Handler{
		LayerService: LayerService,
		Logger:       logger,
	}
}
func (h Handler) healthCheck(c echo.Context) error {
	check, err := h.LayerService.HealthCheckSrv(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Service in Bad mood ):",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": check,
	})
}
