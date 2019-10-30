package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	InitTicker(1)

	time.Sleep(time.Second)

	fmt.Println(time.Now())
	if err := NewTimer("00", 4, false, "time [:00] one", display); err != nil {
		t.Error(err)
	}

	if err := NewTimer("30", 4, false, "time [:30] two", display); err != nil {
		t.Error(err)
	}

	//	NewRunDuration(time.Second*3, -1, nil, oneprint)

	//	NewRunTime(time.Now().Add(time.Second*2), nil, twoprint)

	// for i := 1; i < 10; i++ {
	// 	NewRunTime(time.Now().Add(time.Second), i*111, display)
	// 	time.Sleep(time.Millisecond * 200)
	// }

	time.Sleep(time.Minute * 3)
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
