package table

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Supple map data insert to struct values
func (t *TableStruct) Supple(rows map[string]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	for k, v := range t.elem {
		if d, ok := rows[k]; ok {
			if err := SetValue(v, d); err != nil {
				return err
			}
		}
	}
	return nil
}

// SuppleWithTag Supple map data insert to data values by reflect tag name
func (t *TableStruct) SuppleWithTag(rows map[string]interface{}, tag string) error {
	return SuppleWithTag(t.data, rows, tag)
}

// SuppleWithTag : supple rows insert to data struct
// data is ptr
func SuppleWithTag(data interface{}, rows map[string]interface{}, tag string) error {
	if len(rows) == 0 {
		return nil
	} else if tag == "" {
		return errors.New("Supple tag name can not be null")
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := reflect.TypeOf(data).Elem()
	if rv.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var fileds = make(map[string]reflect.Value)
	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Tag.Get(tag)
		if key != "" {
			fileds[key] = rv.Field(i)
		}
	}

	for k, v := range rows {
		if d, ok := fileds[k]; ok {
			if err := SetValue(d, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetValue reflect.Value insert insert with data interface{}
// Type: string, int, uint, float, bool, time.Time
func SetValue(value reflect.Value, data interface{}) error {
	switch value.Interface().(type) {

	case string:
		value.SetString(fmt.Sprint(data))

	case int, int8, int16, int32, int64:
		if num, ok := TypeIntAndUint(data); ok {
			value.SetInt(num)
		} else {
			return fmt.Errorf("Parse %v to int64 error", data)
		}

	case uint, uint8, uint16, uint32, uint64:
		if num, ok := TypeIntAndUint(data); ok {
			value.SetUint(uint64(num))
		} else {
			return fmt.Errorf("Parse %v to uint64 error", data)
		}

	case float32, float64:
		if num, ok := TypeFloatInt(data); ok {
			value.SetFloat(num)
		} else {
			return fmt.Errorf("Parse %v to float64 error", data)
		}

	case bool:
		if state, ok := TypeBool(data); ok {
			value.SetBool(state)
		} else {
			return fmt.Errorf("Parse %v to bool error", data)
		}

	case time.Time:
		if stamp, ok := TypeTime(data); ok {
			value.Set(reflect.ValueOf(stamp))
		} else {
			return fmt.Errorf("Parse %v to time error", data)
		}

	default:
		return fmt.Errorf("No support %s type", value.Kind())
	}
	return nil
}

// TypeBool check data type and convert to return
func TypeBool(data interface{}) (bool, bool) {
	switch data.(type) {

	case bool:
		return data.(bool), true
	case string:
		str := strings.ToLower(data.(string))
		if str == "true" || str == "0" {
			return true, true
		}
		return false, true
	}

	if num, ok := TypeIntAndUint(data); ok {
		if num == 0 {
			return true, true
		}
		return false, true
	}
	return false, false
}

// TypeIntAndUint check data type and convert to return
func TypeIntAndUint(data interface{}) (int64, bool) {
	switch data.(type) {
	case int:
		return int64(data.(int)), true
	case int8:
		return int64(data.(int8)), true
	case int16:
		return int64(data.(int16)), true
	case int32:
		return int64(data.(int32)), true
	case int64:
		return int64(data.(int64)), true
	case uint:
		return int64(data.(uint)), true
	case uint8:
		return int64(data.(uint8)), true
	case uint16:
		return int64(data.(uint16)), true
	case uint32:
		return int64(data.(uint32)), true
	case uint64:
		return int64(data.(uint64)), true
	case string:
		num, err := strconv.ParseInt(data.(string), 10, 64)
		if err == nil {
			return num, true
		}
	}
	return 0, false
}

// TypeFloatInt check data type and convert to return
func TypeFloatInt(data interface{}) (float64, bool) {
	switch data.(type) {
	case float32:
		return float64(data.(float32)), true
	case float64:
		return data.(float64), true
	case string:
		num, err := strconv.ParseFloat(data.(string), 64)
		if err == nil {
			return num, true
		}
	}

	if a, ok := TypeIntAndUint(data); ok {
		return float64(a), true
	}

	return 0, false
}

// TypeTime check data type and convert to return
func TypeTime(data interface{}) (time.Time, bool) {
	switch data.(type) {

	case time.Time:
		return data.(time.Time), true
	case string:
		if stamp, err := ParseTime(data.(string)); err == nil {
			return stamp, true
		}
	case int64:
		if data.(int64)/1e12 > 0 {
			return time.Unix(0, data.(int64)), true
		}
		return time.Unix(data.(int64), 0), true
	}
	return time.Now(), false
}

// ParseTime string to time.Time
func ParseTime(str string) (time.Time, error) {
	layout := TimeLayout(str)
	if len(layout) == 0 {
		return time.Now(), fmt.Errorf("Not time format")
	}

	var err error
	for _, format := range layout {
		if stamp, err := time.ParseInLocation(format, str, time.Local); err == nil {
			return stamp, nil
		}
	}

	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		if num/1e12 > 0 {
			return time.Unix(0, num), nil
		}
		return time.Unix(num, 0), nil
	}

	return time.Now(), err
}

// TimeLayout check time string lenght  to time format
func TimeLayout(str string) []string {
	switch len(str) {
	case 4:
		return []string{"2006"}
	case 6:
		return []string{"200601"}
	case 8:
		return []string{"20060102"}
	case 10:
		return []string{"2006-01-02", "2006010215"}
	case 12:
		return []string{"200601021504"}
	case 13:
		return []string{"2006-01-02 15"}
	case 14:
		return []string{"200601021504"}
	case 16:
		return []string{"2006-01-02 15:04", "20060102150405"}
	case 19:
		return []string{"2006-01-02 15:04:05"}
	default:
		return nil
	}
}

// SuppleWithMap :
// support type: [string, int, uint, float, bool, time.Time]
// bool type: ture{"ture", "0"}
// name: is struct field tag name defalut: 'ref'
// replace: [true:replace, false:not_replace]
//	replase defalut false: if origin value not nil, not replace in
func SuppleWithMap(data interface{}, rows map[string]string, replace bool, name string) error {
	if len(rows) == 0 {
		return nil
	} else if name == "" {
		return errors.New("struct tag name can not be null")
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	rt := reflect.TypeOf(data)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var fileds = make(map[string]reflect.Value)
	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Tag.Get(name)
		if key != "" && key != "-" {
			fileds[key] = rv.Field(i)
		}
	}

	for k, d := range fileds {
		v, ok := rows[k]
		if !ok {
			continue
		}

		switch d.Interface().(type) {

		case string:
			if !replace && d.Interface().(string) != "" {
				continue
			} else {
				d.SetString(v)
			}

		case int, int8, int16, int32, int64:
			if !replace && TypeInt(d.Interface()) != 0 {
				continue
			} else if n, err := strconv.ParseInt(v, 10, 64); err != nil {
				return fmt.Errorf("Parse %s to int64 error", v)
			} else {
				d.SetInt(n)
			}

		case uint, uint8, uint16, uint32, uint64:
			if !replace && TypeUint(d.Interface()) != 0 {
				continue
			} else if n, err := strconv.ParseUint(v, 10, 64); err != nil {
				return fmt.Errorf("Parse %s to uint64 error", v)
			} else {
				d.SetUint(n)
			}

		case float32, float64:
			if !replace && TypeFloat(d.Interface()) != 0 {
				continue
			} else if n, err := strconv.ParseFloat(v, 64); err != nil {
				return fmt.Errorf("Parse %s to float64 error", v)
			} else {
				d.SetFloat(n)
			}

		case bool:
			if !replace && d.Interface().(bool) == true {
				continue
			} else if strings.ToLower(v) == "true" || v == "0" {
				d.SetBool(true)
			} else {
				d.SetBool(false)
			}

		case time.Time:
			var tt time.Time
			if !replace && d.Interface() != tt {
				continue
			} else if stamp, err := ParseTime(v); err != nil {
				return fmt.Errorf("Parse %s to time error", v)
			} else {
				d.Set(reflect.ValueOf(stamp))
			}

		default:
			return fmt.Errorf("No support %s type", d.Kind())
		}
	}
	return nil
}

// TypeInt : all int type to int
func TypeInt(data interface{}) int {
	switch data.(type) {
	case int:
		return data.(int)
	case int8:
		return int(data.(int8))
	case int16:
		return int(data.(int16))
	case int32:
		return int(data.(int32))
	case int64:
		return int(data.(int64))
	default:
		return 0
	}
}

// TypeUint : all uint to uint
func TypeUint(data interface{}) uint {
	switch data.(type) {
	case uint:
		return data.(uint)
	case uint8:
		return uint(data.(uint8))
	case uint16:
		return uint(data.(uint16))
	case uint32:
		return uint(data.(uint32))
	case uint64:
		return uint(data.(uint64))
	default:
		return 0
	}
}

// TypeFloat : all float type to float64
func TypeFloat(data interface{}) float64 {
	switch data.(type) {
	case float32:
		return float64(data.(float32))
	case float64:
		return data.(float64)
	default:
		return 0
	}
}
