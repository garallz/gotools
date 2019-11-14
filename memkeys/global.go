package memkeys

import (
	"time"
)

var memdata *Memory

func InitMemory(maxMem string, interval int64) {
	var err error
	memdata, err = NewCache(maxMem, interval)
	if err != nil {
		panic(err)
	}
}

func Set(key string, value interface{}) {
	memdata.set(key, value, 0)
}

func SetWithExpire(key string, value interface{}, duration time.Duration) {
	memdata.set(key, value, int64(duration))
}

func Get(key string) (interface{}, bool) {
	return memdata.get(key)
}

func Del(keys ...string) bool {
	var ok bool = true
	for _, key := range keys {
		if key != "" {
			if !memdata.del(key) {
				ok = false
			}
		}
	}
	return ok
}

func FlushAll() bool {
	return memdata.flush()
}

func KeysNum() int64 {
	return int64(memdata.keys)
}
