package logfile

import (
	"sync"
)

const TimeUnit = 1000 * 1000 * 1000
const TimeFormatString = "2006-01-02 15:04:05"

const (
	SecondTimeUnit = 1 * TimeUnit
	MinuteTimeUnit = 60 * SecondTimeUnit
	HourTimeUnit   = 60 * MinuteTimeUnit
	DayTimeUnit    = 24 * HourTimeUnit
)

type TimerFunc struct {
	next     int64 // Next run time
	interval int64 // 时间间隔
	times    int
	fixed    bool //
	function func()
}

type SortFunc []*TimerFunc

func (t SortFunc) Len() int {
	return len(t)
}

func (t SortFunc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type ByNext struct{ SortFunc }

func (t ByNext) Less(i, j int) bool {
	return t.SortFunc[i].next > t.SortFunc[j].next
}

type ByInterval struct{ SortFunc }

func (t ByInterval) Less(i, j int) bool {
	return t.SortFunc[i].interval > t.SortFunc[j].interval
}
