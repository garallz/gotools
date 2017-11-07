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
	Package string     `json:"package`
	Data    []*SqlData `json:"data`
}

type SqlData struct {
	Table   string       `json:"table`
	upTable string       `json:"-`
	Index   string       `json:"index`
	Message string       `json:"message`
	Model   []int        `json:"model`
	Fields  []FieldsData `json:"fields`
}

type FieldsData struct {
	Name string `json:"name`
	Type string `json:"type`
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
		out, err := os.OpenFile(row.Table+".go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		out.WriteString(fmt.Sprintf("//%s\npackage %s\n\n", row.Message, data.Package))
		out.WriteString(dealWithTable(row))

		out.Close()
	}

	exec.Command("sh", "-c", "go fmt").Run()
}

func makeConstFile(data *SqlFunc) {
	out, err := os.OpenFile("sql_const.go", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var context = fmt.Sprintf(ConstFile, data.Package)
	out.WriteString(context)
	out.Close()
}

func dealWithTable(row *SqlData) string {
	row.upTable = CamelCaseString(row.Table)
	var queryString, scanString, timeBool, values string = "", "", "", ""
	for _, field := range row.Fields {
		if field.Type == "time.Time" {
			timeBool = "\"time\"\n"
		}
		upName := CamelCaseString(field.Name)
		values += upName + "\t" + field.Type + "\n"
		queryString += "`" + field.Name + "`, "
		scanString += "\n&data." + upName + ", "
	}
	queryString = queryString[:len(queryString)-2]
	scanString = scanString + "\n"

	var str = fmt.Sprintf(head, timeBool, row.upTable, values)
	str += setConst(row, queryString)
	str += checkModel(row, scanString)
	return str
}

func checkModel(row *SqlData, scanString string) string {
	var str string
	if len(row.Model) == 0 {
		if row.Index != "" {
			str += queryRowData(row, scanString) + "\n"
			str += deleteIndexData(row) + "\n"
			str += updateRowData(row) + "\n"
		}
		str += insertRowData(row) + "\n"
		str += insertArrData(row) + "\n"
		str += queryAllData(row, scanString) + "\n"
		str += queryWhereData(row, scanString) + "\n"
		str += deleteWhereData(row) + "\n"
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
				str += deleteIndexData(row) + "\n"
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

	if data.Index != "" {
		for _, row := range data.Fields {
			if row.Name == data.Index {
				continue
			}
			updateString += "`" + row.Name + "` = ?, "
		}
		str += fmt.Sprintf(constDeleteIndex, data.upTable, data.Table, data.Index)
		str += fmt.Sprintf(constUpdateIndex, data.upTable, data.Table, updateString[:len(updateString)-2], data.Index)
		str += fmt.Sprintf(constSelectIndex, data.upTable, queryString, data.Table, data.Index)
	}
	str += fmt.Sprintf(constInsert, data.upTable, data.Table, queryString, num[:len(num)-2])
	str += fmt.Sprintf(constDeleteWhere, data.upTable, data.Table)
	str += fmt.Sprintf(constOnDuplicate, data.upTable, data.Table, queryString, num[:len(num)-2])
	str += fmt.Sprintf(constSelectAll, data.upTable, queryString, data.Table)
	str += fmt.Sprintf(constSelectWhere, data.upTable, queryString, data.Table)
	return str
}
