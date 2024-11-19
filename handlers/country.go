package handlers

import (
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

	// Fetch country data from Restcountries API
	countryData, err := utils.FetchCountryData(country)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError, err.Error())
	}

	if countryData == nil {
		return utils.HandleError(c, http.StatusNotFound, "Country not found")
	}

	// Build response
	response := map[string]interface{}{
		"country":       countryData["name"].(map[string]interface{})["common"],
		"capital":       countryData["capital"].([]interface{})[0],
		"population":    countryData["population"],
		"area_size_km2": countryData["area"],
		"languages":     countryData["languages"],
		"timezone":      countryData["timezones"],
		"currency":      countryData["currencies"],
		"flag":          countryData["flags"].(map[string]interface{})["png"],
	}

	return utils.HandleSuccess(c, response)

}
