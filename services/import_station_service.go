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
)

// Import ข้อมูลผ่านไฟล์
func ImportFileStations(filename string, data []byte) (int, error) {
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
		return 0, errors.New("unsupported file format use file with .csv, .json, .xlsx")
	}

	//parse data เข้าไปใน raw ซึ่งเป็น slice ของ map[string]interface{}
	//โดยใช้ parser ที่เลือกมา
	//raw จะมีโครงสร้างคล้ายๆ กับ []models.Station แต่ยังไม่ใช่
	var raw []map[string]interface{}
	if err := parser.Parse(data, &raw); err != nil {
		return 0, err
	}

	//แปลงข้อมูลด้วย normalizer ให้อยู่ในรูปแบบของ models.Station
	//เก็บเข้า slice stations
	var stations []interface{}
	for _, item := range raw {
		st := utils.MapToStation(item)
		stations = append(stations, st)
	}

	// Insert เข้า database mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := config.DB.Collection
	_, err := col.InsertMany(ctx, stations)
	if err != nil {
		return 0, err
	}

	return len(stations), nil
}

// Import ข้อมูลผ่าน Url
func ImportUrlStations(apiURL string) (int, error) {
	//ส่ง HTTP GET ไปหา URL เพื่อดึง JSON
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// parse JSON ให้เป็น slice ของ map[string]interface{}
	var raw []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, err
	}

	//แปลงข้อมูลด้วย normalizer ให้อยู่ในรูปแบบของ models.Station
	//เก็บเข้า slice stations
	var stations []interface{}
	for _, item := range raw {
		st := utils.MapToStation(item)
		stations = append(stations, st)
	}

	// Insert เข้า database mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := config.DB.Collection
	_, err = col.InsertMany(ctx, stations)
	if err != nil {
		return 0, err
	}

	return len(stations), nil
}
