package sqlFunc

import (
	"sort"
	"strings"
)

// Convert string to Camel-Case name.
func CamelCaseString(field string) string {
	if field == "" {
		return ""
	}
	var num = 0
	var result = make([]byte, 0)
	for i, r := range []byte(field) {
		// if rune == ('_', '-', ' '), delete and convert next rune.
		if r == 95 || r == 45 || r == 32 {
			num = i + 1
		} else {
			if i == num && r <= 122 && r >= 97 {
				result = append(result, r-32)
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

func SmallCamelCaseString(field string) string {
	if field == "" {
		return ""
	}
	var result = CamelCaseString(field)
	return strings.ToLower(result[:1]) + result[1:]
}

// Delete same int data with array.
func DeleteSameInt(data []int) []int {
	// If array data is null, return null.
	if len(data) == 0 {
		return nil
	}

	// Sort array data to delete.
	sort.Ints(data)
	var result = []int{data[0]}
	var num = data[0]
	for _, i := range data {
		if i == num {
			continue
		} else {
			num = i
			result = append(result, i)
		}
	}
	return result
}
