package parsers

import (
	"encoding/json"
)

type JSONParser struct{}

func (p *JSONParser) Parse(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
