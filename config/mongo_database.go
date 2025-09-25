package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct เก็บข้อมูล DB
type DatabaseType struct {
	Client     *mongo.Client
	DBName     *mongo.Database
	Collection *mongo.Collection
}

var DB *DatabaseType

// ฟังก์ชัน connect DB
func InitDatabase(ctx context.Context, cfg *ConfigType) error {
	// ตั้ง timeout กันการเชื่อมต่อนานเกินไป
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if cfg.MONGO_URI == "" {
		return fmt.Errorf("MongoDB URL is empty - check your .env file")
	}

	// เชื่อมต่อ MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MONGO_URI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// ping เพื่อเช็คการเชื่อมต่อ
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// เลือก database และ collection
	database := client.Database(cfg.DB_NAME)
	collection := database.Collection(cfg.COLLECTION_NAME)

	// เก็บเข้า DB global
	DB = &DatabaseType{
		Client:     client,
		DBName:     database,
		Collection: collection,
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

// ปิดการเชื่อมต่อ
func (d *DatabaseType) Close(ctx context.Context) error {
	if d.Client != nil {
		log.Println("Closing MongoDB connection...")
		return d.Client.Disconnect(ctx)
	}
	return nil
}
