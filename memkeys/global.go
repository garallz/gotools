package memkeys

import (
	"time"
)

var initstatus bool = false
var memdata *Memory

// InitCache is init a global Key-Value Memory Cache
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func InitCache(maxMem string, interval int64) {
	if initstatus {
		return
	}

	var err error
	memdata, err = newCache(maxMem, interval)
	if err != nil {
		panic(err)
	}
	initstatus = true
}

// Set store Key-Value
func Set(key string, value interface{}) {
	memdata.set(key, value, 0)
}

// SetWithExpire store Key-Value and make expire time
func SetWithExpire(key string, value interface{}, duration time.Duration) {
	memdata.set(key, value, int64(duration))
}

// Get key value
func Get(key string) (interface{}, bool) {
	return memdata.get(key)
}

// Del is delete keys
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

// FlushAll : flush all Key-Value Memory
func FlushAll() bool {
	return memdata.flush()
}

// KeysNum is all keys number
func KeysNum() int64 {
	return int64(memdata.keys)
}

// Exist check key exist
func Exist(key string) bool {
	return memdata.exist(key)
}

// MemorySize get memory size
func MemorySize() int64 {
	return memdata.memSize()
}

// MaxMemWarn default: log.Print("Memory overflow maximum preset")
func MaxMemWarn(function func()) {
	memdata.function = function
}
