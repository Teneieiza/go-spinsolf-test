package utils

import (
	"fmt"
	"math"

	"github.com/Teneieiza/go-spinsolf-test/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MapToStation แปลง map[string]interface{} เป็น Station model
// โดยใช้ StationFieldTypes ในการ normalize ค่า
func MapToStation(item map[string]interface{}) models.Station {
	st := models.Station{ID: primitive.NewObjectID()}

	corrupted := false // flag ว่าข้อมูลนี้ถูกแก้

	for field, targetType := range models.FieldTypes {
		if val, ok := item[field]; ok {
			normVal, _ := NormalizeValue(val, targetType)
			switch field {
			case "id":
				st.StationID = normVal.(int)
			case "station_code":
				st.StationCode = normVal.(int)
			case "name":
				st.Name = normVal.(string)
			case "en_name":
				st.EnName = normVal.(string)
			case "th_short":
				st.ThShort = normVal.(string)
			case "en_short":
				st.EnShort = normVal.(string)
			case "chname":
				st.ChName = normVal.(string)
			case "controldivision":
				st.ControlDiv = normVal.(int)
			case "exact_km":
				st.ExactKM = normVal.(int)
			case "exact_distance":
				st.ExactDistance = normVal.(int)
			case "km":
				st.KM = normVal.(int)
			case "class":
				st.Class = normVal.(int)
			case "lat":
				lat := normVal.(float64)
				if lat < -90 || lat > 90 {
					corrupted = true
					lat = 0
				}
				st.Lat = lat
			case "long":
				long := normVal.(float64)
				if long < -180 || long > 180 {
					corrupted = true
					long = 0
				}
				st.Long = long
			case "active":
				st.Active = normVal.(int)
			case "giveway":
				st.Giveway = normVal.(int)
			case "dual_track":
				st.DualTrack = normVal.(int)
			case "comment":
				st.Comment = normVal.(string)
			}
		}
	}

	if corrupted {
		if st.Comment != "" {
			st.Comment += " [corrupted]"
		} else {
			st.Comment = "[corrupted]"
		}
		fmt.Printf("Warning: station_id=%d has corrupted coordinates, set to [0,0]\n", st.StationID)
	}

	// สร้าง field location สำหรับ GeoJSON
	st.Location = map[string]interface{}{
		"type":        "Point",
		"coordinates": []float64{st.Long, st.Lat},
	}

	return st
}


// StationToBsonMap แปลง Station struct เป็น bson.M สำหรับ update
// ไม่รวม _id
func StationToBsonMap(st models.Station, existing map[string]interface{}) map[string]interface{} {
	update := make(map[string]interface{})

	setIfChanged := func(key string, newVal interface{}) {
		oldVal, ok := existing[key]
		if !ok || !isEqual(oldVal, newVal) {
			update[key] = newVal
		}
	}

	setIfChanged("id", st.StationID)
	setIfChanged("station_code", st.StationCode)
	setIfChanged("name", st.Name)
	setIfChanged("en_name", st.EnName)
	setIfChanged("th_short", st.ThShort)
	setIfChanged("en_short", st.EnShort)
	setIfChanged("chname", st.ChName)
	setIfChanged("controldivision", st.ControlDiv)
	setIfChanged("exact_km", st.ExactKM)
	setIfChanged("exact_distance", st.ExactDistance)
	setIfChanged("km", st.KM)
	setIfChanged("class", st.Class)
	setIfChanged("lat", st.Lat)
	setIfChanged("long", st.Long)
	setIfChanged("active", st.Active)
	setIfChanged("giveway", st.Giveway)
	setIfChanged("dual_track", st.DualTrack)
	setIfChanged("comment", st.Comment)
	setIfChanged("location", st.Location)

	return update
}

// ฟังก์ชันช่วยเปรียบเทียบค่าแบบ generic
func isEqual(a, b interface{}) bool {
	switch va := a.(type) {
	case float64:
		vb, ok := b.(float64)
		return ok && math.Abs(va-vb) < 1e-9
	case int, int32, int64:
		return a == b
	case string:
		return a == b
	case nil:
		return b == nil
	default:
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
}



