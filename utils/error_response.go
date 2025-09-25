package utils

import (
	"github.com/Teneieiza/go-spinsolf-test/dto"
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(dto.ErrorResponse{
		Status:  status,
		Error:   httpStatusText(status),
		Message: message,
	})
}

//แปลง status code เป็นข้อความ
func httpStatusText(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "Bad Request"
	case fiber.StatusUnauthorized:
		return "Unauthorized"
	case fiber.StatusForbidden:
		return "Forbidden"
	case fiber.StatusNotFound:
		return "Not Found"
	case fiber.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}
