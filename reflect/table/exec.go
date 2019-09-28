package table

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Query : query row data by database where
func (d *TableStruct) Query(where string) error {
	var rows []string
	var values []interface{}
	for k, v := range d.elem {
		rows = append(rows, k)
		values = append(values, v.Addr().Interface())
	}

	sqlstr := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1",
		strings.Join(rows, ", "), d.table, where)
	return d.GetDB().QueryRow(sqlstr).Scan(values...)
}

// QueryArray : query array data by where
// use null struct ptr to query and return not ptr struct
func (d *TableStruct) QueryArray(where string) ([]interface{}, error) {
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

// QueryByIndex : query data by database index value
func (d *TableStruct) QueryByIndex(arg interface{}) error {
	if d.index == "" {
		panic("Index Not Set")
	} else if arg == nil {
		arg = d.GetIndexValue()
	}

	var rows []string
	var values []interface{}

	for k, v := range d.elem {
		rows = append(rows, k)
		values = append(values, v.Addr().Interface())
	}

	var sqlstr string
	if d.name == POSTGRES_TYPE {
		sqlstr = fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1",
			strings.Join(rows, ", "), d.table, d.index)
	} else {
		sqlstr = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? LIMIT 1",
			strings.Join(rows, ", "), d.table, d.index)
	}
	return d.GetDB().QueryRow(sqlstr, arg).Scan(values...)
}

// Insert data to database
func (d *TableStruct) Insert() error {
	var keys, line []string
	var values []interface{}
	var num = 1

	for key, value := range d.elem {
		keys = append(keys, key)
		values = append(values, value.Interface())
		line = append(line, d.field(num))
		num++
	}

	sqlstr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", d.table,
		strings.Join(keys, ", "), strings.Join(line, ", "))
	_, err := d.GetDB().Exec(sqlstr, values...)
	return err
}

func (d *TableStruct) field(num int) string {
	if d.name == POSTGRES_TYPE {
		return "$" + strconv.FormatInt(int64(num), 10)
	} else {
		return "?"
	}
}

// UpdateWithTag : if rows is null, default all valuse to update
// if rows not null, update rows values name
// type Data struct { Data string `db:"data"`}
// rows is data (tag name) not Data (value name)
func (d *TableStruct) UpdateWithTag(rows []string, where string) error {
	keys, values := d.filterSet(rows)

	sqlstr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", d.table,
		strings.Join(keys, ", "), where)

	_, err := d.GetDB().Exec(sqlstr, values...)
	return err
}

// UpdateValues : type Data struct { Data string `db:"data"`}
// rows is Data (value name) not data (tag name)
func (d *TableStruct) UpdateValues(rows []string, where string) error {
	rows = d.getTag(rows)
	return d.UpdateWithTag(rows, where)
}

// UpdateByIndexTag with rows key by index
// if rows is null, default all
// type Data struct { Data string `db:"data"`}
// rows is data (tag name) not Data (value name)
func (d *TableStruct) UpdateByIndexTag(rows []string) error {
	keys, values := d.filterSet(rows)
	values = append(values, d.GetIndexValue())

	var sqlstr string
	if d.name == POSTGRES_TYPE {
		sqlstr = fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", d.table,
			strings.Join(keys, ", "), d.index, len(values))
	} else {
		sqlstr = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", d.table,
			strings.Join(keys, ", "), d.index)
	}
	_, err := d.GetDB().Exec(sqlstr, values...)
	return err
}

// UpdateByIndex type Data struct { Data string `db:"data"`}
// rows is Data (value name) not data (tag name)
func (d *TableStruct) UpdateByIndex(rows []string) error {
	rows = d.getTag(rows)
	return d.UpdateByIndexTag(rows)
}

func (d *TableStruct) filterSet(rows []string) ([]string, []interface{}) {
	var keys []string
	var values []interface{}
	var num = 1

	if len(rows) > 0 {
		for _, row := range rows {
			if value, ok := d.elem[row]; ok {
				keys = append(keys, d.set(row, num))
				values = append(values, value.Interface())
				num++
			}
		}
	} else {
		for key, value := range d.elem {
			if d.index != "" && key == d.index {
				continue
			}
			keys = append(keys, d.set(key, num))
			values = append(values, value.Interface())
			num++
		}
	}
	return keys, values
}

func (d *TableStruct) set(k string, n int) string {
	if d.name == POSTGRES_TYPE {
		return fmt.Sprintf("%s = $%d", k, n)
	} else {
		return k + " = ?"
	}
}

// GetTagValue : convert data struct ptr to map reflect.Value
func GetTagValue(data interface{}, tag string) map[string]reflect.Value {
	st := reflect.TypeOf(data).Elem()
	rv := reflect.ValueOf(data).Elem()

	var key string
	var result = make(map[string]reflect.Value)

	for i := 0; i < st.NumField(); i++ {
		key = st.Field(i).Tag.Get(tag)
		if key == "" || key == "-" {
			continue
		}
		result[key] = rv.Field(i)
	}
	return result
}

func (d *TableStruct) getTag(rows []string) []string {
	if len(rows) == 0 {
		return nil
	}

	var result []string
	rt := reflect.TypeOf(d.data).Elem()
	for _, row := range rows {
		if rg, ok := rt.FieldByName(row); ok {
			result = append(result, rg.Tag.Get(d.tag))
		}
	}
	return result
}
