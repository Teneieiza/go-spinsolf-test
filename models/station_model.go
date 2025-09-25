package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Station struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StationID     int                `bson:"id" json:"station_id"`
	StationCode   int                `bson:"station_code" json:"station_code"`
	Name          string             `bson:"name" json:"name"`
	EnName        string             `bson:"en_name" json:"en_name"`
	ThShort       string             `bson:"th_short" json:"th_short"`
	EnShort       string             `bson:"en_short" json:"en_short"`
	ChName        string             `bson:"chname" json:"chname"`
	ControlDiv    int                `bson:"controldivision" json:"controldivision"`
	ExactKM       int                `bson:"exact_km" json:"exact_km"`
	ExactDistance int                `bson:"exact_distance" json:"exact_distance"`
	KM            int                `bson:"km" json:"km"`
	Class         int                `bson:"class" json:"class"`
	Lat           float64            `bson:"lat" json:"lat"`
	Long          float64            `bson:"long" json:"long"`
	Active        int                `bson:"active" json:"active"`
	Giveway       int                `bson:"giveway" json:"giveway"`
	DualTrack     int                `bson:"dual_track" json:"dual_track"`
	Comment       string             `bson:"comment" json:"comment"`
}


var FieldTypes = map[string]string{
	"id":              "int",
	"station_code":    "int",
	"name":            "string",
	"en_name":         "string",
	"th_short":        "string",
	"en_short":        "string",
	"chname":          "string",
	"controldivision": "int",
	"exact_km":        "int",
	"exact_distance":  "int",
	"km":              "int",
	"class":           "int",
	"lat":             "float",
	"long":            "float",
	"active":          "int",
	"giveway":         "int",
	"dual_track":      "int",
	"comment":         "string",
}