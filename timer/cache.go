package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

const basicNumber int64 = 100 // 100ms

// Container : save cache data containers
type Container struct {
	cache []TimerStruct
	count int32
	cutNo int64
	lock  sync.Mutex
}

var (
	first  = &Container{cutNo: 5 * basicNumber}
	second = &Container{cutNo: 25 * basicNumber}
	third  = &Container{cutNo: 250 * basicNumber}
)

var level int64 = 5
var initstatus uint32 = 0
var inittest sync.Once

// InitTickerInterval : init timer ticker
// base interval is 100ms, dafualt 100ms, [100ms * interval]
func InitTickerInterval(interval int64) {
	if atomic.AddUint32(&initstatus, 1) != 1 {
		return
	}
	if interval > 0 {
		level = interval
	}

	go func() {
		ticker := time.NewTicker(time.Nanosecond * time.Duration(basicNumber*level))

		for {
			select {
			case <-ticker.C:

				if count := atomic.AddInt32(&second.count, 1); count%4 == 0 {
					// run cache second check
					go checkSecondCache()
				}

				// append first arrge data
				first.lock.Lock()
				mins, maxs := TimeSplit(first.cache, GetNowStamp())
				first.cache = maxs
				first.lock.Unlock()

				for _, row := range mins {
					go row.Run()
				}
			}
		}
	}()
}

var minticker bool
var newpoint chan TimerStruct
var mintimer int64 = 10
var nexttime chan int64

// InitTickerMinTime : init timer ticker
// base interval is 1ms, dafualt 100ms, [100ms * interval]
func InitTickerMinTime(interval int64) {
	if atomic.AddUint32(&initstatus, 1) != 1 {
		return
	}
	minticker = true
	level = 5
	if interval > 0 {
		mintimer = interval
	}

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				// run cache second check
				go checkSecondCache()

			case ts := <-nexttime:

				go func(timestamp int64) {
					time.Sleep(time.Millisecond * time.Duration(timestamp))

					now := GetNowStamp() + mintimer
					// append first arrge data
					first.lock.Lock()
					mins, maxs := TimeSplit(first.cache, now)
					first.cache = maxs
					first.lock.Unlock()

					for _, row := range mins {
						go row.Run()
					}
				}(ts)

			case tmp := <-newpoint:
				// TODO: put to pool, check next run time
				tmp.Run()
			}
		}
	}()
}

func putPools(data TimerStruct) {
	if minticker {
		newpoint <- data
		return
	}

	now := GetNowStamp()

	if data.Next() <= (now + first.cutNo*level) {
		first.lock.Lock()
		first.cache = append(first.cache, data)
		first.lock.Unlock()
	} else if data.Next() > (now + second.cutNo*level) {
		third.lock.Lock()
		third.cache = append(third.cache, data)
		third.lock.Unlock()
	} else {
		second.lock.Lock()
		second.cache = append(second.cache, data)
		second.lock.Unlock()
	}
}

// When insert new data to sort
func checkSecondCache() {
	if count := atomic.AddInt32(&third.count, 1); count%5 == 0 {
		// run cache third check
		go checkThirdCache()
	}

	next := GetNowStamp() + first.cutNo*level

	second.lock.Lock()
	mins, maxs := TimeSplit(second.cache, next)

	first.lock.Lock()
	first.cache = append(first.cache, mins...)
	first.lock.Unlock()

	second.cache = maxs
	second.lock.Unlock()
}

// When insert new data to sort
func checkThirdCache() {
	next := GetNowStamp() + second.cutNo*level

	third.lock.Lock()
	mins, maxs := TimeSplit(third.cache, next)

	second.lock.Lock()
	second.cache = append(second.cache, mins...)
	second.lock.Unlock()

	third.cache = maxs
	third.lock.Unlock()
}
