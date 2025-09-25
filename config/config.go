package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigType struct {
	Port            string
	MONGO_URI       string
	DB_NAME         string
	COLLECTION_NAME string
	API_KEY         string
}

//สร้าง LoadConfig เพื่อโหลดค่าต่างๆจาก .env
func LoadConfig() *ConfigType {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	return &ConfigType{
		Port:            getEnv("PORT", "8080"),
		MONGO_URI:       getEnv("MONGO_URI", "-"),
		DB_NAME:         getEnv("DB_NAME", "location"),
		COLLECTION_NAME: getEnv("COLLECTION_NAME", "station"),
		API_KEY:         getEnv("API_KEY", "-"),
	}
}

//fuc ที่ไว้ใช้ดึงค่าจาก env ถ้าไม่มีจะใช้ defaultValue แทน
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}