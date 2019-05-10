package listen

import (
	"fmt"
	"testing"
)

func TestControlReq(t *testing.T) {
	SetEnvFilePath("./env.json")

	data, err := StatusAll()
	if err != nil {
		t.Error(err)
	}
	for _, d := range data {
		fmt.Println(*d)
	}

	row, err := StatusByName("server_one")
	fmt.Println(row, err)

	result, err := ControlStart("server_one")
	fmt.Println(result, err)
}
