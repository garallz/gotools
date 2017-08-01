package llz_log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	l := LogStruct{
		TimeFormat: "15:04:05",
		FileTime:   TimeMinute,
	}
	l.Init()

	for i := 0; i <= 65; i++ {
		l.WriteError(i, i*10, i*100, i*1000)
		if i%2 == 0 {
			l.WriteDebug("yes", i)
		} else {
			l.WriteWarn("no", i)
		}
		time.Sleep(time.Second)
	}
}
