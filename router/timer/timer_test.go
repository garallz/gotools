package logfile

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	if err := NewTimer("1s", 3, true, funcprint); err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 6)
}

func funcprint() {
	fmt.Println("Testing")
}
