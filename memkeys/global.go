package memkeys

import (
	"time"
)

var memdata *Memory

// 初始化一个全局的 Key-Value Memory
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func InitMemory(maxMem string, interval int64) {
	var err error
	memdata, err = NewCache(maxMem, interval)
	if err != nil {
		panic(err)
	}
}

// 存储一对 Key-Value 键值
func Set(key string, value interface{}) {
	memdata.set(key, value, 0)
}

// 存储一对 Key-Value 键值并设定过期时间
func SetWithExpire(key string, value interface{}, duration time.Duration) {
	memdata.set(key, value, int64(duration))
}

// 获取 Key 的值
func Get(key string) (interface{}, bool) {
	return memdata.get(key)
}

// 删除一组 Key
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

// 清空 Key-Value Memory
func FlushAll() bool {
	return memdata.flush()
}

// 获取 Key 的条数
func KeysNum() int64 {
	return int64(memdata.keys)
}

// 内存溢出最大预设值的警报程序
// default: log.Print("Memory overflow maximum preset")
func MaxMemWarn(function func()) {
	memdata.function = function
}
