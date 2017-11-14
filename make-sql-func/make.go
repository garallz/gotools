package sqlFunc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type SqlFunc struct {
	// package name.
	Package string `json:"package`

	// tables data.
	Data []*SqlData `json:"data`
}

type SqlData struct {
	// sql table name
	Table string `json:"table`

	// table big-camel-case name
	upTable string `json:"-`

	// table index
	Index string `json:"index`

	// sql automatically grow id, equal index
	AutoGrow string `json:"autogrow"`

	// table unique index.
	Unique []string `json:"unique"`

	// mysql duplicate update fields
	Duplicate []string `json:"duplicate"`

	// table explanation message
	Message string `json:"message`

	// according to the need to generate the corresponding sql function.
	// eg: look at check model function.
	Model []int `json:"model`

	// table fields
	Fields []struct {
		// table field name
		Name string `json:"name`
		// table field type
		// eg: int, int32, int64, fload64, string, time.Time, bool...
		Type string `json:"type`
	} `json:"fields`
}

func MakeSqlFunction(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(file)
	}

	var data = &SqlFunc{}
	err = json.Unmarshal(file, data)
	if err != nil {
		panic(err)
	}

	// Make common sql function file.
	makeConstFile(data)

	for _, row := range data.Data {
		newFile, err := os.OpenFile(row.Table+".go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		newFile.WriteString(fmt.Sprintf("//%s\npackage %s\n\n", row.Message, data.Package))
		newFile.WriteString(dealWithTable(row))

		newFile.Close()
	}

	exec.Command("sh", "-c", "go fmt").Run()
}

// Make a sql_const.go file.
func makeConstFile(data *SqlFunc) {
	file, err := os.OpenFile("sql_const.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var context = fmt.Sprintf(ConstFile, data.Package)
	file.WriteString(context)
	file.Close()
}

// Deal with json tables data.
// when table value have time type, need add time package.
// convert table name and values name to camel-case string, using with struct.
// generate query and scan string to sql function to do.
func dealWithTable(row *SqlData) string {
	row.upTable = CamelCaseString(row.Table)
	var queryString, scanString, timeBool, values string
	for _, field := range row.Fields {
		if field.Type == "time.Time" {
			timeBool = "\"time\"\n"
		}
		upName := CamelCaseString(field.Name)
		values += upName + "\t" + field.Type + "\n"
		queryString += "`" + field.Name + "`, "
		scanString += "\n&data." + upName + ", "
	}

	var str = fmt.Sprintf(head, timeBool, row.upTable, values)

	if len(row.Model) == 0 {
		row.Model = []int{1, 2, 3, 4}
	}
	return str + checkModelArray(row, queryString[:len(queryString)-2], scanString+"\n")
}

func checkModelArray(data *SqlData, queryString, scanString string) string {
	var num = strings.Repeat("?, ", len(data.Fields))
	var constString, funcString string

	for _, model := range DeleteSameInt(data.Model) {
		switch model {

		case 1: // Insert
			constString += fmt.Sprintf(constInsert, data.upTable, data.Table, queryString, num[:len(num)-2])
			funcString += insertRowData(data) + "\n"
			funcString += insertArrData(data) + "\n"

		case 2: // Delete
			if data.Index != "" {
				constString += fmt.Sprintf(constDeleteIndex, data.upTable, data.Table, data.Index)
				funcString += deleteIndexData(data) + "\n"
			}
			constString += fmt.Sprintf(constDeleteWhere, data.upTable, data.Table)
			funcString += deleteWhereData(data) + "\n"

		case 3: // Update
			// When update the table, ignore index update.
			if data.Index != "" {
				var updateString []string
				for _, row := range data.Fields {
					if row.Name == data.Index {
						continue
					}
					updateString = append(updateString, "`"+row.Name+"` = ?")
				}
				constString += fmt.Sprintf(constUpdateIndex, data.upTable, data.Table, strings.Join(updateString, ", "), data.Index)
				funcString += updateRowData(data) + "\n"
			}
			// When update the table, ignore unique update.
			if len(data.Unique) != 0 {
				var updateString []string
				for _, row := range data.Fields {
					for _, x := range data.Unique {
						if row.Name == x {
							continue
						}
					}
					updateString = append(updateString, "`"+row.Name+"` = ?")
				}
				constString += fmt.Sprintf(constUpdateUnique, data.upTable, data.Table, strings.Join(updateString, ", "), data.Index)
				funcString += updateRowData(data) + "\n"
			}

		case 4: //Select
			if data.Index != "" {
				constString += fmt.Sprintf(constSelectIndex, data.upTable, queryString, data.Table, data.Index)
				funcString += queryRowData(data, scanString) + "\n"
			}
			constString += fmt.Sprintf(constSelectAll, data.upTable, queryString, data.Table)
			constString += fmt.Sprintf(constSelectWhere, data.upTable, queryString, data.Table)
			funcString += queryAllData(data, scanString) + "\n"
			funcString += queryWhereData(data, scanString) + "\n"

		case 5: // Duplicate
			constString += fmt.Sprintf(constOnDuplicate, data.upTable, data.Table, queryString, num[:len(num)-2])
			funcString += DuplicateData(data) + "\n"
		default:
			continue
		}
	}

	return constString + funcString
}
