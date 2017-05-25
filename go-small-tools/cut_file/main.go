package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please put-in file same string, if null meat all：")
	data, _, _ := reader.ReadLine()
	fileNameReg := string(data)

	var length int
	var err error
	for {
		fmt.Print("Please put-in cat up line number：")
		data, _, _ = reader.ReadLine()
		if len(data) == 0 {
			fmt.Print("Error: line number is null, please put-in again! \n")
			continue
		}
		if length, err = strconv.Atoi(string(data)); err != nil {
			fmt.Println(err)
			return
		}
		break
	}

	dirs := make([]string, 0)
	if fileDirs, err := ioutil.ReadDir("./"); err != nil {
		panic(err)
	} else {
		for _, file := range fileDirs {
			dirs = append(dirs, checkDir(file, "./")...)
		}
	}

	for _, file := range dirs {
		if strings.Contains(file, fileNameReg) {
			fmt.Printf("-------------- %s -------------- \n", file)
			if err := cutAddWriteFile(file[2:], length); err == nil {
				fmt.Println("OK!")
			} else {
				fmt.Println("Error:", err)
			}
		}
	}
}

func checkDir(file os.FileInfo, path string) []string {
	dirs := make([]string, 0)
	if file.IsDir() {
		path += file.Name() + "/"
		fileDirs, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		} else {
			for _, file := range fileDirs {
				dirs = append(dirs, checkDir(file, path)...)
			}
		}
		return dirs
	} else {
		return []string{path + file.Name()}
	}
}

func cutAddWriteFile(fileName string, length int) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println(file.Name())
	fileArr, err := csv.NewReader(file).ReadAll()
	if len(fileArr) <= length {
		return errors.New("file line number small than cut length!")
	}

	number := len(fileArr) / length

	for i := 0; i <= number; i++ {
		f, err := os.Create(fmt.Sprintf("new%d_%s", i+1, file.Name()))
		if err != nil {
			fmt.Println("create file was wrong:", i, err)
			return err
		}
		defer f.Close()

		f.WriteString("\xEF\xBB\xBF") // write UTF-8 BOM
		w := csv.NewWriter(f)

		if i != number {
			w.WriteAll(fileArr[i*length : (i+1)*length])
		} else {
			w.WriteAll(fileArr[number*length:])
		}
	}
	return nil
}
