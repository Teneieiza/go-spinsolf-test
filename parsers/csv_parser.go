package parsers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type CSVParser struct{}

//Parse แปลงข้อมูล csv
func (p *CSVParser) Parse(data []byte, target interface{}) error {
    slice := target.(*[]map[string]interface{})

    r := csv.NewReader(bytes.NewReader(data))
    rows, err := r.ReadAll()
    if err != nil {
        return err
    }
    if len(rows) < 1 {
        return errors.New("empty CSV")
    }

    headers := rows[0]

    for _, row := range rows[1:] {
        m := make(map[string]interface{})
        for i, h := range headers {
            if i >= len(row) {
                continue
            }
            m[h] = strings.TrimSpace(row[i])
        }
        *slice = append(*slice, m)
    }

    return nil
}



func setFieldValue(f reflect.Value, v string) {
	v = strings.ReplaceAll(v, ",", "")
	switch f.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		if n, err := strconv.Atoi(v); err == nil {
			f.SetInt(int64(n))
		}
	case reflect.Float32, reflect.Float64:
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			f.SetFloat(n)
		}
	case reflect.String:
		f.SetString(v)
	}
}
