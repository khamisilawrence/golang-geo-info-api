package handlers

import (
	"encoding/json"
	"geo-info-api/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func request(t *testing.T, target string) *httptest.ResponseRecorder {
	t.Helper()
	utils.InitCache()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, GetCountryInfo(c))
	return rec
}

func TestGetCountryInfo(t *testing.T) {
	rec := request(t, "/info?country=Germany")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Germany")
}

func TestGetCountryInfoJamaica(t *testing.T) {
	rec := request(t, "/info?country=jamaica")
	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Data map[string]interface{} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "Jamaica", body.Data["country"])
	assert.Equal(t, "Kingston", body.Data["capital"])
	assert.Equal(t, "https://flagcdn.com/w320/jm.png", body.Data["flag"])
	assert.Greater(t, body.Data["population"], float64(1_000_000))
	assert.Equal(t, float64(10991), body.Data["area_size_km2"])
}

func TestGetCountryInfoMissingParam(t *testing.T) {
	rec := request(t, "/info")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetCountryInfoUnknownCountry(t *testing.T) {
	rec := request(t, "/info?country=Atlantis")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCountryInfoNoCapital(t *testing.T) {
	// Antarctica has no capital; the handler must not panic.
	rec := request(t, "/info?country=Antarctica")
	assert.Equal(t, http.StatusOK, rec.Code)
}
