package handlers

import (
	"geo-info-api/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCountryInfo(t *testing.T) {
	utils.InitCache() // Initialize cache

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/info?country=Germany", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	err := GetCountryInfo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Germany")
}
