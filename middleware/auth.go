package middleware

import (
	"os"
	"strings"

	"github.com/Teneieiza/go-spinsolf-test/config"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// APIKeyMiddleware ตรวจสอบ API Key ใน header หรือ query param
func APIKeyMiddleware(c *fiber.Ctx) error {
	//โหลดค่าจาก env
	cfg := config.LoadConfig()

	//โหดค่า API_KEY จาก .env
	apiKeyEnv := cfg.API_KEY

	//ถ้าไม่พบค่าใน env ให้ลองโหลดจากไฟล์ .env ใหม่อีกรอบ
	if apiKeyEnv == "" {
		_ = godotenv.Load(".env")
		apiKeyEnv = os.Getenv("API_KEY")
	}

	//ถ้ายังไม่พบค่า API Key ใน env ให้ส่งกลับ 500 Internal Server Error
	if apiKeyEnv == "" {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "missing server api key config")
	}

	//ตรวจสอบค่า API Key ใน header "x-api-key" หรือ query param "api_key"
	apiKey := c.Get("x-api-key")

	if apiKey == "" {
		apiKey = c.Query("api_key")
	}

	//ถ้าไม่พบค่า API Key ใน header หรือ query param ให้ส่งกลับ 401 Unauthorized
	if apiKey == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "missing api key")
	}

	//ตรวจสอบว่า API Key ที่ส่งมาตรงกับค่าที่ตั้งไว้ใน env หรือไม่ (trim space ด้วย)
	//ถ้าไม่ตรงให้ส่งกลับ 401 Unauthorized ถ้าตรงก็ให้ไปทำตัวอื่นต่อ
	if strings.TrimSpace(apiKey) != apiKeyEnv {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid api key")
	}

	return c.Next()
}
