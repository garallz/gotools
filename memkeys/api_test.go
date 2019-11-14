package memkeys

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	mem, err := NewCache("2MB", 2)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(*mem)
	}

	mem.SetWithExpire("test", "dafkjasdlfka", time.Second*2)
	mem.SetWithExpire("test", "dafkjasdlfka", time.Second*2)
	mem.SetWithExpire("one", "dafkjasdlfka", time.Second*2)
	mem.Set("test", "dafkjasdlfka")
	mem.Set("two", "dafkjasdlfka")
	mem.SetWithExpire("three", "33333333333", time.Second*2)

	value, ok := mem.Get("three")
	fmt.Println(value, ok)

	fmt.Println(mem.KeysNum())

	value, ok = mem.get("test")
	fmt.Println(value, ok)

	mem.del("test")
	fmt.Println(mem.KeysNum())

	time.Sleep(time.Second * 3)

	fmt.Println(mem.memSize)

	value, ok = mem.get("three")
	fmt.Println(value, ok)
}

// 2066326               646 ns/op             243 B/op         4 allocs/op
func BenchmarkSet(b *testing.B) {
	mem, err := NewCache("2GB", 2)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		mem.SetWithExpire(strconv.FormatInt(int64(i), 10), "dafkjasdlfka", time.Second*2)
		// mem.Set(strconv.FormatInt(int64(i), 10), "ahfghfghsfghs")
	}
}

// delete null key
// 25208442                46.6 ns/op             1 B/op         0 allocs/op
func BenchmarkDel(b *testing.B) {
	mem, err := NewCache("2MB", 2)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		mem.Del("test")
	}
}
