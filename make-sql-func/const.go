package sqlFunc

import (
	"database/sql"
	"errors"
)

const head string = "//%s\npackage %s\n\nimport (\n\"database/sql\"\n)\n\ntype %sTable struct{\n%s\n}\n\n"

const (
	constInsert      string = "const insert%s = \"INSERT INTO `%s` (%s) VALUES (%s)\"\n"
	constDeleteIndex string = "const delete%sIndex = \"DELETE FROM `%s` WHERE `%s` = ?\"\n"
	constDeleteWhere string = "const delete%sWhere = \"DELETE FROM `%s` WHERE \"\n"
	constUpdateIndex string = "const update%sIndex = \"UPDATE `%s` SET %s WHERE `%s` = ?\"\n"
	constOnDuplicate string = "const duplicate%s = \"INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE \"\n"
	constSelectAll   string = "const select%sAll = \"SELECT %s FROM `%s`\"\n"
	constSelectIndex string = "const select%sIndex = \"SELECT %s FROM `%s` WHERE `%s` = ?\"\n"
	constSelectWhere string = "const select%sWhere = \"SELECT %s FROM `%s` WHERE \"\n"
)

const ConstFile = "package %s\n\nimport (\n\"errors\"\n)\n\nvar (\n\tErrValuesLength = errors.New(\"Update values not eq fields!\")\n\t" +
	"ErrFieldNameNull = errors.New(\"Update data was wrong than field name is null!\")\n)\n\n" +
	"type UpdateData struct {\nName string\nValue interface{}\n}\n\n"

const (
	// 1&3: UpTable
	// 2&4: same sql data.
	insertRowFunc string = `
// Insert one row.
func InsertRow%s(db *sql.DB, %s interface{}) error {
	_, err := db.Exec(insert%s, %s)
	return err
}`

	// 1&2&3: UpTable
	// 4: array data.
	insertArrFunc string = `
// Insert array data with struct.
func Insert%sArray(db *sql.DB, data []*%sTable) error {
	if tx, err := db.Begin(); err != nil {
		return err
	} else {
		for _, row := range data {
			if _, err = tx.Exec(insert%s, %s); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

	// 1&2&3: UpTable
	// 4: scan string
	queryRowFunc string = `
// Query one row by index
func Query%sByIndex(db *sql.DB, index interface{}) (data *%sTable, err error) {
	err = db.QueryRow(select%sIndex, index).Scan(%s)
	return data, err
}`

	// 1&2&3&4&5: UpTable
	// 6: scan string
	queryAllFunc string = `
// Get all table rows.
func Query%sAll(db *sql.DB) (data []*%sTable, err error) {
	r, err := db.Query(select%sAll)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		var data = &%sTable{}
		for r.Next() {
			if err = r.Scan(%s); err != nil {
				return result, err
			} else {
				result = append(result, data)
			}
		}
		return result, nil
	}
}`

	// 1&2&3&4&5: UpTable
	// 6: scan string
	queryAllWhereFunc string = `
// Query data by where query.
func Query%sWhere(db *sql.DB, where string, query ...interface{}) (data []*%sTable, err error) {
	r, err := db.Query(select%sWhere + where, query)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		var data = &%sTable{}
		for r.Next() {
			if err = r.Scan(%s); err != nil {
				return result, err
			} else {
				result = append(result, data)
			}
		}
		return result, nil
	}
}`

	// 1&2: UpTable
	deleteRowFunc string = `
// Delete one row data by index.
func Delete%sByIndex(db *sql.DB, index interface{}) error {
	_, err := db.Exec(delete%sIndex, index)
	return err
}`

	// 1&2: UpTable
	deleteWhereFunc string = `
// Delete some data by when where == query.
func Delete%sWhere(db *sql.DB, where string, query ...interface{}) error {
	_, err := db.Exec(delete%sWhere + where, query)
	return err
}`

	// 1: UpTable
	updateRowFunc string = `
// Update one row by index.
func UpdateRow%s(db *sql.DB, index, %s interface{}) error {
	_, err := db.Exec(update%sIndex, index, %s)
	return err
}`
)

func Duplicate(db *sql.DB, fields []string, data ...interface{}) error {
	if len(fields) != len(data) {
		return errors.New("Update values not eq fields!")
	}

	return nil
}
