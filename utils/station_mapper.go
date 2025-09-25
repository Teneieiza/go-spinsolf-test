package utils

import (
	"github.com/Teneieiza/go-spinsolf-test/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MapToStation แปลง map[string]interface{} เป็น Station model
// โดยใช้ StationFieldTypes ในการ normalize ค่า
func MapToStation(item map[string]interface{}) models.Station {
	st := models.Station{ID: primitive.NewObjectID()}

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
				st.Lat = normVal.(float64)
			case "long":
				st.Long = normVal.(float64)
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
	return st
}
