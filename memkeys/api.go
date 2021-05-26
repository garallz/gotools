package memkeys

import (
	"time"
)

// Cache insterface with memory key-values
type Cache interface {
	// set key-value, no time expire
	// Will keep the original expiration time
	Set(key string, val interface{})

	// set key-value with time expire
	// Update the original expiration time
	SetWithExpire(key string, val interface{}, expire time.Duration)

	// get value by key
	Get(key string) (interface{}, bool)

	// delete key-values
	Del(keys ...string) bool

	// check key exist
	Exist(key string) bool

	// flush all keys
	FlushAll() bool

	// return keys number
	KeysNum() int64

	// memory size
	MemorySize() int64

	// memory overflow maximum preset
	// default: log.Print("Memory overflow maximum preset")
	MaxMemWarn(function func())
}

// NewCache : make a new Key-Value Memory
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func NewCache(maxMem string, interval int64) (Cache, error) {
	if d, err := newCache(maxMem, interval); err != nil {
		return nil, err
	} else {
		var data Cache = d
		return data, nil
	}
}
