package listen

import (
	"fmt"
	"testing"
)

func TestControlReq(t *testing.T) {
	SetEnvFilePath("./env.json")
	SetTimeout(3)

	data, err := StatusAll()
	if err != nil {
		t.Error(err)
	}
	for _, d := range data {
		fmt.Println(*d)
	}

	data, err := StatusByServer("server_two")
	if err != nil {
		t.Error(err)
	}
	for _, d := range data {
		fmt.Println(*d)
	}

	row, err := StatusByName("uniq_one")
	fmt.Println(row, err)

	result, err := ControlStart("uniq_one")
	fmt.Println(result, err)

	result, err := ControlByUnique("uniq_one")
	fmt.Println(result, err)

	result, err := ControlByServer("Server_one")
	fmt.Println(result, err)

}
