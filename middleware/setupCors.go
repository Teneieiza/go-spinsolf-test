package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupCorsMiddleware(app *fiber.App) {
	//ตั้งต่า cors ให้อนุญาต domain, method, header ที่จะเข้ามาใช้ API
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", //domain
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	//ตั้งค่า logger เพื่อ log request ทีทุกครั้ง
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		// ตัวอย่าง output: [127.0.0.1]:3000 200 - GET /api/users | 2.1ms
	}))

	//ตั้งค่า recover เพื่อป้องกัน server crash เมื่อเกิด panic
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true, // แสดง stack trace เมื่อเกิด panic
	}))
}
