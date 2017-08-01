package llz_log

import (
	"os"
	"sync"
)

type (
	LogLevel int
	LogTime  string
)

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

	buf   []byte
	file  *os.File
	stamp int64
	tc    bool
	mu    sync.Mutex
}

const (
	LevelInfo  LogLevel = 1
	LevelDebug LogLevel = 2
	LevelWarn  LogLevel = 3
	LevelError LogLevel = 4
	LevelFatal LogLevel = 5
)

const (
	TimeMonth  LogTime = "200601"
	TimeDay    LogTime = "20060102"
	TimeHour   LogTime = "2006010215"
	TimeMinute LogTime = "200601021504"
)
