package timer

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const basicNumber int64 = 100 * MilliTimeUnit // 100ms

var (
	firstCut  int64 = 5 * basicNumber
	secondCut int64 = 25 * basicNumber
	thirdCut  int64 = 250 * basicNumber
)

var (
	cacheFirst  = make([]*TimerFunc, 0)
	cacheSecond = make([]*TimerFunc, 0)
	cacheThird  = make([]*TimerFunc, 0)
	bufferCache = make([]*TimerFunc, 0)
)

var (
	secodeCount int32 = 0
	thirdCount  int32 = 0
)

var (
	firstLock  sync.Mutex
	secondLock sync.Mutex
	thirdLock  sync.Mutex
	bufferLock sync.Mutex
)

var level int64 = 5
var runStatus bool

// InitTicker : init timer ticker
// base interval is 100ms, dafualt 100ms, [100ms * interval]
func InitTicker(interval int64) {
	if interval > 0 {
		level = interval
	}

	go func() {
		ticker := time.NewTicker(time.Nanosecond * time.Duration(basicNumber*level))

		for {
			select {
			case <-ticker.C:

				if count := atomic.AddInt32(&secodeCount, 1); count >= 4 {
					// run cache second check
					go checkSecondCache()
					atomic.SwapInt32(&secodeCount, 0)
				}

				if runStatus {
					continue
				}
				runStatus = true

				firstLock.Lock()

				// append buffer arrge data
				bufferLock.Lock()
				cacheFirst = append(cacheFirst, bufferCache...)
				bufferCache = make([]*TimerFunc, 0)
				bufferLock.Unlock()
				// sort array data
				sort.Sort(ByNext(cacheFirst))

				now := time.Now().UnixNano()
				length := len(cacheFirst) - 1
				for i, row := range cacheFirst {
					if row.next <= now {
						go run(row)
					} else {
						cacheFirst = cacheFirst[i:]
						break
					}
					if i == length {
						cacheFirst = make([]*TimerFunc, 0)
					}
				}
				runStatus = false

				firstLock.Unlock()
			}
		}
	}()
}

func run(data *TimerFunc) {
	switch data.times {
	case 0:
		return
	case 1:
		data.function(data.msg)
		return
	case -1:
		data.next += data.interval
	default:
		data.times--
		data.next += data.interval
	}

	go putInto(data)
	data.function(data.msg)
}

func putInto(data *TimerFunc) {
	now := time.Now().UnixNano()

	if data.next <= (now + firstCut*level) {
		bufferLock.Lock()
		bufferCache = append(bufferCache, data)
		bufferLock.Unlock()
	} else if data.next > (now + thirdCut*level) {
		secondLock.Lock()
		cacheThird = append(cacheThird, data)
		secondLock.Unlock()
	} else {
		thirdLock.Lock()
		cacheSecond = append(cacheSecond, data)
		thirdLock.Unlock()
	}
}

// When insert new data to sort
func checkSecondCache() {
	var next = time.Now().UnixNano() + firstCut*level
	// Sort arrge
	sort.Sort(ByNext(cacheSecond))

	length := len(cacheSecond) - 1
	for i, row := range cacheSecond {
		if row.next < next && i < length {
			continue
		}

		bufferLock.Lock()
		bufferCache = append(bufferCache, cacheSecond[:i+1]...)
		bufferLock.Unlock()

		secondLock.Lock()
		cacheSecond = cacheSecond[i+1:]
		secondLock.Unlock()
		return
	}

	if count := atomic.AddInt32(&thirdCount, 1); count >= 10 {
		// run cache third check
		go checkThirdCache()
		atomic.SwapInt32(&thirdCount, 0)
	}
}

// When insert new data to sort
func checkThirdCache() {
	var next = time.Now().UnixNano() + secondCut*level
	// Sort arrge
	sort.Sort(ByNext(cacheThird))

	length := len(cacheThird) - 1
	for i, row := range cacheThird {
		if row.next < next && i < length {
			continue
		}

		secondLock.Lock()
		cacheSecond = append(cacheSecond, cacheThird[:i+1]...)
		secondLock.Unlock()

		thirdLock.Lock()
		cacheThird = cacheThird[i+1:]
		thirdLock.Unlock()
		return

	}
}
