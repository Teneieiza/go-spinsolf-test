package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type StationWithDistance struct {
	ID          primitive.ObjectID `json:"id"`
	StationCode int                `json:"station_code"`
	Name        string             `json:"name"`
	EnName      string             `json:"en_name"`
	Lat         float64            `json:"lat"`
	Long        float64            `json:"long"`
	DistanceKM  float64            `json:"distance_km"`
}


type ImportStationResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Count   int    `json:"count"`
}