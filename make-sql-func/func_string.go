package sqlFunc

import (
	"fmt"
	"strings"
)

func insertRowData(data *SqlData) string {
	var values []string
	for _, field := range data.Fields {
		values = append(values, strings.ToLower(field.upName[:1])+field.upName[1:])
	}
	return fmt.Sprintf(insertRowFunc, data.upTable, strings.Join(values, ", "), data.upTable, strings.Join(values, ", "))
}

func insertArrData(data *SqlData) string {
	var values []string
	for _, field := range data.Fields {
		values = append(values, "row."+field.upName)
	}
	return fmt.Sprintf(insertArrFunc, data.upTable, data.upTable, data.upTable, strings.Join(values, ", "))
}

func queryRowData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryRowFunc, data.upTable, data.upTable, data.upTable, scanString)
}

func queryAllData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryAllFunc, data.upTable, data.upTable, data.upTable, data.upTable, data.upTable, scanString)
}

func queryWhereData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryAllWhereFunc, data.upTable, data.upTable, data.upTable, data.upTable, data.upTable, scanString)
}

func deleteDataById(data *SqlData) string {
	return fmt.Sprintf(deleteRowFunc, data.upTable, data.upTable)
}

func deleteWhereData(data *SqlData) string {
	return fmt.Sprintf(deleteWhereFunc, data.upTable, data.upTable)
}

func updateRowData(data *SqlData) string {
	var values []string
	for _, field := range data.Fields {
		if field.Name == data.Index {
			continue
		}
		values = append(values, strings.ToLower(field.upName[:1])+field.upName[1:])
	}
	return fmt.Sprintf(updateRowFunc, data.upTable, strings.Join(values, ", "), data.upTable, strings.Join(values, ", "))
}

func DuplicateData(data *SqlData) string {
	return ""
}
