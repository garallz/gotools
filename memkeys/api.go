package memkeys

import (
	"time"
)

type Cache interface {
	// set cache max memory store
	// size 是一个字符串。支持以下参数: 1KB，100KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool

	// set key-value, no time expire
	Set(key string, val interface{})

	// set key-value with time expire
	// 设置一个缓存项，并且在expire时间之后过期
	SetWithExpire(key string, val interface{}, expire time.Duration)

	// get value by key
	// 获取一个值
	Get(key string) (interface{}, bool)

	// delete key-values
	// 删除一个值
	Del(keys ...string) bool

	// check key exist
	// 检测一个值 是否存在
	Exist(key string) bool

	// flush all keys
	// 清空所有值
	FlushAll() bool

	// return keys number
	// 返回所有的key多少
	KeysNum() int64
}
