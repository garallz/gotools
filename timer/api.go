package timer

import (
	"time"
)

type TimerStruct interface {
	// set time-zone
	SetTimeZone(time.Duration)
	// turn time string to timer run
	AddTimerArgRun(string, int, interface{}, func(interface{})) (int64, error)
	// custom running time (less than day time)
	AddTimerCustom(string, int, func()) (int64, error)
	// add new boolean check function to ticker run
	AddTimerBoolean(time.Duration, func() bool) int64
	// add new boolean check function to ticker run with argument input
	AddTimerArgBool(time.Duration, interface{}, func(interface{}) bool) int64
	// add new timer function
	AddTimerFunction(time.Duration, int, func()) int64
	// add new timer function with argument input
	AddTimerArgument(time.Duration, int, interface{}, func(interface{})) int64
	// add new timer with each month day
	// put in month day, timestamp, times, function
	// eg: each month of 1th day, 00:00:00 time, 10 times, function
	AddTimerMonth(int, string, int, func() bool) (int64, error)

	// delete timer function by unique id
	DeleteTimer(int64)
	// cancel timer goroutine
	CloseAndExit()
}

// if null, default: 500ms
func NewTimer(duration time.Duration) TimerStruct {
	if duration == 0 {
		duration = time.Duration(time.Millisecond * 500)
	} else if duration.Milliseconds() < 100 {
		duration = time.Duration(time.Millisecond * 100)
	}
	var dur = duration.Milliseconds()

	var tmp = &timerData{
		root:  &chainNode{durmax: dur * 4, serial: 1},
		timer: time.NewTicker(duration),
	}
	var one = &chainNode{durmax: dur * 10, durmin: dur * 4, nummax: 3, serial: 2}
	tmp.root.next, one.last = one, tmp.root

	var two = &chainNode{durmax: dur * 100, durmin: dur * 10, nummax: 3, serial: 3}
	one.next, two.last = two, one

	var three = &chainNode{durmin: dur*100 + 1, nummax: 10, serial: 4}
	two.next, three.last, tmp.last = three, two, three

	go tmp.TimerRun()
	return tmp
}
