package services

import (
	"context"
	"sort"
	"time"

	"github.com/Teneieiza/go-spinsolf-test/config"
	"github.com/Teneieiza/go-spinsolf-test/dto"
	"github.com/Teneieiza/go-spinsolf-test/models"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetNearbyStations ดึงสถานีที่ใกล้ที่สุด
// รับ lat long page และ limit คืนค่าเป็น slice ของ StationWithDistance ในรูปแบบของ PaginatedResponse(มาจากไฟล์ dto/station_response.go นะจ้ะ)
func GetNearbyStations(lat, long float64, page, limit int) (*dto.PaginatedResponse[dto.StationWithDistance], error) {
	// ตั้ง context set timeout กัน query ค้าง
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//filter เฉพาะสถานีที่ active
	filter := bson.M{
		"active": 1,
	}

	col := config.DB.Collection

	//ค้นหาข้อมูลใน collection โดยใส่ context ไว้ว่าถ้าเวลาเกิน 10วิ ให้ cancel ไป
	cur, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//สร้าง slice มาเก็บข้อมูล ให้อยู่ในรูปแบบของ station_model.go
	//แปลงข้อมูล cur ตาม station_model.go
	var stations []models.Station
	if err := cur.All(ctx, &stations); err != nil {
		return nil, err
	}

	//สร้าง slice มาเก็บข้อมูล ให้อยู่ในรูปแบบของ station_response.go
	results := make([]dto.StationWithDistance, 0, len(stations))
	// loop stations ทั้งหมดโดยแทน s คือ station แต่ละตัว
	// โดยเข้า func คำนวณระยะทาง lat long คือค่าที่รับมา s.lat s.long คือค่าใน station
	for _, s := range stations {
		dist := utils.Haversine(lat, long, s.Lat, s.Long)
		//ถ้าระยะทางที่คำนวณได้น้อยกว่าหรือเท่ากับรัศมีที่กำหนดไว้
		//เพิ่มข้อมูลที่ผ่านการคำนวณระยะทางแล้ว ไปใส่ไว้ใน results ตามรูปแบบของ struct station_response.go
			results = append(results, dto.StationWithDistance{
				ID:          s.ID,
				StationCode: s.StationCode,
				Name:        s.Name,
				EnName:      s.EnName,
				Lat:         s.Lat,
				Long:        s.Long,
				DistanceKM:  dist,
			})
	}

	//เรียงลำดับข้อมูล results จากน้อยไปมาก
	sort.Slice(results, func(i, j int) bool {
		return results[i].DistanceKM < results[j].DistanceKM
	})

	//ตัดข้อมูลตาม page และ limit ที่ส่งมา
	//เช่น limit 10 page 1 ก็จะได้ 0-9, limit 10 page 2 ก็จะได้ 10-19
	start := (page - 1) * limit
	if start > len(results) {
		start = len(results)
	}
	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	//ตัด slice results ตาม start end ที่คำนวณได้
	resultPaginated := results[start:end]

	//คืนค่าเป็น struct PaginatedResponse ที่กำหนดไว้ใน dto/station_response.go
	return &dto.PaginatedResponse[dto.StationWithDistance]{
		Page:     page,
		PageSize: limit,
		Total:    len(results),
		Start:    start + 1,
		End:      end,
		Data:     resultPaginated,
	}, nil
}
