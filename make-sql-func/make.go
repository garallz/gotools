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

	// table convert to big-camel-case name
	upTable string `json:"-`

	// table index
	Index string `json:"index`

	// sql automatically grow id, equal index
	AutoGrow string `json:"autogrow"`

	// table unique index.
	Unique []string `json:"unique"`

	// mysql duplicate update fields
	//Duplicate []string `json:"duplicate"`

	// table explanation message
	Message string `json:"message`

	// according to the need to generate the corresponding sql function.
	// eg: look at check model function.
	Model []int `json:"model`

	// table fields
	Fields []struct {
		// table field name
		Name string `json:"name`
		// field convert to big-camel-case name
		upName string `json:"-"`
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

		for i, field := range row.Fields {
			row.Fields[i].upName = CamelCaseString(field.Name)
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
		values += field.upName + "\t" + field.Type + "\n"
		queryString += "`" + field.Name + "`, "
		scanString += "\n&data." + field.upName + ", "
	}

	var str = fmt.Sprintf(head, timeBool, row.upTable, values)

	if len(row.Model) == 0 {
		row.Model = []int{1, 2, 3, 4, 5}
	}
	return str + checkModelArray(row, queryString[:len(queryString)-2], scanString+"\n")
}

func checkModelArray(data *SqlData, queryString, scanString string) string {
	var num = strings.Repeat("?, ", len(data.Fields))
	var constString, funcString []string

	for _, model := range DeleteSameInt(data.Model) {
		switch model {

		case 1: // Insert
			constString = append(constString, fmt.Sprintf(constInsert, data.upTable, data.Table, queryString, num[:len(num)-2]))
			funcString = append(funcString, insertRowData(data))
			funcString = append(funcString, insertArrData(data))

		case 2: // Delete
			if data.Index != "" {
				constString = append(constString, fmt.Sprintf(constDeleteIndex, data.upTable, data.Table, data.Index))
				funcString = append(funcString, deleteIndexData(data))
				funcString = append(funcString, deleteArrayIndexData(data))
			}
			constString = append(constString, fmt.Sprintf(constDeleteWhere, data.upTable, data.Table))
			funcString = append(funcString, deleteWhereData(data))

		case 3: // Update
			// When update the table, ignore index update.
			if data.Index != "" {
				var updateSet []string
				for _, row := range data.Fields {
					if row.Name == data.Index {
						continue
					}
					updateSet = append(updateSet, "`"+row.Name+"` = ?")
				}
				constString = append(constString, fmt.Sprintf(constUpdateIndex, data.upTable, data.Table, strings.Join(updateSet, ", "), data.Index))
				funcString = append(funcString, updateRowData(data))
			}
			// When update the table, ignore unique update.
			if len(data.Unique) != 0 {
				var setQuery, whereQuery, setData, whereData []string
			UpdateGo:
				for _, row := range data.Fields {
					for _, x := range data.Unique {
						if row.Name == x {
							whereQuery = append(whereQuery, "`"+row.Name+"` = ?")
							whereData = append(whereData, "row."+row.upName)
							continue UpdateGo
						}
					}
					setQuery = append(setQuery, "`"+row.Name+"` = ?")
					setData = append(setData, "row."+row.upName)
				}
				constString = append(constString, fmt.Sprintf(constUpdateUnique, data.upTable, data.Table, strings.Join(setQuery, ", "), strings.Join(whereQuery, ", ")))
				funcString = append(funcString, updateUniqueArrayData(data, strings.Join(append(setData, whereData...), ",\n")))
			}

		case 4: //Select
			if data.Index != "" {
				constString = append(constString, fmt.Sprintf(constSelectIndex, data.upTable, queryString, data.Table, data.Index))
				funcString = append(funcString, queryRowData(data, scanString))
			}
			constString = append(constString, fmt.Sprintf(constSelectAll, data.upTable, queryString, data.Table))
			constString = append(constString, fmt.Sprintf(constSelectWhere, data.upTable, queryString, data.Table))
			funcString = append(funcString, queryAllData(data, scanString))
			funcString = append(funcString, queryWhereData(data, scanString))

		case 5: // Duplicate
			if len(data.Unique) != 0 {
				var setQuery, setData, execData []string
			DuplicateGo:
				for _, row := range data.Fields {
					execData = append(execData, "row."+row.upName)
					for _, x := range data.Unique {
						if row.Name == x || row.Name == data.Index {
							continue DuplicateGo
						}
					}
					setQuery = append(setQuery, "`"+row.Name+"` = ?")
					setData = append(setData, "row."+row.upName)
				}
				constString = append(constString, fmt.Sprintf(constDupliUnique, data.upTable, data.Table, queryString, num[:len(num)-2], strings.Join(setQuery, ", ")))
				funcString = append(funcString, DuplicateUniqueData(data, strings.Join(append(execData, setData...), ",\n")))
			}
			//constString = append(constString, fmt.Sprintf(constOnDuplicate, data.upTable, data.Table, queryString, num[:len(num)-2]))
			//funcString = append(funcString, DuplicateData(data))
		default:
			continue
		}
	}

	return strings.Join(append(constString, funcString...), "\n")
}
