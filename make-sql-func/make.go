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
	Package string
	Data    []*SqlData
}

type SqlData struct {
	Table   string
	upTable string
	Index   string
	Message string
	Model   []int
	Fields  []FieldsData
}

type FieldsData struct {
	Name   string
	upName string
	Type   string
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

	makeConstFile(data)

	for _, row := range data.Data {
		out, err := os.OpenFile(row.Table+".go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(err)
		}

		var values string

		row.upTable = CamelCaseString(row.Table)
		for i, _ := range row.Fields {
			row.Fields[i].upName = CamelCaseString(row.Fields[i].Name)
			values += row.Fields[i].upName + "\t" + row.Fields[i].Type + "\n"
		}
		str := fmt.Sprintf(head, row.Message, data.Package, row.upTable, values)

		str += checkModel(row)

		out.WriteString(str)
		out.Close()
	}

	exec.Command("sh", "-c", "go fmt").Run()
}

func makeConstFile(data *SqlFunc) {
	out, err := os.OpenFile("sql_const.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	var context = fmt.Sprintf(ConstFile, data.Package)
	out.WriteString(context)
	out.Close()
}

func checkModel(row *SqlData) string {
	var queryString, scanString string
	for _, row := range row.Fields {
		queryString += "`" + row.Name + "`, "
		scanString += "\n&data." + row.upName + ", "
	}
	queryString = queryString[:len(queryString)-2]
	scanString = scanString + "\n"

	var str = setConst(row, queryString)

	if len(row.Model) == 0 {
		str += insertRowData(row) + "\n"
		str += insertArrData(row) + "\n"
		str += queryRowData(row, scanString) + "\n"
		str += queryAllData(row, scanString) + "\n"
		str += queryWhereData(row, scanString) + "\n"
		str += deleteDataById(row) + "\n"
		str += deleteWhereData(row) + "\n"
		str += updateRowData(row) + "\n"
	} else {
		for _, num := range DeleteSameInt(row.Model) {
			switch num {
			case 0:
				str += insertRowData(row) + "\n"
			case 1:
				str += queryRowData(row, scanString) + "\n"
			case 2:
				str += queryAllData(row, scanString) + "\n"
			case 3:
				str += queryWhereData(row, scanString) + "\n"
			case 4:
				str += deleteDataById(row) + "\n"
			case 5:
				str += deleteWhereData(row) + "\n"
			case 6:
				str += updateRowData(row) + "\n"
			default:
				continue
			}
		}
	}
	return str
}

func setConst(data *SqlData, queryString string) string {
	var num = strings.Repeat("?, ", len(data.Fields))
	var str, updateString string
	for _, row := range data.Fields {
		if row.Name == data.Index {
			continue
		}
		updateString += "`" + row.Name + "` = ?, "
	}

	str += fmt.Sprintf(constInsert, data.upTable, data.Table, queryString, num[:len(num)-2])
	str += fmt.Sprintf(constDeleteIndex, data.upTable, data.Table, data.Index)
	str += fmt.Sprintf(constDeleteWhere, data.upTable, data.Table)
	str += fmt.Sprintf(constUpdateIndex, data.upTable, data.Table, updateString[:len(updateString)-2], data.Index)
	str += fmt.Sprintf(constOnDuplicate, data.upTable, data.Table, queryString, num[:len(num)-2])
	str += fmt.Sprintf(constSelectAll, data.upTable, queryString, data.Table)
	str += fmt.Sprintf(constSelectIndex, data.upTable, queryString, data.Table, data.Index)
	str += fmt.Sprintf(constSelectWhere, data.upTable, queryString, data.Table)
	return str
}
