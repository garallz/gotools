package sqlFunc

import (
	"fmt"
	"strings"
)

func insertRowData(data *SqlData) string {
	var values []string
	for _, field := range data.Fields {
		values = append(values, SmallCamelCaseString(field.Name))
	}
	return fmt.Sprintf(insertRowFunc, data.upTable, strings.Join(values, ", "), data.upTable, strings.Join(values, ", "))
}

func insertArrData(data *SqlData) string {
	var values []string
	for _, field := range data.Fields {
		values = append(values, "row."+CamelCaseString(field.Name))
	}
	return fmt.Sprintf(insertArrFunc, data.upTable, data.upTable, data.upTable, strings.Join(values, ", "))
}

func queryRowData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryIndexFunc, data.upTable, data.upTable, data.upTable, scanString)
}

func queryAllData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryAllFunc, data.upTable, data.upTable, data.upTable, data.upTable, data.upTable, scanString)
}

func queryWhereData(data *SqlData, scanString string) string {
	return fmt.Sprintf(queryAllWhereFunc, data.upTable, data.upTable, data.upTable, data.upTable, data.upTable, scanString)
}

func deleteIndexData(data *SqlData) string {
	return fmt.Sprintf(deleteIndexFunc, data.upTable, data.upTable)
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
		values = append(values, SmallCamelCaseString(field.Name))
	}
	return fmt.Sprintf(updateIndexFunc, data.upTable, strings.Join(values, ", "), data.upTable, strings.Join(values, ", "))
}

func DuplicateData(data *SqlData) string {
	return ""
}
