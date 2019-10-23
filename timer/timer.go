package timer

import (
	"errors"
	"strconv"
	"time"
)

// NewTimer ：make new ticker function
// 时间定时 stamp: 15:04:05; 04:05; 05;
// 时间间隔 stamp: s-m-h-d:  10s; 30m; 60h; 7d;
// 执行次数 times: run times (defalut -1:forever)
// 立即执行 run:   defalut: running next time
func NewTimer(stamp string, times int, run bool, msg interface{}, function func(interface{})) error {
	if next, interval, err := checkTime(stamp); err != nil {
		return err
	} else {
		if times < 0 {
			times = -1
		} else if times == 0 {
			return errors.New("ticker run times can not be zero")
		}

		if run {
			switch times {
			case 1:
				function(msg)
				return nil
			case -1:
				function(msg)
			default:
				times--
				function(msg)
			}
		}

		putInto(&TimerFunc{
			function: function,
			times:    times,
			next:     next,
			interval: interval,
			msg:      msg,
		})
	}
	return nil
}

// NewRunDuration : Make a new function run just only one times
func NewRunDuration(duration time.Duration, msg interface{}, function func(interface{})) {
	var data = &TimerFunc{
		next:     time.Now().Add(duration).UnixNano(),
		times:    1,
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// NewRunTime : Make a new function run time
func NewRunTime(timestamp time.Time, msg interface{}, function func(interface{})) {
	var data = &TimerFunc{
		next:     timestamp.UnixNano(),
		times:    1,
		function: function,
		msg:      msg,
	}
	putInto(data)
}

// check timerstamp value
func checkTime(stamp string) (int64, int64, error) {
	var (
		err      error
		temp     int
		next     int64 = time.Now().UnixNano()
		interval int64
	)

	switch stamp[len(stamp)-1:] {
	case "s", "S":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * SecondTimeUnit)

	case "m", "M":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * MinuteTimeUnit)

	case "h", "H":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * HourTimeUnit)

	case "d", "D":
		temp, err = strconv.Atoi(stamp[:len(stamp)-1])
		interval = int64(temp * DayTimeUnit)

	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":

		timeString := time.Now().Format(TimeFormatString)
		timeLen := len(timeString)
		var t time.Time

		switch len(stamp) {
		case 2: // second
			t, err = time.Parse(TimeFormatString, timeString[:timeLen-2]+stamp)
			next = t.UnixNano()
			interval = MinuteTimeUnit

		case 5: // min
			t, err = time.Parse(TimeFormatString, timeString[:timeLen-5]+stamp)
			next = t.UnixNano()
			interval = HourTimeUnit

		case 8: // hour
			t, err = time.Parse(TimeFormatString, timeString[:timeLen-8]+stamp)
			next = t.UnixNano()
			interval = DayTimeUnit

		default:
			err = errors.New("Can't parst time, please check it")
		}

	default:
		err = errors.New("Can't parst stamp value, please check it")
	}

	return next, interval, err
}
