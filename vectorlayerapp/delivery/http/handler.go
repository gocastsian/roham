package http

import (
	"github.com/gocastsian/roham/vectorlayerapp/service"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type Handler struct {
	LayerService service.Service
	Logger       *slog.Logger
}

func NewHandler(LayerService service.Service, logger *slog.Logger) Handler {
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

func (h Handler) ImportLayer(c echo.Context) error {
	res, err := h.LayerService.ScheduleImportLayer(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":    "success",
		"workflowId": res.WorkflowId,
	})

}
