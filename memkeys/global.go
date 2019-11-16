package memkeys

import (
	"time"
)

var memdata *Memory

// init a global Key-Value Memory Cache
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func InitCache(maxMem string, interval int64) {
	var err error
	memdata, err = newCache(maxMem, interval)
	if err != nil {
		panic(err)
	}
}

// store Key-Value
func Set(key string, value interface{}) {
	memdata.set(key, value, 0)
}

// store Key-Value and make expire time
func SetWithExpire(key string, value interface{}, duration time.Duration) {
	memdata.set(key, value, int64(duration))
}

// get key value
func Get(key string) (interface{}, bool) {
	return memdata.get(key)
}

// delete keys
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

// flush all Key-Value Memory
func FlushAll() bool {
	return memdata.flush()
}

// get all keys number
func KeysNum() int64 {
	return int64(memdata.keys)
}

// check key exist
func Exist(key string) bool {
	return memdata.exist(key)
}

// get memory size
func MemorySize() int64 {
	return memdata.memSize()
}

// default: log.Print("Memory overflow maximum preset")
func MaxMemWarn(function func()) {
	memdata.function = function
}
