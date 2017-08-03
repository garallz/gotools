# Log (English)

## Use Way
```
    data := llz_log.LogStruct{}  // Initialized structure and assignment.
    data.Init()                  // Initialized func.
    data.WriteError("message")   // Write log data with level.
```

## Struct Definition
```
// It can be null.
type LogStruct struct {
            // if true, mean log data first put in cache, than cache full put in file.
            // if false, mean log data put in file as first time.
            // when LogStruct was null, Cache is true.
	Cache bool
            // cache save size (byte).
            // when cache was true and cache was null, mean cache eq 1024*1024 byte.
	CacheSize int
            // log time format, eg: "2006-01-02 15:04:05"
            // when TimeFormat was null, mean eq "2006-01-02 15:04:05",
 	TimeFormat string
            // log file pre name.
            // when FileName was null, mean pre name eq "log"
      	FileName string
            // file save path.
      	FilePath string
            // log save level.
            // when FileName was null, mean pre name eq LevelError.
	Level LogLevel
            // how long about file create.
            // when FileName was null, mean pre name eq TimeDay.
	FileTime LogTime
	    // whether create dir to save log file.
	Dir bool
}
```

# 日志 （Chinese）

## 使用方法
```
    data := llz_log.LogStruct{}  // 结构体初始化并定义赋值
    data.Init()                  // 程序初始化
    data.WriteError("message")   // 级别日志记录
```

## 结构定义
```
// 允许定义时为空，采用默认值。
type LogStruct struct {
            // 是否开启缓存写入
      	Cache bool
            // 缓存大小，单位为字节
            // 默认值为 1024*1024 bytes.
      	CacheSize int
            // 日志中的时间格式，默认为： "2006-01-02 15:04:05",
      	TimeFormat string
            // 日志文件前缀，默认为 log
      	FileName string
            // 日志路径，默认为当前路径
      	FilePath string
            // 日志记录等级，默认为 Error.
      	Level LogLevel
            // 文件生成间隔，默认为 TimeDay.
      	FileTime LogTime
	        // 是否开启创建日志文件夹， 默认不创建。
	Dir DirLevel
}
```
