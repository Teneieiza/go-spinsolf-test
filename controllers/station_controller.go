package controllers

import (
	"net/http"
	"strconv"

	"github.com/Teneieiza/go-spinsolf-test/services"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"github.com/gofiber/fiber/v2"
)

func GetNearbyStations(c *fiber.Ctx) error {
	// แปลง lat เป็น float64
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid lat")
	}

	// แปลง long เป็น float64
	long, err := strconv.ParseFloat(c.Query("long"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid long")
	}

	// limit (default = 10)
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// เรียกใช้งาน service GetNearbyStations แล้วส่ง lat, long, limit เข้าไป
	stations, err := services.GetNearbyStations(lat, long, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// ส่งกลับข้อมูล JSON
	return c.JSON(stations)
}



// GetNearbyStations ดึงสถานีที่ใกล้ที่สุดตามพิกัดที่ระบุ
// กำหนด context ของ fiber เพื่อใช้ดึง query param body มาใช้งาน
func GetNearbyStationsPage(c *fiber.Ctx) error {
	// แปลง lat เป็น float64
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid lat")
	}

	// แปลง long เป็น float64
	long, err := strconv.ParseFloat(c.Query("long"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid long")
	}

	// แปลง page, limit เป็น int
	// กำหนด fallback page=1, limit=10 เวลาได้รับข้อมูลผิด
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

		limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

//เรียกใช้งาน service GetNearbyStations แล้วส่ง lat long page limit เข้าไป
	stations, err := services.GetNearbyStationsPage(lat, long, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// ส่งกลับข้อมูลในรูปแบบ JSON
	return c.JSON(stations)
}
