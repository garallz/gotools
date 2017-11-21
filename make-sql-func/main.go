package main

import (
	"os"
)

func main() {
	var fields = os.Args

	if len(fields) < 2 {
		panic("Not enoght args!")
	} else if len(fields) == 2 {
		MakeSqlFunction(fields[1], "")
	} else {
		MakeSqlFunction(fields[1], fields[2])
	}
}
