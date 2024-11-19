package utils

import "github.com/labstack/echo/v4"

// SuccessResponse represents the structure of a successful API response
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// HandleSuccess formats and sends a successful response
func HandleSuccess(c echo.Context, data interface{}) error {
	return c.JSON(200, SuccessResponse{Data: data})
}
