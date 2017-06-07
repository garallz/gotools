package lifekey

import (
	"fmt"
	"testing"
	"time"
)

func TestLifeKey(t *testing.T) {
	a := LifeData{}
	a.GcData(time.Second)

	a.SetAddData("one", "a", 5)
	fmt.Println("SetAddData: ", a.Get("a"))
	a.UpdateData("two", "a")
	fmt.Println("UpdateData :", a.Get("a"))
	a.Delete("a")
	fmt.Println("Delete a: ", a.Get("a"))

	a.Set("b", 3)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Println(i+1, "Second:", a.Check("b"))
	}
}
