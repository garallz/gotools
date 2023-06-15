package timer

import "time"

// parse time format
const (
	TimeFormatString = "2006-01-02 15:04:05"
	TimeFormatLength = 19
)

// time unit
const (
	TimeUnitSec  = 1000
	TimeUnitMin  = 60 * TimeUnitSec
	TimeUnitHour = 60 * TimeUnitMin
	TimeUnitDay  = 24 * TimeUnitHour
	TimeUnitWeek = 7 * TimeUnitDay
)

type runFunc interface {
	Run() bool
	Uuid() int64
	Next() int64
}

type funcOnlyBody struct {
	next     int64 // Next run time
	uniq     int64 // function uuid
	times    int   //
	interval int64 // 时间间隔
	function func()
}

func (f *funcOnlyBody) Uuid() int64 {
	if f == nil {
		return 0
	}
	return f.uniq
}
func (f *funcOnlyBody) Next() int64 {
	if f == nil {
		return 0
	}
	return f.next
}
func (f *funcOnlyBody) Run() bool {
	if f.times == 0 {
		return false
	}
	f.times--
	f.next += f.interval
	f.function()
	return true
}

type funcArgsBody struct {
	next     int64 // Next run time
	uniq     int64 // function uuid
	times    int   //
	interval int64 // 时间间隔
	message  interface{}
	function func(interface{})
}

func (f *funcArgsBody) Uuid() int64 { return f.uniq }
func (f *funcArgsBody) Next() int64 { return f.next }
func (f *funcArgsBody) Run() bool {
	if f.times == 0 {
		return false
	}
	f.times--
	f.next += f.interval
	f.function(f.message)
	return true
}

type funcBoolBody struct {
	next     int64 // Next run time
	uniq     int64 // function uuid
	times    int   //
	interval int64 // 时间间隔
	function func() bool
}

func (f *funcBoolBody) Uuid() int64 { return f.uniq }
func (f *funcBoolBody) Next() int64 { return f.next }
func (f *funcBoolBody) Run() bool {
	if f.times == 0 {
		return false
	}
	f.times--
	f.next += f.interval
	return f.function()
}

type funcArgsBool struct {
	next     int64 // Next run time
	uniq     int64 // function uuid
	times    int   //
	interval int64 // 时间间隔
	message  interface{}
	function func(interface{}) bool
}

func (f *funcArgsBool) Uuid() int64 { return f.uniq }
func (f *funcArgsBool) Next() int64 { return f.next }
func (f *funcArgsBool) Run() bool {
	if f.times == 0 {
		return false
	}
	f.times--
	f.next += f.interval
	return f.function(f.message)
}

type funcMonthDay struct {
	next     int64 // Next run time
	uniq     int64 // function uuid
	times    int   //
	function func() bool
}

func (f *funcMonthDay) Uuid() int64 { return f.uniq }
func (f *funcMonthDay) Next() int64 { return f.next }
func (f *funcMonthDay) Run() bool {
	if f.times == 0 {
		return false
	}
	f.times--
	f.next = time.UnixMilli(f.next).AddDate(0, 1, 0).UnixMilli()
	return f.function()
}
