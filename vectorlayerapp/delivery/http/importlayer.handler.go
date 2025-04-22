package http

import (
	"github.com/gocastsian/roham/pkg/context"
	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/statuscode"
	"github.com/gocastsian/roham/pkg/validator"
	"github.com/gocastsian/roham/types"
	"github.com/gocastsian/roham/vectorlayerapp/service/importlayer"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type ImportLayerHandler struct {
	Service importlayer.Service
	Logger  *slog.Logger
}

func NewImportLayerHandler(LayerService importlayer.Service, logger *slog.Logger) ImportLayerHandler {
	return ImportLayerHandler{
		Service: LayerService,
		Logger:  logger,
	}
}
func (h ImportLayerHandler) healthCheck(c echo.Context) error {
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
func (h ImportLayerHandler) createJob(c echo.Context) error {
	var req importlayer.CreateJobRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorResponse{Message: errmsg.ErrInvalidRequestFormat.Error()})
	}

	userInfo, err := context.ExtractUserInfo(c)
	if err != nil {
		return handleServiceError(c, err)
	}
	req.UserId = types.ID(userInfo.ID)

	response, err := h.Service.CreateJob(c.Request().Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, response)

}
func handleServiceError(c echo.Context, err error) error {
	if vErr, ok := err.(validator.Error); ok {
		return c.JSON(vErr.StatusCode(), vErr)
	}
	if eResp, ok := err.(errmsg.ErrorResponse); ok {
		return c.JSON(statuscode.MapToHTTPStatusCode(eResp), eResp)
	}
	return c.JSON(http.StatusInternalServerError, errmsg.ErrorResponse{
		Message: "Internal server error",
	})
}
