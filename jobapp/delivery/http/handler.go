package http

import (
	"context"
	"github.com/gocastsian/roham/adapter/temporal"
	"github.com/gocastsian/roham/jobapp/service/job"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	jobService      job.Service
	temporalAdapter temporal.Adapter
}

func NewHandler(jobSvc job.Service, temporalAdp temporal.Adapter) Handler {
	return Handler{
		jobService:      jobSvc,
		temporalAdapter: temporalAdp,
	}
}

func (h Handler) Test(c echo.Context) error {
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
	we, err := h.temporalAdapter.Client.ExecuteWorkflow(ctx, options, h.jobService.Greeting, body.Name)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	var res string
	if err := we.Get(ctx, &res); err != nil {
		log.Fatalf("Failed to get result: %v", err)
	}

	return c.JSON(http.StatusOK, res)
}
