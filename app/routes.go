package app

import (
	"time"

	"github.com/Teneieiza/go-spinsolf-test/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Stations
	api.Get("/stations/nearby", controllers.GetNearbyStations)
	api.Post("/stations/import/url", controllers.ImportUrlStations)
	api.Post("/stations/import/file", controllers.ImportFileStations)

	// health check
	api.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
					"status": "ok",
					"timestamp": time.Now().Format(time.RFC3339),
			})
	})
}
