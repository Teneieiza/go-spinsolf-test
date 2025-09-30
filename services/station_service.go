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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetNearbyStations ดึงสถานีใกล้ที่สุด
// รับ lat long และ limit คืนค่าเป็น slice ของ StationWithDistance(มาจากไฟล์ dto/station_response.go นะจ้ะ)
func GetNearbyStations(lat, long float64, limit int) ([]dto.StationWithDistance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := config.DB.Collection

  //filter เฉพาะสถานีที่ active และ ใช้ &near เพื่อดูว่าตัวไหนใกล้เคียงกับ lat long ที่ใส่มา
	filter := bson.M{
		"active": 1,
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type": "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}

	//setlimit ไว้เป็นจำนวณของ limit
	cur, err := col.Find(ctx, filter,
		options.Find().SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var stations []models.Station
	if err := cur.All(ctx, &stations); err != nil {
		return nil, err
	}

	results := make([]dto.StationWithDistance, 0, len(stations))
	for _, s := range stations {
		dist := utils.Haversine(lat, long, s.Lat, s.Long)
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

	return results, nil
}


// GetNearbyStations ดึงสถานีที่ใกล้ที่สุดที่มี pagination
// รับ lat long page และ limit คืนค่าเป็น slice ของ StationWithDistance ในรูปแบบของ PaginatedResponse(มาจากไฟล์ dto/station_response.go นะจ้ะ)
func GetNearbyStationsPage(lat, long float64, page, limit int) (*dto.PaginatedResponse[dto.StationWithDistance], error) {
	// ตั้ง context set timeout กัน query ค้าง
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := config.DB.Collection

	//filter เฉพาะสถานีที่ active และ ใช้ &near เพื่อดูว่าตัวไหนใกล้เคียงกับ lat long ที่ใส่มา
	filter := bson.M{
		"active": 1,
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type": "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}

	//set start ไว้
	start := int64((page - 1) * limit)

	//ค้นหาข้อมูลใน collection โดยใส่ context ไว้ว่าถ้าเวลาเกิน 10วิ ให้ cancel ไป
	//ให้ไปค้าหาข้อมูลที่เริ่มต้นด้วย start และสิ้นสุดที่จำนวณ limit
	cur, err := col.Find(ctx, filter,
		options.Find().SetSkip(start).SetLimit(int64(limit)),
	)
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

	total, err := col.CountDocuments(ctx, bson.M{"active": 1})
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedResponse[dto.StationWithDistance]{
		Page:     page,
		PageSize: limit,
		Total:    int(total),
		Start:    int(start) + 1,
		End:      int(start) + len(results),
		Data:     results,
	}, nil
}
