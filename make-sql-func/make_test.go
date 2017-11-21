package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestMakeSqlFunction(t *testing.T) {
	err := exec.Command("sh", "-c", "rm sql_const.go test.go").Run()
	if err != nil {
		t.Error(err)
	}

	if err = MakeSqlFunction("env.json", ""); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Make sql function success!")
	}
}
