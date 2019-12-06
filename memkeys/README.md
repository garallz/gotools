# Memkeys

## Key-Values Memory Cache

```go
type Cache interface {
	// set key-value, no time expire
	// keep the original expiration time
	// 保持原有过期时间
	// 如果是新的将永远存活
	Set(key string, val interface{})

	// Add key-value with time expire
	// Update the original expiration time
	// 添加一个Key-Value缓存项，并且在expire时间之后过期
	// 更新原有的过期时间
	SetWithExpire(key string, val interface{}, expire time.Duration)

	// get value by key
	// 获取一个key值
	Get(key string) (interface{}, bool)

	// delete key-values
	// 删除一组key值对
	Del(keys ...string) bool

	// check key exist
	// 检测一个值是否存在
	Exist(key string) bool

	// flush all keys
	// 清空所有keys值
	FlushAll() bool

	// return keys number
	// 返回所有的keys的数量
	KeysNum() int64

	// memory size
	// 总占用内存大小
	MemorySize() int64

	// memory overflow maximum preset
	// default: log.Print("Memory overflow maximum preset")
	// 内存溢出最大预设值的警报程序
	// 每 100 * interval 间隔时间触发一次检查警报
	MaxMemWarn(function func())
}

// maxMem: Set Max Memory Size: ['512KB', '3MB', '432MB', '2GB']
// interval unit is 100ms, eg: 10 => 1s, default: 500ms
func NewCache(maxMem string, interval int64) (Cache, error)

// init a global key-values cache to use
func InitCache(maxMem string, interval int64)
```