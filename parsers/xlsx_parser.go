package parsers

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

type XLSXParser struct{}

func (p *XLSXParser) Parse(data []byte, target interface{}) error {
    f, err := excelize.OpenReader(bytes.NewReader(data))
    if err != nil {
        return err
    }
    defer f.Close()

    sheet := f.GetSheetName(0)
    rows, err := f.GetRows(sheet)
    if err != nil {
        return err
    }
    if len(rows) < 1 {
        return nil
    }

    headers := rows[0]

    // ตรวจสอบว่า target เป็น *[]map[string]interface{}
    sliceVal := reflect.ValueOf(target)
    if sliceVal.Kind() != reflect.Ptr || sliceVal.Elem().Kind() != reflect.Slice {
        return errors.New("target must be a pointer to a slice")
    }
    sliceVal = sliceVal.Elem()

    for _, row := range rows[1:] {
        m := make(map[string]interface{})
        for i, h := range headers {
            if i >= len(row) {
                continue
            }
            m[h] = strings.TrimSpace(row[i])
        }
        sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(m)))
    }
    return nil
}

