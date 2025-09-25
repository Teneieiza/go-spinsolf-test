package controllers

import (
	"net/http"
	"strconv"

	"github.com/Teneieiza/go-spinsolf-test/services"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"github.com/gofiber/fiber/v2"
)

// GetNearbyStations ดึงสถานีที่ใกล้ที่สุดตามพิกัดที่ระบุ
// กำหนด context ของ fiber เพื่อใช้ดึง query param body มาใช้งาน
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

	//ตั้ง limit ไว้สูงสุดที่ 10
	//ถ้าใส่น้อยกว่า 0 หรือ มากกว่า 10 และ ใส่ค่าผิดมาแล้วไม่สามารถแปลงเป็น int ได้ก็ให้แสดง error
	//แต่ถ้าไม่ได้ส่ง limit มาก็ให้ใช้ deafault เป็น 1
	limit, err := strconv.Atoi(c.Query("limit", "1"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid limit")
	}
	if limit <= 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "limit must be greater than 0")
	}
	if limit > 10 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "limit cannot be greater than 10")
	}

	//เรียกใช้งาน service GetNearbyStations แล้วส่ง lat long limit เข้าไป
	stations, err := services.GetNearbyStations(lat, long, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	//แปลงข้อมูลเป็น JSON แล้วส่งกลับไปให้ client ในรูปแบบของ
	//Content-Type: application/json
	return c.JSON(stations)
}

// GetNearbyStationsPaginated ดึงสถานีที่ใกล้ที่สุดแบบแบ่งหน้า
// กำหนด context ของ fiber เพื่อใช้ดึง query param body มาใช้งาน
func GetNearbyStationsPaginated(c *fiber.Ctx) error {
	//รับ lat มาแปลงเป็น float64
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid lat")
	}

	//รับ long มาแปลงเป็น float64
	long, err := strconv.ParseFloat(c.Query("long"), 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid long")
	}

	//รับ page มาแปลงเป็น int
	//แปลงเสร็จแล้ว ถ้า error หรือ page น้อยกว่า 1 ให้ fallback page เป็น 1
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	//รับ page_size มาแปลงเป็น int
	//ถ้าแปลงแล้ว error หรือ page_size น้อยกว่า 10 หรือ มากกว่า 100 ให้แสดง error
	//แต่ถ้าไม่ได้ส่ง page_size มาก็ให้ใช้ default เป็น 10
	pageSizeStr := c.Query("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid page_size")
	}
	if pageSize < 10 || pageSize > 100 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "page_size must be between 10 and 100")
	}

	//เรียกใช้งาน service GetNearbyStationsPaginated แล้วส่ง lat long page page_size เข้าไป
	stations, err := services.GetNearbyStationsPaginated(lat, long, page, pageSize)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	//แปลงข้อมูลเป็น JSON แล้วส่งกลับไปให้ client ในรูปแบบของ
	//Content-Type: application/json
	return c.JSON(fiber.Map{
		"page":      page,
		"page_size": pageSize,
		"data":      stations,
	})
}
