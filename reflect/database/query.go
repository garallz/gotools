package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var ErrNotSureType = errors.New("This type of parsing is temporarily not supported.")

// QueryWithNull : query data by where default values null
func (d *TableStruct) QueryWithNull(where string) error {
	var rows []string
	var values []interface{}
	for k, v := range d.elem {
		rows = append(rows, k)
		switch v.Interface().(type) {
		case bool:
			values = append(values, &sql.NullBool{})
		case string:
			values = append(values, &sql.NullString{})
		case int8, uint8, int16, uint16, int32, uint32, int, uint, int64, uint64:
			values = append(values, &sql.NullInt64{})
		case float32, float64:
			values = append(values, &sql.NullFloat64{})
		case time.Time:
			values = append(values, &sql.NullTime{})
		default:
			return ErrNotSureType
		}
	}

	sqlstr := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1",
		strings.Join(rows, ", "), d.table, where)
	err := d.GetDB().QueryRow(sqlstr).Scan(values...)
	if err != nil {
		return err
	}

	for i, k := range rows {
		v := d.elem[k]
		if err = setSqlNull(v, values[i]); err != nil {
			return err
		}
	}
	return nil
}

// scan sql null struct to set reflect.value interface
func setSqlNull(value reflect.Value, data interface{}) error {
	switch value.Interface().(type) {
	case bool:
		if d, ok := data.(*sql.NullBool); ok && d.Valid {
			value.SetBool(d.Bool)
		} else {
			value.SetBool(false)
		}
	case string:
		if d, ok := data.(*sql.NullString); ok && d.Valid {
			value.SetString(d.String)
		} else {
			value.SetString("")
		}
	case int8, int16, int32, int, int64:
		if d, ok := data.(*sql.NullInt64); ok && d.Valid {
			value.SetInt(d.Int64)
		} else {
			value.SetInt(0)
		}
	case uint8, uint16, uint32, uint, uint64:
		if d, ok := data.(*sql.NullInt64); ok && d.Valid {
			value.SetUint(uint64(d.Int64))
		} else {
			value.SetUint(0)
		}
	case float32, float64:
		if d, ok := data.(*sql.NullFloat64); ok && d.Valid {
			value.SetFloat(d.Float64)
		} else {
			value.SetFloat(0)
		}
	case time.Time:
		if d, ok := data.(*sql.NullTime); ok && d.Valid {
			value.Set(reflect.ValueOf(d.Time))
		} else {
			value.Set(reflect.Value{})
		}
	default:
		return ErrNotSureType
	}
	return nil
}

// QueryArrayWithNull : query array data by where default values null
// use null struct ptr to query and return not ptr struct
func (d *TableStruct) QueryArrayWithNull(where string) ([]interface{}, error) {
	var rows []string
	var values []interface{}
	for k, v := range d.elem {
		rows = append(rows, k)
		values = append(values, v.Addr().Interface())
	}

	sqlstr := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(rows, ", "), d.table, where)

	res, err := d.GetDB().Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var result []interface{}
	for res.Next() {
		if err = res.Scan(values...); err != nil {
			return nil, err
		} else {
			result = append(result, reflect.ValueOf(d.data).Elem().Interface())
		}
	}
	return result, nil
}
