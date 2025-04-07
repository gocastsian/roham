package http

import (
	"context"
	"fmt"
	"github.com/gocastsian/roham/adapter/temporal"
	job "github.com/gocastsian/roham/vectorlayerapp/job/temporal"
	"github.com/gocastsian/roham/vectorlayerapp/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	LayerService service.Service
	Logger       *slog.Logger
	Temporal     temporal.Adapter
}

func NewHandler(LayerService service.Service, logger *slog.Logger, temproal temporal.Adapter) Handler {
	return Handler{
		LayerService: LayerService,
		Logger:       logger,
		Temporal:     temproal,
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

func (h Handler) healthCheckJob(c echo.Context) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	workflowId := "test" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: job.GREETING_QUEUE_NAME,
	}

	we, err := h.Temporal.Client.ExecuteWorkflow(ctx, options, h.LayerService.HealthCheckJob, body.Name)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": fmt.Sprintf("Failed to start workflow: %v", err),
		})
	}

	var res string
	if err := we.Get(ctx, &res); err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": fmt.Sprintf("Failed to get result: %v", err),
		})
	}

	return c.JSON(http.StatusOK, res)
}
