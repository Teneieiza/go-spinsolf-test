package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Teneieiza/go-spinsolf-test/config"
	"github.com/Teneieiza/go-spinsolf-test/parsers"
	"github.com/Teneieiza/go-spinsolf-test/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Import ข้อมูลผ่านไฟล์
func ImportFileStations(filename string, data []byte) (inserted int, updated int, corrupted int, totalImported int, err error) {
	//เลือก parser ตามนามสกุลของไฟล์
	//รองรับ .csv .json .xlsx
	//ถ้าไม่รองรับให้ return error
	ext := strings.ToLower(filepath.Ext(filename))
	var parser parsers.Parser

	switch ext {
	case ".csv":
		parser = &parsers.CSVParser{}
	case ".json":
		parser = &parsers.JSONParser{}
	case ".xlsx":
		parser = &parsers.XLSXParser{}
	default:
		return 0, 0, 0, 0, errors.New("unsupported file format use file with .csv, .json, .xlsx")
	}

	//parse data เข้าไปใน raw ซึ่งเป็น slice ของ map[string]interface{}
	//โดยใช้ parser ที่เลือกมา
	//raw จะมีโครงสร้างคล้ายๆ กับ []models.Station แต่ยังไม่ใช่
	var raw []map[string]interface{}
	if err := parser.Parse(data, &raw); err != nil {
		return 0, 0, 0, 0, err
	}

	// Insert เข้า database mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	col := config.DB.Collection

	//ประกาศ models สำหรับทำ buckwrite
	var models []mongo.WriteModel

	for _, item := range raw {
		st := utils.MapToStation(item)

		// นับ corrupted
		if strings.Contains(st.Comment, "[corrupted]") {
			corrupted++
		}

		// filter station_code เพื่อหาข้อมูลซ้ำ
		filter := bson.M{"station_code": st.StationCode}

		// query document เดิมจาก DB
		var existingDoc map[string]interface{}
		_ = col.FindOne(ctx, filter).Decode(&existingDoc)

		//แปลง Station struct เป็น bson.M
		updateData := utils.StationToBsonMap(st, existingDoc)

		//สร้าง NewUpdateOneModel ไว้ใน model
		//filter สร้างเงื่อนไข
		//field ที่จะให้ update
		//update or insert ถ้ามีข้อมูลตรงตาม filter ก็ update แต่ถ้าไม่มีก็ insert
		if len(updateData) > 0 {
			model := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(bson.M{"$set": updateData}).
				SetUpsert(true) // insert ถ้าไม่มี record
			models = append(models, model)
		}
	}

	//ส่ง models ไปทำ buckwrite ให้ทำ update/insert ทีเดียว
	//ทำ buckwrite เพื่อลด round-trip
	if len(models) > 0 {
		res, err := col.BulkWrite(ctx, models)
		if err != nil {
			return 0, 0, 0, corrupted, err
		}
		inserted = int(res.UpsertedCount)
		updated = int(res.ModifiedCount)
		totalImported = len(models)
	}

	// สร้าง index location 2dsphere หลัง import
	indexModel := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}
	_, _ = col.Indexes().CreateOne(ctx, indexModel)

	return inserted, updated, corrupted, totalImported, nil
}

// Import ข้อมูลผ่าน Url
func ImportUrlStations(apiURL string) (inserted int, updated int, corrupted int, totalImported int, err error) {
	//ส่ง HTTP GET ไปหา URL เพื่อดึง JSON
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer resp.Body.Close()

	// parse JSON ให้เป็น slice ของ map[string]interface{}
	var raw []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, 0, 0, 0, err
	}

	// Insert เข้า database mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	col := config.DB.Collection

	//ประกาศ models สำหรับทำ buckwrite
	var models []mongo.WriteModel

	for _, item := range raw {
		st := utils.MapToStation(item)

		// นับ corrupted
		if strings.Contains(st.Comment, "[corrupted]") {
			corrupted++
		}

		// filter station_code เพื่อหาข้อมูลซ้ำ
		filter := bson.M{"station_code": st.StationCode}

		// query document เดิมจาก DB
		var existingDoc map[string]interface{}
		_ = col.FindOne(ctx, filter).Decode(&existingDoc)

		//แปลง Station struct เป็น bson.M
		updateData := utils.StationToBsonMap(st, existingDoc)

		//สร้าง NewUpdateOneModel ไว้ใน model
		//filter สร้างเงื่อนไข
		//field ที่จะให้ update
		//update or insert ถ้ามีข้อมูลตรงตาม filter ก็ update แต่ถ้าไม่มีก็ insert
		if len(updateData) > 0 {
			model := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(bson.M{"$set": updateData}).
				SetUpsert(true) // insert ถ้าไม่มี record
			models = append(models, model)
		}
	}

	//ส่ง models ไปทำ buckwrite ให้ทำ update/insert ทีเดียว
	//ทำ buckwrite เพื่อลด round-trip
	if len(models) > 0 {
		res, err := col.BulkWrite(ctx, models)
		if err != nil {
			return 0, 0, corrupted, 0, err
		}
		inserted = int(res.UpsertedCount)
		updated = len(models) - inserted
		totalImported = len(models)
	}

	// สร้าง index location 2dsphere หลัง import
	indexModel := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}
	_, _ = col.Indexes().CreateOne(ctx, indexModel)

	return inserted, updated, corrupted, totalImported, nil
}
