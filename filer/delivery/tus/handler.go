package tus

import (
	"github.com/gocastsian/roham/filer/service/upload"
	"github.com/labstack/echo/v4"
	"net/http"

	//"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

type Handler struct {
	UploadService upload.Service
	TusHandler    *tusd.Handler
}

func NewHandler(srv upload.Service, tusHandler *tusd.Handler) Handler {
	return Handler{
		UploadService: srv,
		TusHandler:    tusHandler,
	}
}

func (h Handler) HandleUpload(c echo.Context) error {
	http.StripPrefix("/uploads/", h.TusHandler).ServeHTTP(c.Response(), c.Request())
	return nil

}
