package timer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// NewTimer ：make new ticker function
// stamp -> Time timing: 15:04:05; 04:05; 05;
// times: 	run times [-1:forever, 0:return not run]
func (s *timerData) AddTimerArgRun(stamp string, times int, arg interface{},
	function func(interface{})) (int64, error) {
	if stamp == "" || function == nil || times == 0 {
		return 0, errors.New("time stamp or function wrong")
	}
	if next, interval, err := s.checkTime(stamp); err != nil {
		return 0, err
	} else {
		fmt.Printf("- %d, %d, %v", next, interval, time.UnixMilli(next))

		var uniq = atomic.AddInt64(&s.uniq, 1)
		s.last.AddFunc(next-time.Now().UnixMilli(),
			&funcArgsBody{next: next, uniq: uniq, interval: interval,
				function: function, message: arg, times: times},
		)
		return uniq, nil
	}
}

func (s *timerData) AddTimerCustom(stamp string, times int,
	function func()) (int64, error) {
	if stamp == "" || function == nil || times == 0 {
		return 0, errors.New("time stamp or function wrong")
	}
	if next, interval, err := s.checkTime(stamp); err != nil {
		return 0, err
	} else {
		var uniq = atomic.AddInt64(&s.uniq, 1)
		s.last.AddFunc(next-time.Now().UnixMilli(),
			&funcOnlyBody{next: next, interval: interval, uniq: uniq,
				function: function, times: times})
		return uniq, nil
	}
}

// NewRunDuration : Make a new function run
// times: [-1 meas forever], [0 meas not run]
// if argument or function is nil, return 0 with not run
func (s *timerData) AddTimerArgument(duration time.Duration, times int,
	arg interface{}, function func(interface{})) int64 {
	if function == nil || times == 0 || arg == nil {
		return 0
	}
	var uniq = atomic.AddInt64(&s.uniq, 1)
	s.last.AddFunc(duration.Milliseconds(), &funcArgsBody{
		uniq:     uniq,
		next:     time.Now().Add(duration).UnixMilli(),
		interval: duration.Milliseconds(),
		function: function, message: arg, times: times,
	})
	return uniq
}

func (s *timerData) AddTimerFunction(duration time.Duration, times int, function func()) int64 {
	if function == nil || times == 0 {
		return 0
	}
	var uniq = atomic.AddInt64(&s.uniq, 1)
	s.last.AddFunc(duration.Milliseconds(), &funcOnlyBody{
		uniq:     uniq,
		next:     time.Now().Add(duration).UnixMilli(),
		interval: duration.Milliseconds(),
		function: function, times: times,
	})
	return uniq
}

func (s *timerData) AddTimerBoolean(duration time.Duration, function func() bool) int64 {
	var uniq = atomic.AddInt64(&s.uniq, 1)
	s.last.AddFunc(duration.Milliseconds(),
		&funcBoolBody{uniq: uniq,
			next:     time.Now().Add(duration).UnixMilli(),
			interval: duration.Milliseconds(),
			function: function, times: -1,
		})
	return uniq
}

func (s *timerData) AddTimerArgBool(duration time.Duration, arg interface{}, function func(interface{}) bool) int64 {
	var uniq = atomic.AddInt64(&s.uniq, 1)
	s.last.AddFunc(duration.Milliseconds(),
		&funcArgsBool{uniq: uniq,
			next:     time.Now().Add(duration).UnixMilli(),
			interval: duration.Milliseconds(),
			function: function, message: arg, times: -1,
		})
	return uniq
}

func (s *timerData) AddTimerMonth(day int, stamp string, times int, function func() bool) (int64, error) {
	if day < 0 || day > 28 || times == 0 {
		return 0, errors.New("timer month day scope was wrong or times is null")
	}
	var now = time.Now()
	var next time.Time
	if stamp == "" {
		next = time.Date(now.Year(), now.Month()-1, day, 0, 0, 0, 0, time.UTC).Add(s.zone)
	} else {
		stamp = strings.ReplaceAll(stamp, ":", "")
		timestamp, err := strconv.ParseInt(stamp, 10, 64)
		if err != nil || timestamp < 0 || timestamp > 240000 ||
			timestamp/10000 < 60000 && timestamp/100 < 60 {
			return 0, errors.New("parse timestamp wrong: " + err.Error())
		}
		next = time.Date(now.Year(), now.Month(), day, int(timestamp)/10000,
			int(timestamp)%10000/100, int(timestamp)%100, 0, time.UTC).Add(s.zone)
	}
	for time.Since(next).Seconds() > 0 {
		next = next.AddDate(0, 1, 0)
	}
	var uniq = atomic.AddInt64(&s.uniq, 1)
	s.last.AddFunc(time.Until(next).Milliseconds(), &funcMonthDay{
		uniq: uniq, next: next.UnixMilli(), function: function, times: times})
	return uniq, nil
}

func (s *timerData) CloseAndExit() { s.timer.Stop() }

// check timerstamp value, [ 15:04:05 | 04:05 | 05 ]
func (s *timerData) checkTime(stamp string) (int64, int64, error) {
	if len(stamp) < 2 {
		return 0, 0, errors.New("timestamp was wrong")
	}
	stamp = strings.ReplaceAll(stamp, ":", "")
	timestamp, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return 0, 0, errors.New("parse timestamp wrong: " + err.Error())
	}

	var now = time.Now()
	var interval int64
	var next time.Time
	if timestamp > 10000 && timestamp < 240000 &&
		timestamp/10000 < 60000 && timestamp/100 < 60 {
		interval = 24 * 60 * 60 * 1000
		next = time.Date(now.Year(), now.Month(), now.Day()-1, int(timestamp)/10000,
			int(timestamp)%10000/100, int(timestamp)%100, 0, time.UTC)
	} else if timestamp > 100 && timestamp < 6000 && timestamp/100 < 60 {
		interval = 60 * 60 * 1000
		next = time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(),
			int(timestamp)/100, int(timestamp)%100, 0, time.UTC)
	} else if timestamp < 60 {
		interval = 60 * 1000
		next = time.Date(now.Year(), now.Month(), now.Day()-1, now.Hour(),
			now.Minute(), int(timestamp), 0, time.UTC)
	} else {
		return 0, 0, errors.New("timestamp format was wrong")
	}
	next = next.Add(s.zone).Add(time.Hour * -24)
	for time.Since(next).Seconds() > 0 {
		next = next.Add(time.Millisecond * time.Duration(interval))
	}
	return next.UnixMilli(), interval, nil
}
