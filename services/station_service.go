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
// รับ lat long และ limit คืนค่าเป็น slice ของ StationWithDistance(มาจากไฟล์ dto/station_response.go นะจ้ะ)
func GetNearbyStations(lat, long float64, limit int) ([]dto.StationWithDistance, error) {
	// ตั้ง context set timeout กัน query ค้าง
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//collection ที่จะใช้ query
	col := config.DB.Collection

	//ค้นหาข้อมูลใน collection โดยใส่ context ไว้ว่าถ้าเวลาเกิน 10วิ ให้ cancel ไป
	//และหา field ที่ "active": 1 ข้อมูลที่ออกมาก็จะอยู่ในรูปแบบของ map[string]interface{}
	cur, err := col.Find(ctx, bson.M{"active": 1})
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
	var results []dto.StationWithDistance
	// loop stations ทั้งหมดโดยแทน s คือ station แต่ละตัว
	// โดยเข้า func คำนวณระยะทาง lat long คือค่าที่รับมา s.lat s.long คือค่าใน station
	for _, s := range stations {
		dist := utils.Haversine(lat, long, s.Lat, s.Long)
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

	//ถ้าข้อมูลที่ออกมามีมากว่า limit ก็ให้ตัด slice เหลือจำนวณตาม limit
	if len(results) > limit {
		results = results[:limit]
	}

	//return results
	return results, nil
}


// GetNearbyStationsPaginated ดึงสถานีที่ใกล้ที่สุดแบบมี pagination
// รับ lat long page pageSize คืนค่าเป็น slice ของ StationWithDistance(มาจากไฟล์ dto/station_response.go นะจ้ะ)
func GetNearbyStationsPaginated(lat, long float64, page, pageSize int) ([]dto.StationWithDistance, error) {
	//ดึงข้อมูลสถานีที่ใกล้ที่สุดมาเก็บไว้ใน all
	//โดยกำหนด limit 1000 เพื่อให้แน่ใจว่ามีข้อมูลเพียงพอสำหรับการแบ่งหน้า
	all, err := GetNearbyStations(lat, long, 1000)
	if err != nil {
		return nil, err
	}

	//คำนวณตำแหน่งเริ่มต้นและสิ้นสุดของข้อมูลที่จะแสดงในแต่ละหน้า
	//ถ้า start มากกว่าข้อมูลจริง ก็ให้ return arrayว่าง
	start := (page - 1) * pageSize
	if start > len(all) {
		return []dto.StationWithDistance{}, nil
	}

	//ถ้า end มากกว่าข้อมูลจริง ก็ให้ end = ความยาวของข้อมูลจริง
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}

	//return เป็น slice ย่อยเฉพาะหน้านั้นๆ
	return all[start:end], nil
}