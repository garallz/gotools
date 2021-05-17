package table

import (
	"database/sql"
	"errors"
	"reflect"
)

type DBName string

const (
	MYSQL_TYPE    DBName = "mysql"
	POSTGRES_TYPE DBName = "postgres"
)

type ParseTimeType string

const (
	ParseTimeUTC      ParseTimeType = "UTC"
	ParseTimeLocation ParseTimeType = "Location"
)

var (
	databaseName   = MYSQL_TYPE
	defaultTagName = "db"
	parseTimeType  = ParseTimeLocation
	defaultSqlDB   *sql.DB
)

var (
	ErrValueIsNull = errors.New("values can't be null")
	ErrIndexIsNull = errors.New("index not set")
)

// TableStruct function use by this struct
type TableStruct struct {
	data  interface{}
	table string
	db    *sql.DB
	name  DBName // database name: [mysql, postgres], default: mysql
	index string // database index name
	tag   string // struct tag name

	elem map[string]reflect.Value
}

// NewTable make new table interface to use
func NewTable(data interface{}, table, tag string, db *sql.DB) *TableStruct {
	if table == "" {
		panic("Table name Can't be Null")
	}
	if db == nil {
		panic("sql.DB is null")
	}

	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		panic("This not ptr interface")
	}

	if tag == "" {
		tag = defaultTagName
	}

	elem := GetTagValue(data, tag)
	if len(elem) == 0 {
		panic("Tag Values is Null")
	}

	return &TableStruct{
		data:  data,
		table: table,
		db:    db,
		name:  databaseName,
		tag:   tag,
		elem:  elem,
	}
}

func NewDefault(data interface{}, table string) *TableStruct {
	if table == "" || defaultSqlDB == nil || defaultTagName == "" {
		panic("Set Default Values Wrong")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		panic("This not ptr interface")
	}
	elem := GetTagValue(data, defaultTagName)
	if len(elem) == 0 {
		panic("Tag Values is Null")
	}
	return &TableStruct{
		data:  data,
		table: table,
		db:    defaultSqlDB,
		name:  databaseName,
		tag:   defaultTagName,
		elem:  elem,
	}
}

// SetDefaultSqlDB : setting default database connect
func SetDefaultSqlDB(db *sql.DB) {
	defaultSqlDB = db
}

// SetDefaultDBName : setting default database, defaut [mysql]
func SetDefaultDBName(name DBName) {
	databaseName = name
}

// SetDefaultTagName : setting default data struct reflect tag name, defaut [db]
func SetDefaultTagName(tag string) {
	defaultTagName = tag
}

// SetParseTimeType : set parse time string to time type
func SetParseTimeType(name ParseTimeType) {
	parseTimeType = name
}

// SetDBName : [postgres, mysql]
// Default: mysql
func (d *TableStruct) SetDBName(name DBName) {
	d.name = name
}

// SetIndex : set data struct select or update index name: (uniq_index)
func (d *TableStruct) SetIndex(index string) {
	if _, ok := d.elem[index]; ok {
		d.index = index
	} else {
		rt := reflect.TypeOf(d.data).Elem()
		if rf, ok := rt.FieldByName(index); ok {
			d.index = rf.Tag.Get(d.tag)
		} else {
			panic("Set Index not Found Value")
		}
	}
}

// GetTable : get database table name
func (d *TableStruct) GetTable() string {
	return d.table
}

// GetDB : get connect this table *sql.DB
func (d *TableStruct) GetDB() *sql.DB {
	return d.db
}

// GetIndex : get table index set
func (d *TableStruct) GetIndex() string {
	return d.index
}

// GetIndexValue : get table index set value
func (d *TableStruct) GetIndexValue() interface{} {
	if d.index == "" {
		panic("Index Not Set")
	} else {
		return d.elem[d.index].Interface()
	}
}
