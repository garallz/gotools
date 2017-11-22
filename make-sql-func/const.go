package main

const head string = "//%s\npackage %s\n\nimport (\n\t\"database/sql\"%s\n)\n\ntype %sTable struct{\n%s\n}\n\n"

const (
	constInsert       string = "const Insert%s = \"INSERT INTO `%s` (%s) VALUES (%s)\""
	constDeleteIndex  string = "const DeleteIndex%s = \"DELETE FROM `%s` WHERE `%s` = ?\""
	constDeleteWhere  string = "const DeleteWhere%s = \"DELETE FROM `%s` WHERE \""
	constUpdateIndex  string = "const UpdateIndex%s = \"UPDATE `%s` SET %s WHERE `%s` = ?\""
	constUpdateUnique string = "const UpdateUnique%s = \"UPDATE `%s` SET %s WHERE %s\""
	constOnDuplicate  string = "const DuplicateWhere%s = \"INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE \""
	constDupliUnique  string = "const DuplicateUnique%s = \"INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s\""
	constSelectAll    string = "const SelectAll%s = \"SELECT %s FROM `%s`\""
	constSelectIndex  string = "const SelectIndex%s = \"SELECT %s FROM `%s` WHERE `%s` = ?\""
	constSelectWhere  string = "const SelectWhere%s = \"SELECT %s FROM `%s` WHERE \""
)

// 1&3: upTable
// 2&4: same sql data.
const insertRowFunc string = `
// Insert one row.
func InsertRow%s(db *sql.DB, %s interface{}) error {
	_, err := db.Exec(Insert%s, %s)
	return err
}`

// 1&2&3: upTable
// 4: array data.
const insertArrFunc string = `
// Insert array data with struct.
func Insert%sArray(db *sql.DB, data []*%sTable) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, row := range data {
			if _, err = tx.Exec(Insert%s,
				%s,
			); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

// 1&2&3: upTable
// 4: scan string
const queryIndexFunc string = `
// Query one row by index
func Query%sIndex(db *sql.DB, index interface{}) (*%sTable, error) {
	var row = new(%sTable)
	var err = db.QueryRow(SelectIndex%s, index).Scan(
	%s,
	)
	return row, err
}`

// 1&2&3&4&5: upTable
// 6: scan string
const queryAllFunc string = `
// Get all table rows.
func Query%sAll(db *sql.DB) ([]*%sTable, error) {
	r, err := db.Query(SelectAll%s)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		for r.Next() {
			var row = &%sTable{}
			if err = r.Scan(
				%s,
			); err != nil {
				return result, err
			} else {
				result = append(result, row)
			}
		}
		return result, nil
	}
}`

// 1&2&3&4&5: upTable
// 6: scan string
const queryAllWhereFunc string = `
// Query data by where query.
func Query%sWhere(db *sql.DB, where string, query ...interface{}) ([]*%sTable, error) {
	r, err := db.Query(SelectWhere%s + where, query...)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		for r.Next() {
			var row = &%sTable{}
			if err = r.Scan(
				%s,
			); err != nil {
				return result, err
			} else {
				result = append(result, row)
			}
		}
		return result, nil
	}
}`

// 1&2: upTable
const deleteIndexFunc string = `
// Delete one row data by index.
func Delete%sIndex(db *sql.DB, index interface{}) error {
	_, err := db.Exec(DeleteIndex%s, index)
	return err
}`

// 1&2: upTable
const deleteArrayIndexFunc string = `
// Delete data by index array.
func Delete%sArrayIndex(db *sql.DB, data []interface{}) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, index := range data {
			if _, err := tx.Exec(DeleteIndex%s, index); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

// 1&2: upTable
const deleteWhereFunc string = `
// Delete some data by when where == query.
func Delete%sWhere(db *sql.DB, where string, query ...interface{}) error {
	_, err := db.Exec(DeleteWhere%s + where, query...)
	return err
}`

// 1: upTable
// 2&3: table fields (small Camel-Case name)
const updateIndexFunc string = `
// Update one row by index.
func Update%sIndex(db *sql.DB, index, %s interface{}) error {
	_, err := db.Exec(UpdateIndex%s, index, %s)
	return err
}`

// 1&2&3: upTable
// 4: sql exec data.
const updateUniqueFunc string = `
// Update table set where unique fields.
func Update%sUnique(db *sql.DB, data []*%sTable) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, row := range data {
			if _, err = tx.Exec(UpdateUnique%s, 
				%s,
			); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

// 1&2&3: upTable
// 4: tabel value string
const duplicateArrayUniqueFunc string = `
// Mysql: On Duplicate Key Update
// insert or update by unique and index.
func Duplicate%sUnique(db *sql.DB, data []*%sTable) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, row := range data {
			if _, err = tx.Exec(DuplicateUnique%s, 
				%s,
			); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

// 1: Function name
// 2: upTable
// 3: query string
// 4: exec data.
const sqlExecArrayFunc string = `
func %s(db *sql.DB, data []*%sTable) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, row := range data {
			if _, err = tx.Exec(%s, %s); err != nil {
				return err
			}
		}
	}
	return tx.Commit
}`

// Generate common function file.
const fileCommonConst = `// Sql Common Set
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
