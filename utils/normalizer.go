package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// NormalizeInt แปลงค่าให้เป็น int
func NormalizeInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case string:
		clean := strings.TrimSpace(strings.ReplaceAll(val, ",", ""))
		if clean == "" {
			return 0, nil
		}
		return strconv.Atoi(clean)
	case float64:
		return int(val), nil
	case int:
		return val, nil
	default:
		return 0, errors.New("cannot convert to int")
	}
}

// NormalizeFloat แปลงค่าให้เป็น float64
func NormalizeFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case string:
		clean := strings.TrimSpace(strings.ReplaceAll(val, ",", ""))
		if clean == "" {
			return 0, nil
		}
		return strconv.ParseFloat(clean, 64)
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	default:
		return 0, errors.New("cannot convert to float64")
	}
}

// NormalizeString แปลงค่าให้เป็น string
func NormalizeString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

// Generic Normalizer
func NormalizeValue(v interface{}, targetType string) (interface{}, error) {
	switch targetType {
	case "int":
		return NormalizeInt(v)
	case "float":
		return NormalizeFloat(v)
	case "string":
		return NormalizeString(v), nil
	default:
		return v, nil
	}
}
