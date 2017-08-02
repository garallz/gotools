package llz_log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// check struct data and supplement.
func (l *LogStruct) checkStruct() {
	if l == nil {
		l.Cache = true
		l.CacheSize = 1024 * 1024
		l.FileName = "log"
		l.FileTime = TimeDay
		l.Level = LevelError
		l.TimeFormat = "2006-01-02 15:04:05"
	} else {
		if l.FileName == "" {
			l.FileName = "log"
		}
		if l.FileTime == "" {
			l.FileTime = TimeDay
		}
		if l.Level == 0 {
			l.Level = LevelError
		}
		if l.TimeFormat == "" {
			l.TimeFormat = "2006-01-02 15:04:05"
		}
		if l.Cache && l.CacheSize == 0 {
			l.CacheSize = 1024 * 1024
		}
		if l.CacheSize != 0 && !l.Cache {
			l.Cache = true
		}
		if l.FilePath != "" {
			path := l.FilePath[len(l.FilePath)-1:]
			if path != "/" || path != `\` {
				if runtime.GOOS == "windows" {
					l.FilePath += `\`
				} else {
					l.FilePath += "/"
				}
			}
		}
	}
}

// open file and put in struct with *os.file
// init cache
func (l *LogStruct) open() {
	var err error
	name := l.FilePath + l.FileName + "." + time.Now().Format(fmt.Sprint(l.FileTime))

	l.file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("Open log file error: " + err.Error())
	}

	if l.Cache {
		l.buf = l.buf[:0]
	}

	go l.upFile()
}

// sleep time to make new file open.
func (l *LogStruct) upFile() {
	var last string
	var format = fmt.Sprint(l.FileTime)

	switch l.FileTime {
	case TimeMonth:
		last = time.Now().UTC().AddDate(0, 1, 0).Format(format)
	case TimeDay:
		last = time.Now().UTC().Add(time.Hour * 24).Format(format)
	case TimeHour:
		last = time.Now().UTC().Add(time.Hour * 1).Format(format)
	case TimeMinute:
		last = time.Now().UTC().Add(time.Minute).Format(format)
	}

	if stamp, err := time.Parse(format, last); err != nil {
		panic("Time parse error: " + err.Error())
	} else {
		l.stamp = stamp.UTC().Unix()
		if sleep := stamp.Sub(time.Now().UTC()).Seconds(); sleep > 5 {
			time.Sleep(time.Second * time.Duration(sleep-5))
		}
		l.tc = true
	}
}

// put log data in cache or file.
func (l *LogStruct) put(level string, msg ...interface{}) error {
	d := make([]string, len(msg)+1)
	d[0] = time.Now().Format(l.TimeFormat) + level
	for i, r := range msg {
		d[i+1] = fmt.Sprint(r)
	}

	f := []byte(strings.Join(d, " ") + "\n")

	var err error

	if l.tc {
		if err = l.check(); err != nil {
			return err
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Cache {
		l.buf = append(l.buf, f...)
		if len(l.buf) >= l.CacheSize {
			_, err = l.file.Write(l.buf)
			l.buf = l.buf[:0]
		}
	} else {
		_, err = l.file.Write(f)
	}
	return err
}

// check new file open.
func (l *LogStruct) check() error {
	if l.stamp <= time.Now().UTC().Unix() {
		l.mu.Lock()
		if l.Cache {
			if _, err := l.file.Write(l.buf); err != nil {
				return err
			}
			l.buf = l.buf[:0]
		}
		l.file.Close()

		var name = l.FilePath + l.FileName + "." + time.Now().Format(fmt.Sprint(l.FileTime))
		var err error
		l.file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic("Open new log file error: " + err.Error())
		}

		l.tc = false
		l.mu.Unlock()

		go l.upFile()
	}
	return nil
}
