package main

import (
	"fmt"
	"testing"
)

func TestGetPid(t *testing.T) {
	pid := getpid("garalinluzhi", "gocode")
	fmt.Println(pid)

	status := getProcessStatus(pid)
	fmt.Println(status)

	comp := getcomputer()
	fmt.Println(comp)
}
