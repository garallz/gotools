package test

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

const (
	fileName string = "baimi.csv"
	length   int    = 1000000
)

func TestCutCsv(t *testing.T) {
	if err := cutCsvFile(); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Cut file success!")
	}
}

func cutCsvFile() error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileArr, err := csv.NewReader(file).ReadAll()

	number := len(file) / length

	for i := 0; i <= length; i++ {
		f, err := os.Create(fmt.Sprintf("file%d.csv", i+1))
		if err != nil {
			fmt.Println("Create file was wrong:", i, err)
			return err
		}
		defer f.Close()

		f.WriteString("\xEF\xBB\xBF") // write UTF-8 BOM
		w := csv.NewWriter(f)

		if i != length {
			w.WriteAll(fileArr[i*length : (i+1)*length])
		} else {
			w.WriteAll(fileArr[number*length:])
		}
	}
	return nil
}
