package app

import (
	"fmt"
	"log"

	"github.com/Teneieiza/go-spinsolf-test/config"
	"github.com/Teneieiza/go-spinsolf-test/middleware"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"github.com/gofiber/fiber/v2"
)

type ApplicationType struct {
	fiber  *fiber.App
	config *config.ConfigType
}

func NewApplication(cfg *config.ConfigType) *ApplicationType {
	//สร้าง fiber app พร้อมตั้งค่า error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			log.Printf("Error: %v", err)
			return utils.ErrorResponse(c, code, message)
		},
	})

	application := &ApplicationType{
		fiber:  app,
		config: cfg,
	}

	// Setup CORS middleware
	middleware.SetupCorsMiddleware(app)

	// Setup API Key middleware สำหรับทุก route /api
	app.Use("/api", middleware.APIKeyMiddleware)

	// Setup API routes
	RegisterRoutes(app)

	return application
}

func (app *ApplicationType) Start() error {
	addr := fmt.Sprintf(":%s", app.config.Port)
	log.Printf("listening on %s", addr)
	return app.fiber.Listen(addr)
}

func (app *ApplicationType) Shutdown() error {
	log.Println("Gracefully shutting down Fiber server...")
	return app.fiber.Shutdown()
}