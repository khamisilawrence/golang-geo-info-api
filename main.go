package main

import (
	"geo-info-api/handlers"
	"geo-info-api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "geo-info-api/docs" // Import the generated Swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Initialize the cache
	utils.InitCache()

	// Initialize Echo instance
	e := echo.New()

	// Add Middleware
	e.Use(middleware.Logger())  // Logs HTTP requests
	e.Use(middleware.Recover()) // Recovers from panics and logs them

	// Apply rate limiting (20 requests per second per client)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Add Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Define Routes
	defineRoutes(e)

	// Start the server
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func defineRoutes(e *echo.Echo) {
	// Health check route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Geospatial Country Information API is running")
	})

	// Country info route
	e.GET("/info", handlers.GetCountryInfo)
}
