package http

import (
	"github.com/gocastsian/roham/filer/service/filestorage"
	"net/http"
	"time"

	errmsg "github.com/gocastsian/roham/pkg/err_msg"
	"github.com/gocastsian/roham/pkg/statuscode"
	"github.com/gocastsian/roham/pkg/validator"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	storageService filestorage.Service
}

func NewHandler(srv filestorage.Service) Handler {
	return Handler{
		storageService: srv,
	}
}

func (h Handler) DownloadFile(c echo.Context) error {

	key := c.Param("key")

	body, err := h.storageService.GetFileByKey(c.Request().Context(), key)
	if err != nil {
		err := handleError(c, err)
		if err != nil {
			return err
		}
	}

	// todo set headers using res.metadata
	//w.Header().Set("Content-Type", metadata.ContentType)
	//w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileInfo.OriginalName))
	//w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.ContentLength))

	return c.Stream(http.StatusOK, "application/octet-stream", body)
}

func (h Handler) DownloadFileUsingPreSignedURL(c echo.Context) error {
	//todo get pre-signed duration using config

	storageName := "temp-storage"
	url, err := h.storageService.GeneratePreSignedURL(c.Request().Context(), storageName, c.Param("key"), 30*time.Minute)
	if err != nil {
		if vErr, ok := err.(validator.Error); ok {
			return c.JSON(vErr.StatusCode(), vErr)
		}
		return c.JSON(statuscode.MapToHTTPStatusCode(err.(errmsg.ErrorResponse)), err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h Handler) CreateStorage(c echo.Context) error {

	var input filestorage.CreateStorageInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorResponse{Message: err.Error()})
	}
	_, err := h.storageService.CreateStorage(c.Request().Context(), filestorage.CreateStorageInput{
		Name: input.Name,
		Kind: input.Kind,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "everything is good!",
	})
}

func handleError(c echo.Context, err error) error {
	if vErr, ok := err.(validator.Error); ok {
		return c.JSON(vErr.StatusCode(), vErr)
	}
	if eResp, ok := err.(errmsg.ErrorResponse); ok {
		return c.JSON(statuscode.MapToHTTPStatusCode(eResp), eResp)
	}
	return c.JSON(http.StatusInternalServerError, errmsg.ErrorResponse{
		//todo show error if debug is true
		Message: "Internal server error : " + err.Error(),
	})
}
