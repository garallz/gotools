package memkeys

import (
	"time"
)

// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func NewCache(maxMem string, interval int64) (*Memory, error) {
	var data = &Memory{}
	size, err := parseMaxMemory(maxMem)
	if err != nil {
		return nil, err
	}
	data.maxMem = size

	if data.maxMem <= 10*ByteSizeMB {
		data.oneCache = &MemoryData{cache: make(map[string]*KeyValue)}
	} else {
		num := int(data.maxMem/(10*ByteSizeMB) + 1)
		for i := 1; ; i *= 2 {
			if num <= i {
				data.pages = i - 1
				break
			}
		}
		data.paging = true
		data.allCache = make([]*MemoryData, data.pages+1)
		for j, _ := range data.allCache {
			data.allCache[j] = &MemoryData{cache: make(map[string]*KeyValue)}
		}
	}

	data.newExp()
	if interval <= 0 {
		data.expCache.interval = 5
	} else {
		data.expCache.interval = interval
	}
	data.initExpire()

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
