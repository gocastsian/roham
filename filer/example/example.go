package main

import (
	"github.com/gocastsian/roham/filer/adapter/tusdadapter"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// Setup Echo
	e := echo.New()

	tusdHandler := tusdadapter.New()

	// Create uploads dir if not exist
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("./uploads", os.ModePerm); err != nil {
			log.Fatalf("Could not create upload dir: %v", err)
		}
	}

	// Route TUS uploads through Echo
	e.Any("/files/*", echo.WrapHandler(http.StripPrefix("/files/", tusdHandler)))

	// Example health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "TUS + Echo server running")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8060"))
}
