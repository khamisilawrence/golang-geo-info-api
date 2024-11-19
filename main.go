package main

import (
	"geo-info-api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Add Middleware
	e.Use(middleware.Logger())  // Logs HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics and logs them

	// Define Routes
	defineRoutes(e)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

func defineRoutes(e *echo.Echo) {
	// Health check route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Geospatial Country Information API is running")
	})

	// Country info route
	e.GET("/info", handlers.GetCountryInfo)
}
