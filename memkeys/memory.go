package memkeys

import (
	"log"
	"time"
)

// 初始化一个全局的 Key-Value Memory
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
	data.function = func() { log.Print("Memory overflow maximum preset") }

	return data, nil
}

// 存储一对 Key-Value 键值
func (m *Memory) Set(key string, value interface{}) {
	m.set(key, value, 0)
}

// 存储一对 Key-Value 键值并设定过期时间
func (m *Memory) SetWithExpire(key string, value interface{}, duration time.Duration) {
	m.set(key, value, int64(duration))
}

// 获取 Key 的值
func (m *Memory) Get(key string) (interface{}, bool) {
	return m.get(key)
}

// 删除一组 Key
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

// 清空 Key-Value Memory
func (m *Memory) FlushAll() bool {
	return m.flush()
}

// 获取 Key 的条数
func (m *Memory) KeysNum() int64 {
	return int64(m.keys)
}

// 内存溢出最大预设值的警报程序
// default: log.Print("Memory overflow maximum preset")
// 每 100 * interval 间隔时间触发一次检查
func (m *Memory) MaxMemWarn(function func()) {
	m.function = function
}
