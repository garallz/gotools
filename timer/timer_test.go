package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	InitTicker(2)

	time.Sleep(time.Second)

	fmt.Println(time.Now())
	if err := NewTimer("1s", 4, false, nil, eachprint); err != nil {
		t.Error(err)
	}

	NewRunDuration(time.Second*1, nil, oneprint)

	NewRunTime(time.Now().Add(time.Second*2), nil, twoprint)

	for i := 1; i < 10; i++ {
		NewRunTime(time.Now().Add(time.Second), i*111, display)
		time.Sleep(time.Millisecond * 200)
	}

	time.Sleep(time.Second * 6)
}

func display(data interface{}) {
	fmt.Println(time.Now(), data)
}

func eachprint(data interface{}) {
	fmt.Println(time.Now(), "each second to display")
}

func oneprint(data interface{}) {
	fmt.Println(time.Now(), "one second last display")
}

func twoprint(data interface{}) {
	fmt.Println(time.Now(), "two second last display")
}
