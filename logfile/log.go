package logfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var LineByte = []byte("\n")

// check struct data and supplement.
func (l *LogStruct) checkStruct() *LogData {
	var d = &LogData{
		cache:  true,
		size:   1024 * 1024,
		name:   "log",
		time:   TimeDay,
		level:  LevelError,
		path:   "./",
		format: "2006-01-02 15:04:05",
		types:  DataTypeJson,
	}

	if l == nil {
		return d
	} else {
		if l.FileName != "" {
			d.name = l.FileName
		}
		if l.FileTime != "" {
			d.time = l.FileTime
		}
		if l.Level != 0 {
			d.level = l.Level
		}
		if l.TimeFormat != "" {
			d.format = l.TimeFormat
		}
		if l.CacheSize == 0 && !l.Cache {
			d.cache = false
		} else {
			if l.CacheSize != 0 {
				d.size = l.CacheSize
			}
			d.buf = d.buf[:0]
		}
		if l.Dir {
			d.dir = true
		}
		if l.FilePath != "" {
			d.path = l.FilePath
		}
		if l.DataType == DataTypeByte {
			d.types = DataTypeByte
		}
		return d
	}
}

// open file and put in struct with *os.file
func (l *LogData) open() {
	var err error
	name := filepath.Join(l.path, l.name+"."+time.Now().Format(fmt.Sprint(l.time)))

	if l.dir {
		d := l.createDir()
		name = filepath.Join(l.path, d, l.name+"."+time.Now().Format(fmt.Sprint(l.time)))
	}

	l.file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("Open log file error: " + err.Error())
	}
}

// sleep time to make new file open.
func (l *LogData) upFile() {
	var last string
	var format = fmt.Sprint(l.time)

	switch l.time {
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

// JsonData : DataType is json to wirte
type JsonData struct {
	Time  string `json:"time"`
	Level string `json:"level,omitempty"`
	Body  string `json:"body"`
}

// put log data and level in buffer.
func (l *LogData) put(level string, args []interface{}) error {
	if l.types == DataTypeJson {
		bts, _ := json.Marshal(&JsonData{
			Time:  time.Now().Format(l.format),
			Level: level,
			Body:  fmt.Sprint(args...),
		})
		return l.putByte(append(bts, LineByte...))
	} else {
		message := fmt.Sprintf("%s\t%s\t%s\n", time.Now().Format(l.format), level, fmt.Sprint(args...))
		return l.putByte([]byte(message + "\n"))
	}
}

// put log data and level in buffer by string.
func (l *LogData) putf(level string, msg string) error {
	if l.types == DataTypeJson {
		bts, _ := json.Marshal(&JsonData{
			Time:  time.Now().Format(l.format),
			Level: level,
			Body:  msg,
		})
		return l.putByte(append(bts, LineByte...))
	} else {
		msg = fmt.Sprintf("%s\t%s\t%s\n", time.Now().Format(l.format), level, msg)
		return l.putByte([]byte(msg))
	}
}

func (l *LogData) putPanic(bts []byte) {
	if l.cache {
		l.buf = append(l.buf, bts...)
		l.file.Write(l.buf)
	} else {
		l.file.Write(bts)
	}
	l.file.Close()
}

// put byte in cache or file.
func (l *LogData) putByte(bts []byte) error {
	var err error
	// check new file create time status.
	// TODO : make new file
	if l.tc {
		if err = l.check(); err != nil {
			return err
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cache {
		l.buf = append(l.buf, bts...)
		if len(l.buf) >= l.size {
			_, err = l.file.Write(l.buf)
			l.buf = l.buf[:0]
		}
	} else {
		_, err = l.file.Write(bts)
	}
	return err
}

// check new file open.
func (l *LogData) check() error {
	if l.stamp <= time.Now().UTC().Unix() {
		l.mu.Lock()
		if l.cache {
			if _, err := l.file.Write(l.buf); err != nil {
				return err
			}
			l.buf = l.buf[:0]
		}
		l.file.Close()

		l.open()
		l.tc = false
		l.mu.Unlock()

		go l.upFile()
	}
	return nil
}

// make dir about FileTime.
func (l *LogData) createDir() string {
	// Create log file dir with year.
	l.create(time.Now().Format("2006"))

	// Create log file dir with month.
	if l.time != TimeMonth {
		l.create(time.Now().Format("2006/01"))
	} else {
		return time.Now().Format("2006/")
	}
	// Create log file dir with day.
	if l.time != TimeDay {
		l.create(time.Now().Format("2006/01/02"))
	} else {
		return time.Now().Format("2006/01/")
	}
	// Create log file dir with hour.
	if l.time != TimeHour {
		l.create(time.Now().Format("2006/01/02/15"))
	} else {
		return time.Now().Format("2006/01/02/")
	}
	return time.Now().Format("2006/01/02/15/")
}

// check dir exist and create.
func (l *LogData) create(path string) {
	path = filepath.Join(l.path, path)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, 0666); err != nil {
				panic("Create log file dir error!")
			}
		}
	}
}
