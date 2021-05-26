package timer

import (
	"sync"
	"sync/atomic"
)

type TimerType int

const (
	TypeIsNilFunc TimerType = iota
	TypeParamFunc
	TypeEventBool
)

type TimerStruct interface {
	Type() TimerType
	Next() int64
	Run()
	Del()
}

type TimerEventBool struct {
	next     int64 // Next run time
	interval int64 // 时间间隔
	times    int32 // run times
	function func(interface{}) bool
	msg      interface{}
	lock     sync.Mutex
}

// Type : event types
func (d *TimerEventBool) Type() TimerType {
	return TypeEventBool
}

// Next :
func (d *TimerEventBool) Next() int64 {
	return int64(d.next)
}

// Run : run function
func (d *TimerEventBool) Run() {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.times == 0 {
		return
	}
	if d.function(d.msg) {
		return
	} else if d.times >= 1 {
		if num := atomic.AddInt32(&d.times, -1); num == 0 {
			return
		}
	}
	// Input next
	atomic.AddInt64(&d.next, d.interval)
	putPools(d)
}

// Del :
func (d *TimerEventBool) Del() {
	d.times, d.next = 0, GetNowStamp()
}

// TimerIntervalFunc : Run the program at regular intervals
type TimerIntervalFunc struct {
	next     int64 // Next run time
	interval int64 // 时间间隔
	times    int32 // run times
	function func()
	lock     sync.Mutex
}

// Type : event types
func (d *TimerIntervalFunc) Type() TimerType {
	return TypeIsNilFunc
}

// Next :
func (d *TimerIntervalFunc) Next() int64 {
	return int64(d.next)
}

// Run : run function
func (d *TimerIntervalFunc) Run() {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch d.times {
	case 0:
		return
	case 1:
		d.function()
		return
	case -1:
		// forever
	default:
		atomic.AddInt32(&d.times, -1)
	}

	d.function()
	atomic.AddInt64(&d.next, d.interval)
	putPools(d)
}

// Del :
func (d *TimerIntervalFunc) Del() {
	d.times, d.next = 0, GetNowStamp()
}

// TimerParamFunc : Run the program at regular intervals
type TimerParamFunc struct {
	next     int64 // Next run time
	interval int64 // 时间间隔
	times    int32 // run times
	function func(param interface{})
	param    interface{}
	lock     sync.Mutex
}

// Type : event types
func (d *TimerParamFunc) Type() TimerType {
	return TypeParamFunc
}

// Next :
func (d *TimerParamFunc) Next() int64 {
	return int64(d.next)
}

// Run : run function
func (d *TimerParamFunc) Run() {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch d.times {
	case 0:
		return
	case 1:
		d.function(d.param)
		return
	case -1:
		// forever
	default:
		atomic.AddInt32(&d.times, -1)
	}

	d.function(d.param)
	atomic.AddInt64(&d.next, d.interval)
	putPools(d)
}

// Del :
func (d *TimerParamFunc) Del() {
	d.times, d.next = 0, GetNowStamp()
}
