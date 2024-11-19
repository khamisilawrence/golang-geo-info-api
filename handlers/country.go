package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetCountryInfo handles the /info route
func GetCountryInfo(c echo.Context) error {
	// Extract 'country' from the query string
	country := c.QueryParam("country")
	if country == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Country query parameter is required",
		})
	}

	// Placeholder response (will be replaced with API integration)
	response := map[string]interface{}{
		"country": country,
		"message": "Country information will be available soon",
	}

	return c.JSON(http.StatusOK, response)
}
