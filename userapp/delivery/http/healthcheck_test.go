package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	server := Server{}
	if assert.NoError(t, server.healthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var body map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &body)
		assert.NoError(t, err)

		expectedMessage := "everything is good!"
		assert.Equal(t, expectedMessage, body["message"])
	}
}
