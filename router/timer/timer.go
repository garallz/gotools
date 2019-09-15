package logfile

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"
)

var timerFunc []*TimerFunc
var accuracy int64 = 100 * 1000 * 1000 // 精确度 ms
var timerLock = new(sync.Mutex)

func init() {

	go func() {
		for {
			var temp int64 = time.Now().UnixNano()
			var next int64

			timerLock.Lock()

			for _, f := range timerFunc {

				if f.next < temp {
					go f.function()

					if f.fixed {
						f.next += f.interval
					} else {
						// 	f.next = time.Unix(0, next).AddDate(0, int(f.interval), 0).UnixNano()
					}

					if f.next < next || next == 0 {
						next = f.next
					}

					if f.times == 0 {
						go deleteFunc()
						continue
					} else if f.times > 0 {
						f.times--
					}
				}
			}
			timerLock.Unlock()

			if (next - temp) <= accuracy {
				time.Sleep(time.Nanosecond * time.Duration(accuracy))
			} else {
				time.Sleep(time.Nanosecond * time.Duration(next-time.Now().UnixNano()))
			}
		}
	}()
}

// 时间定时 stamp: 15:04:05; 04:05; 05;
// 时间间隔 stamp: s-m-h-d:  10s; 30m; 60h; 7d;
// 执行次数 times: run times (defalut -1:forever)
// 立即执行 run:   defalut: running next time
func NewTimer(stamp string, times int, run bool, function func()) error {
	if next, interval, err := checkTime(stamp); err != nil {
		return err
	} else {
		if times == 0 {
			times = -1
		}

		if run {
			if next > time.Now().UnixNano() {
				next -= interval
			}
		} else {
			if next < time.Now().UnixNano() {
				next += interval
			}
		}

		// and function and sort functions by interval
		timerLock.Lock()

		timerFunc = append(timerFunc, &TimerFunc{
			next:     next,
			interval: interval,
			times:    times,
			fixed:    true,
			function: function,
		})
		sort.Sort(ByInterval{timerFunc})

		timerLock.Unlock()
	}
	return nil
}

// default accuracy to run, min 10ms
func SetAccuracy(ms uint) error {
	if ms < 10 {
		return errors.New("Accuracy min is 10ms")
	}
	accuracy = int64(ms * 1000 * 1000)
	return nil
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

// delete function by times eq 0
func deleteFunc() {
	timerLock.Lock()
	for i, f := range timerFunc {
		if f.next == 0 {
			timerFunc = append(timerFunc[:i], timerFunc[i+1:]...)
		}
		return
	}
	timerLock.Unlock()
}
