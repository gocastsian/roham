package http

import (
	"github.com/gocastsian/roham/vectorlayerapp/service/generatetile"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type GenLayerHandler struct {
	Service generatetile.Service
	Logger  *slog.Logger
}

func NewGenLayerHandler(LayerService generatetile.Service, logger *slog.Logger) GenLayerHandler {
	return GenLayerHandler{
		Service: LayerService,
		Logger:  logger,
	}
}
func (h GenLayerHandler) healthCheck(c echo.Context) error {
	check, err := h.Service.HealthCheckSrv(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Service in Bad mood ):",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": check,
	})
}
