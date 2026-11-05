package handlers

import (
	"errors"
	"geo-info-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCountryInfo godoc
// @Summary Retrieve country information
// @Description Get detailed geospatial and demographic information about a country.
// @Tags country
// @Accept  json
// @Produce  json
// @Param   country query string true "Country name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /info [get]
func GetCountryInfo(c echo.Context) error {
	// Extract 'country' from query parameter
	country := c.QueryParam("country")
	if country == "" {
		return utils.HandleError(c, http.StatusBadRequest, "Country query parameter is required")
	}

	// Look the country up in the embedded dataset
	countryData, err := utils.FetchCountryData(country)
	if errors.Is(err, utils.ErrCountryNotFound) {
		return utils.HandleError(c, http.StatusNotFound, "Country not found")
	}
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Build response
	response := map[string]interface{}{
		"country":       nestedString(countryData, "name", "common"),
		"capital":       firstString(countryData["capital"]),
		"population":    countryData["population"],
		"area_size_km2": countryData["area"],
		"languages":     countryData["languages"],
		"timezone":      countryData["timezones"],
		"currency":      countryData["currencies"],
		"flag":          nestedString(countryData, "flags", "png"),
	}

	return utils.HandleSuccess(c, response)
}

// nestedString reads data[outer][inner] as a string, or "" if it is missing.
func nestedString(data map[string]interface{}, outer, inner string) string {
	m, ok := data[outer].(map[string]interface{})
	if !ok {
		return ""
	}
	s, _ := m[inner].(string)
	return s
}

// firstString returns the first element of a string slice, or "" if empty.
// Some territories (for example Antarctica) have no capital.
func firstString(v interface{}) string {
	switch list := v.(type) {
	case []string:
		if len(list) > 0 {
			return list[0]
		}
	case []interface{}:
		if len(list) > 0 {
			s, _ := list[0].(string)
			return s
		}
	}
	return ""
}
