package sqlFunc

const head string = "import (\n\"database/sql\"\n%s)\n\ntype %sTable struct{\n%s\n}\n\n"

const (
	constInsert       string = "const Insert%s = \"INSERT INTO `%s` (%s) VALUES (%s)\"\n"
	constDeleteIndex  string = "const DeleteIndex%s = \"DELETE FROM `%s` WHERE `%s` = ?\"\n"
	constDeleteWhere  string = "const DeleteWhere%s = \"DELETE FROM `%s` WHERE \"\n"
	constUpdateIndex  string = "const UpdateIndex%s = \"UPDATE `%s` SET %s WHERE `%s` = ?\"\n"
	constUpdateUnique string = "const UpdateUnique%s = \"UPDATE `%s` SET %s WHERE `%s` = ?\"\n"
	constOnDuplicate  string = "const Duplicate%s = \"INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE \"\n"
	constSelectAll    string = "const SelectAll%s = \"SELECT %s FROM `%s`\"\n"
	constSelectIndex  string = "const SelectIndex%s = \"SELECT %s FROM `%s` WHERE `%s` = ?\"\n"
	constSelectWhere  string = "const SelectWhere%s = \"SELECT %s FROM `%s` WHERE \"\n"
)

const ConstFile = "package %s\n\nimport (\n\"errors\"\n)\n\nvar (\n\t" +
	"ErrValuesLength = errors.New(\"Update values not eq fields!\")\n\t" +
	"ErrFieldNameNull = errors.New(\"Update data was wrong than field name is null!\")\n\t" +
	"ErrFieldsNull = errors.New(\"Fields array was null.\")\n)\n\n" +
	"type UpdateData struct {\nName string\nValue interface{}\n}\n\n"

const (
	// 1&3: UpTable
	// 2&4: same sql data.
	insertRowFunc string = `
// Insert one row.
func InsertRow%s(db *sql.DB, %s interface{}) error {
	_, err := db.Exec(Insert%s, %s)
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
			if _, err = tx.Exec(Insert%s, %s); err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}`

	// 1&2&3: UpTable
	// 4: scan string
	queryIndexFunc string = `
// Query one row by index
func Query%sIndex(db *sql.DB, index interface{}) (data *%sTable, err error) {
	err = db.QueryRow(SelectIndex%s, index).Scan(%s)
	return data, err
}`

	// 1&2&3&4&5: UpTable
	// 6: scan string
	queryAllFunc string = `
// Get all table rows.
func Query%sAll(db *sql.DB) (data []*%sTable, err error) {
	r, err := db.Query(SelectAll%s)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		for r.Next() {
			var data = &%sTable{}
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
	r, err := db.Query(SelectWhere%s + where, query)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()

		var result = make([]*%sTable, 0)
		for r.Next() {
			var data = &%sTable{}
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
	deleteIndexFunc string = `
// Delete one row data by index.
func Delete%sIndex(db *sql.DB, index interface{}) error {
	_, err := db.Exec(DeleteIndex%s, index)
	return err
}`

	// 1&2: UpTable
	deleteWhereFunc string = `
// Delete some data by when where == query.
func Delete%sWhere(db *sql.DB, where string, query ...interface{}) error {
	_, err := db.Exec(DeleteWhere%s + where, query)
	return err
}`

	// 1: UpTable
	updateIndexFunc string = `
// Update one row by index.
func UpdateRow%s(db *sql.DB, index, %s interface{}) error {
	_, err := db.Exec(UpdateIndex%s, index, %s)
	return err
}`

	// 1&2&3: UpTable
	// 4: tabel value string
	duplicateFunc string = `
// Mysql: On Duplicate Key Update
// fields is need update table field, and rows is update data,
// fields and rows the order needs one-to-one correspondence.
func Duplicate%sUnique(db *sql.DB, data *%sTable, fields []string, rows ...interface{}) error {
	if len(fields) != len(rows) {
		return ErrValuesLength
	} else if len(fields) == 0 {
		return ErrFieldsNull
	}
	
	var query string
	for _, field := range fields {
		query +=  field + " = ?, "
	}
	
	_, err := db.Exec(Duplicate%s + query[:len(query)-2], %s, rows)
	return err
}`

	// 1&2&3: UpTable
	// 4: tabel value string
	duplicateWhereFunc string = `
// Mysql: On Duplicate Key Update	
// fields is need update table field strint, eg: "name = ?, age = ?"
// fields and rows the order needs one-to-one correspondence.
func Duplicate%sIndex(db *sql.DB, data *%sTable, fields string, rows ...interface{}) error {
	_, err := db.Exec(Duplicate%s + fields, %s, rows)
	return err
}`
)
