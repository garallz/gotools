package memkeys

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	ByteSizeKB = 1024
	ByteSizeMB = ByteSizeKB * 1024
	ByteSizeGB = ByteSizeMB * 1024
)

type Memory struct {
	// 总锁
	lock  sync.Mutex
	state bool

	// max key-value save bytes
	maxMem int64

	// not page to store key-values
	oneCache *MemoryData
	// all key-values store
	allCache []*MemoryData

	// 是否分页计算地址
	paging bool

	// 分页异或数
	pages int

	// 定时过期时间key-values删除处理存放
	expCache *ExpireCache

	// keys number
	keys int32

	// 内存超出预警程序
	function func()
}

type MemoryData struct {
	lock  sync.RWMutex
	cache map[string]*KeyValue
}

type KeyValue struct {
	key    string
	value  interface{}
	expire int64
}

// sort by expire time
type KVS []*KeyValue

func (t KVS) Len() int {
	return len(t)
}

func (t KVS) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t KVS) Less(i, j int) bool {
	return t[i].expire < t[j].expire
}

func parseMaxMemory(str string) (int64, error) {
	if str == "" {
		return 0, nil
	} else if len(str) < 3 {
		return 0, fmt.Errorf("Max Memory Size string: %s was wrong", str)
	}

	num, err := strconv.ParseFloat(str[:len(str)-2], 64)
	if err != nil {
		return 0, fmt.Errorf("Max Memory Size number: %s was error: %v", str, err)
	}

	switch strings.ToUpper(str[len(str)-2:]) {
	case "KB":
		return int64(num * ByteSizeKB), nil
	case "MB":
		return int64(num * ByteSizeMB), nil
	case "GB":
		return int64(num * ByteSizeGB), nil
	default:
		return 0, fmt.Errorf("Max Memory Size Unit: %s was wrong", str)
	}
}

func (m *Memory) memory(key string) *MemoryData {
	if m.paging {
		pag := pagination(key, m.pages)
		return m.allCache[pag]
	} else {
		return m.oneCache
	}
}

func (m *Memory) set(key string, value interface{}, duration int64) {
	var mem = m.memory(key)
	if m.state {
		m.lock.Lock()
		defer m.lock.Unlock()
	}

	mem.lock.RLock()
	v, ok := mem.cache[key]
	mem.lock.RUnlock()

	if !ok {
		atomic.AddInt32(&m.keys, 1)
		v = &KeyValue{
			key:   key,
			value: value,
		}
		mem.lock.Lock()
		mem.cache[key] = v
		mem.lock.Unlock()

		if duration != 0 {
			v.expire = time.Now().UnixNano() + duration
			m.expCache.putInto(v)
		}
	} else {
		if duration != 0 {
			v.expire = time.Now().UnixNano() + duration
		}
		v.value = value
	}
}

func (m *Memory) get(key string) (interface{}, bool) {
	var mem = m.memory(key)
	if m.state {
		m.lock.Lock()
		defer m.lock.Unlock()
	}

	mem.lock.RLock()
	result, ok := mem.cache[key]
	mem.lock.RUnlock()

	if ok {
		return result.value, true
	}
	return nil, false
}

func (m *Memory) del(key string) bool {
	var data = m.memory(key)
	if m.state {
		m.lock.Lock()
		defer m.lock.Unlock()
	}

	data.lock.Lock()
	defer data.lock.Unlock()

	if _, ok := data.cache[key]; ok {
		data.cache[key].expire = 0
		delete(data.cache, key)
		atomic.AddInt32(&m.keys, -1)
		return true
	} else {
		return false
	}
}

func (m *Memory) exist(key string) bool {
	var mem = m.memory(key)
	if m.state {
		m.lock.Lock()
		defer m.lock.Unlock()
	}

	mem.lock.RLock()
	_, ok := mem.cache[key]
	mem.lock.RUnlock()
	return ok
}

func (m *Memory) flush() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.state = true

	m.newExp()
	if m.paging {
		var cache = make([]*MemoryData, m.pages+1)
		for i, _ := range cache {
			cache[i] = &MemoryData{cache: make(map[string]*KeyValue)}
		}
		m.allCache = cache
	} else {
		m.oneCache = &MemoryData{cache: make(map[string]*KeyValue)}
	}
	atomic.SwapInt32(&m.keys, 0)
	m.state = false

	return true
}

func (m *Memory) keysNum() int64 {
	return int64(m.keys)
}

func (m *Memory) memSize() int64 {
	return int64(unsafe.Sizeof(*m))
}

// key not be null
func pagination(key string, pages int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) & pages
}

func (m *Memory) checkCacheSize() {
	if m.maxMem < m.memSize() {
		m.function()
	}
}
