package main

const ConstFile = `// Sql Common Set
package %s 

import (
	"database/sql"
	"errors"
)

var (
	ErrValuesLength		= errors.New("Update values not eq fields!")
	ErrFieldNameNull	= errors.New("Update data was wrong than field name is null!")
	ErrFieldsNull		= errors.New("Fields array was null.")
)

type FieldData struct {
	Name	string
	Value	interface{}
}

func SqlExec(db *sql.DB, query string, data ...interface{}) error {
	_, err := db.Exec(query, data...)
	return err
}
`
