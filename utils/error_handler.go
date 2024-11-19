package utils

import (
	"github.com/labstack/echo/v4"
)

// ErrorResponse represents the structure of an API error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// HandleError formats and sends an error response
func HandleError(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorResponse{Message: message})
}
