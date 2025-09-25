package parsers

type Parser interface {
	// data: raw file data
	// target: pointer ของ slice เช่น *[]models.Station
	Parse(data []byte, target interface{}) error
}
