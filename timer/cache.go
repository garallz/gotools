package timer

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const basicNumber int64 = 100 * MilliTimeUnit // 100ms

// Container : save cache data containers
type Container struct {
	cache []*TimerFunc
	count int32
	cutNo int64
	lock  sync.Mutex
}

var (
	buffer = &Container{}
	first  = &Container{cutNo: 5 * basicNumber}
	second = &Container{cutNo: 25 * basicNumber}
	third  = &Container{cutNo: 250 * basicNumber}
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

				if count := atomic.AddInt32(&second.count, 1); count >= 4 {
					// run cache second check
					go checkSecondCache()
					atomic.SwapInt32(&second.count, 0)
				}

				if runStatus {
					continue
				}
				runStatus = true

				first.lock.Lock()

				// append buffer arrge data
				buffer.lock.Lock()
				first.cache = append(first.cache, buffer.cache...)
				buffer.cache = make([]*TimerFunc, 0)
				buffer.lock.Unlock()

				// sort array data
				sort.Sort(ByNext(first.cache))

				now := time.Now().UnixNano()
				length := len(first.cache) - 1
				for i, row := range first.cache {
					if row.next <= now {
						go run(row)
					} else {
						first.cache = first.cache[i:]
						break
					}
					if i == length {
						first.cache = make([]*TimerFunc, 0)
					}
				}
				first.lock.Unlock()
				runStatus = false
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

	if data.next <= (now + first.cutNo*level) {
		buffer.lock.Lock()
		buffer.cache = append(buffer.cache, data)
		buffer.lock.Unlock()
	} else if data.next > (now + second.cutNo*level) {
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
	if count := atomic.AddInt32(&third.count, 1); count >= 5 {
		// run cache third check
		go checkThirdCache()
		atomic.SwapInt32(&third.count, 0)
	}

	var next = time.Now().UnixNano() + first.cutNo*level
	second.lock.Lock()
	defer second.lock.Unlock()
	// Sort arrge
	sort.Sort(ByNext(second.cache))

	if length := len(second.cache); length == 0 {
		return
	} else if second.cache[0].next > next {
		return
	} else if second.cache[length-1].next < next {
		buffer.lock.Lock()
		buffer.cache = append(buffer.cache, second.cache...)
		buffer.lock.Unlock()
		second.cache = make([]*TimerFunc, 0)
	} else {
		for i, row := range second.cache {
			if row.next > next {
				buffer.lock.Lock()
				buffer.cache = append(buffer.cache, second.cache[:i]...)
				buffer.lock.Unlock()
				second.cache = second.cache[i:]
				break
			}
		}
	}
}

// When insert new data to sort
func checkThirdCache() {
	var next = time.Now().UnixNano() + second.cutNo*level
	// Sort arrge
	third.lock.Lock()
	defer third.lock.Unlock()
	sort.Sort(ByNext(third.cache))

	if length := len(third.cache); length == 0 {
		return
	} else if third.cache[0].next > next {
		return
	} else if third.cache[length-1].next < next {
		second.lock.Lock()
		second.cache = append(second.cache, third.cache...)
		second.lock.Unlock()
		third.cache = make([]*TimerFunc, 0)
	} else {
		for i, row := range third.cache {
			if row.next > next {
				second.lock.Lock()
				second.cache = append(second.cache, third.cache[:i]...)
				second.lock.Unlock()
				third.cache = third.cache[i:]
				break
			}
		}
	}
}
