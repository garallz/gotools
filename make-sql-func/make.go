package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
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

// Make sql function import.
func MakeSqlFunction(fileName, path string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(file)
	}

	var data = &SqlFunc{}
	err = json.Unmarshal(file, data)
	if err != nil {
		panic(err)
	}

	// Check path string.
	if path != "" {
		if path[len(path)-1:] != `/` || path[len(path)-1:] != `\` {
			if runtime.GOOS == "windows" {
				path += `\`
			} else {
				path += `/`
			}
		}
	}

	// Make common sql function file.
	makeConstFile(data, path)

	for _, row := range data.Data {
		// Generate function file.
		newFile, err := os.OpenFile(path+row.Table+".go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		row.upTable = CamelCaseString(row.Table)
		var timeBool, values string
		for i, field := range row.Fields {
			// If table fields have time type, import time.Time.
			if field.Type == "time.Time" {
				timeBool = "\n\t\"time\""
			}
			row.Fields[i].upName = CamelCaseString(field.Name)
			values += row.Fields[i].upName + "\t" + field.Type + "\n"
		}

		newFile.WriteString(fmt.Sprintf(head, row.Message, data.Package, timeBool, row.upTable, values))
		newFile.WriteString(dealWithTable(row))

		newFile.Close()
	}

	return exec.Command("sh", "-c", "go fmt").Run()
}

// Make a sql_const.go file.
func makeConstFile(data *SqlFunc, path string) {
	file, err := os.OpenFile(path+"sql_const.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var context = fmt.Sprintf(fileCommonConst, data.Package)
	file.WriteString(context)
	file.Close()
}

// Deal with json tables data.
// when table value have time type, need add time package.
// convert table name and values name to camel-case string, using with struct.
// generate query and scan string to sql function to do.
func dealWithTable(row *SqlData) string {
	var queryString, scanString []string
	for _, field := range row.Fields {
		queryString = append(queryString, "`"+field.Name+"`")
		scanString = append(scanString, "&row."+field.upName)
	}

	if len(row.Model) == 0 {
		row.Model = []int{1, 2, 3, 4, 5}
	}
	return checkModelArray(row, strings.Join(queryString, ", "), strings.Join(scanString, ",\n"))
}

// check and range model to generate function.
// if index or unique index not null, will make the index function.
func checkModelArray(data *SqlData, queryString, scanString string) string {
	var num = strings.Repeat("?, ", len(data.Fields))
	var constString, funcString []string

	for _, model := range DeleteSameInt(data.Model) {
		switch model {

		case 1: // Insert Functions
			constString = append(constString, fmt.Sprintf(constInsert, data.upTable, data.Table, queryString, num[:len(num)-2]))
			funcString = append(funcString, insertRowData(data))
			funcString = append(funcString, insertArrData(data))

		case 2: // Delete Functions
			if data.Index != "" {
				constString = append(constString, fmt.Sprintf(constDeleteIndex, data.upTable, data.Table, data.Index))
				funcString = append(funcString, deleteIndexData(data))
				funcString = append(funcString, deleteArrayIndexData(data))
			}
			constString = append(constString, fmt.Sprintf(constDeleteWhere, data.upTable, data.Table))
			funcString = append(funcString, deleteWhereData(data))

		case 3: // Update Functions
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
				// If Unique == Index, continue.
				if len(data.Unique) == 1 && data.Unique[0] == data.Index {
					continue
				}
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

		case 4: //Select Functions
			if data.Index != "" {
				constString = append(constString, fmt.Sprintf(constSelectIndex, data.upTable, queryString, data.Table, data.Index))
				funcString = append(funcString, queryRowData(data, scanString))
			}
			constString = append(constString, fmt.Sprintf(constSelectAll, data.upTable, queryString, data.Table))
			constString = append(constString, fmt.Sprintf(constSelectWhere, data.upTable, queryString, data.Table))
			funcString = append(funcString, queryAllData(data, scanString))
			funcString = append(funcString, queryWhereData(data, scanString))

		case 5: // Duplicate Functions
			if len(data.Unique) != 0 || data.Index != "" {
				var setQuery, setData, execData []string
			DuplicateGo:
				for _, row := range data.Fields {
					execData = append(execData, "row."+row.upName)
					if row.Name == data.Index {
						continue
					} else {
						for _, x := range data.Unique {
							if row.Name == x {
								continue DuplicateGo
							}
						}
					}
					setQuery = append(setQuery, "`"+row.Name+"` = ?")
					setData = append(setData, "row."+row.upName)
				}
				constString = append(constString, fmt.Sprintf(constDupliUnique, data.upTable, data.Table, queryString, num[:len(num)-2], strings.Join(setQuery, ", ")))
				funcString = append(funcString, DuplicateUniqueData(data, strings.Join(append(execData, setData...), ",\n")))
			}
			constString = append(constString, fmt.Sprintf(constOnDuplicate, data.upTable, data.Table, queryString, num[:len(num)-2]))

		default: // Undefined Function
			continue
		}
	}

	return strings.Join(append(constString, funcString...), "\n")
}
