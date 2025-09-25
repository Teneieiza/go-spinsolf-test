package utils

import "math"

// ฟังก์ชันคำนวณระยะทาง Haversine
func Haversine(lat1, long1, lat2, long2 float64) float64 {
	// รัศมีของโลกในหน่วยกิโลเมตร
	const R = 6371.0
	// ฟังก์ชันแปลงองศาเป็นเรเดียน
	radian := func(degree float64) float64 { return degree * math.Pi / 180.0 }

	//แปลง lat long เป็นเรเดียน
	radianLat1 := radian(lat1)
	radianLat2 := radian(lat2)

	//คำนวณความต่างของ lat long
	diffRLat := radian(lat2 - lat1)
	diffRLong := radian(long2 - long1)

	a := 	math.Sin(diffRLat/2)*math.Sin(diffRLat/2) +
				math.Cos(radianLat1)*math.Cos(radianLat2)*math.Sin(diffRLong/2)*math.Sin(diffRLong/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}