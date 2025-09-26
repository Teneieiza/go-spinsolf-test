package utils

func GetMultiplier(limit int) int {
    switch {
    case limit <= 50:
        return 5
    case limit <= 150:
        return 3
    default:
        return 2
    }
}

func GetRadiusKM(limit int) float64 {
	switch {
	case limit <= 50:
		return 50
	case limit <= 150:
		return 120
	default:
		return 250
	}
}