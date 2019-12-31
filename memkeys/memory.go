package memkeys

import (
	"log"
	"time"
)

// make a new Key-Value Memory
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func newCache(maxMem string, interval int64) (*Memory, error) {
	var data = &Memory{}
	size, err := parseMaxMemory(maxMem)
	if err != nil {
		return nil, err
	}
	data.maxMem = size

	if data.maxMem <= 10*byteSizeMB {
		data.oneCache = &memoryData{cache: make(map[string]*keyValues)}
	} else {
		num := int(data.maxMem/(10*byteSizeMB) + 1)
		for i := 1; ; i *= 2 {
			if num <= i {
				data.pages = i - 1
				break
			}
		}
		data.paging = true
		data.allCache = make([]*memoryData, data.pages+1)
		for j, _ := range data.allCache {
			data.allCache[j] = &memoryData{cache: make(map[string]*keyValues)}
		}
	}

	data.newExp()
	if interval <= 0 {
		data.expCache.interval = 5
	} else {
		data.expCache.interval = interval
	}
	data.initExpire()
	data.function = func() { log.Print("Memory overflow maximum preset") }

	return data, nil
}

func (m *Memory) Set(key string, value interface{}) {
	m.set(key, value, 0)
}

func (m *Memory) SetWithExpire(key string, value interface{}, duration time.Duration) {
	m.set(key, value, int64(duration))
}

func (m *Memory) Get(key string) (interface{}, bool) {
	return m.get(key)
}

func (m *Memory) Del(keys ...string) bool {
	var ok bool = true
	for _, key := range keys {
		if key != "" {
			if !m.del(key) {
				ok = false
			}
		}
	}
	return ok
}

func (m *Memory) FlushAll() bool {
	return m.flush()
}

func (m *Memory) KeysNum() int64 {
	return int64(m.keys)
}

func (m *Memory) Exist(key string) bool {
	return m.exist(key)
}

func (m *Memory) MemorySize() int64 {
	return m.memSize()
}

func (m *Memory) MaxMemWarn(function func()) {
	m.function = function
}
