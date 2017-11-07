package sqlFunc

import (
	"os/exec"
	"testing"
)

func TestMakeSqlFunction(t *testing.T) {
	err := exec.Command("sh", "-c", "rm node.go sql_const.go").Run()
	if err != nil {
		t.Error(err)
	}
	MakeSqlFunction("env.json")
}
