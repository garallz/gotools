package timer

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(1)
	timer.SetTimeZone(time.Hour * 8)
	fmt.Println(time.Now(), time.Now().UTC().Unix())

	num, err := timer.AddTimerArgRun("00", 4, "time [:00] one", argDisplay)
	if err != nil {
		t.Error(err)
	} else {
		log.Println(num)
	}

	num1 := timer.AddTimerArgument(time.Second, -1, "Each Second", argDisplay)
	log.Println(num1)

	num2 := timer.AddTimerArgument(time.Second*2, -1, "Two Second", argDisplay)
	log.Println(num2)

	num3 := timer.AddTimerFunction(time.Millisecond*1500, -1, dirDisplay)
	log.Println(num3)

	time.Sleep(time.Second * 20)
	timer.DeleteTimer(num1)

	time.Sleep(time.Second * 40)
	timer.DeleteTimer(num2)

	time.Sleep(time.Minute * 5)
	timer.CloseAndExit()
}

func argDisplay(data interface{}) {
	fmt.Println(time.Now(), data)
}

func dirDisplay() {
	fmt.Println(time.Now(), "------------ display")
}
