package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Teneieiza/go-spinsolf-test/dto"
	"github.com/Teneieiza/go-spinsolf-test/services"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"github.com/gofiber/fiber/v2"
)

// Import ข้อมูลผ่านไฟล์
func ImportFileStations(c *fiber.Ctx) error {
	//ดึงไฟล์จาก from-data
	//ถ้าไม่มีไฟล์ให้ส่งกลับ 400 Bad Request
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "file is required")
	}
	//เปิดไฟล์เพื่ออ่านข้อมูล
	f, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	defer f.Close()

	//อ่านข้อมูลทั้งหมดจากไฟล์
	//แล้วส่งข้อมูลไปที่ service ImportFileStations
	//ถ้ามี error ให้ส่งกลับ 500 Internal Server Error
	data, err := io.ReadAll(f)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	inserted, updated, corrupted,totalImported, err := services.ImportFileStations(file.Filename, data)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	message := fmt.Sprintf("Imported %d new, updated %d existing, %d records corrupted", inserted, updated, corrupted)

	//ส่งกลับจำนวนข้อมูลที่ import ได้ พร้อม message ในรูปแบบ JSON
	return c.JSON(dto.ImportStationResponse{
		Status:  200,
		Message: message,
		Count:   inserted + updated,
		TotalImport:  totalImported,
	})
}

// Import ข้อมูลผ่าน URL
func ImportUrlStations(c *fiber.Ctx) error {
	//ดึงค่า url จาก query param
	//ถ้าไม่มี url ให้ส่งกลับ 400 Bad Request
	apiURL := c.Query("url")
	if apiURL == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "url is required")
	}

	//เรียกใช้งาน service ImportUrlStations แล้วส่ง url เข้าไป
	//ถ้ามี error ให้ส่งกลับ 500 Internal Server Error
	inserted, updated, corrupted, totalImported, err := services.ImportUrlStations(apiURL)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	message := fmt.Sprintf("Imported %d new, updated %d existing, %d records corrupted", inserted, updated, corrupted)

	//ส่งกลับจำนวนข้อมูลที่ import ได้ พร้อม message ในรูปแบบ JSON
	return c.JSON(dto.ImportStationResponse{
		Status:  200,
		Message: message,
		Count:   inserted + updated,
		TotalImport:  totalImported,
	})
}
