package memkeys

import (
	"sync"
	"sync/atomic"
	"time"
)

const basicNumber int64 = 100 * 1000 * 1000 // 100ms

// Container : save cache data containers
type Container struct {
	cache []*KeyValue
	count int32
	cutNo int64
	lock  sync.Mutex
}

type ExpireCache struct {
	interval int64
	first    *Container
	second   *Container
	third    *Container
}

func (d *Memory) newExp() {
	d.expCache = &ExpireCache{
		third:  &Container{cutNo: 250 * basicNumber},
		second: &Container{cutNo: 25 * basicNumber},
		first:  &Container{cutNo: 5 * basicNumber},
	}
}

// InitTicker : init timer ticker
// base interval is 100ms, dafualt 100ms, [100ms * interval]
func (m *Memory) initExpire() {
	var times int32
	go func(m *Memory) {
		ticker := time.NewTicker(time.Nanosecond * time.Duration(basicNumber*m.expCache.interval))

		for {
			select {
			case <-ticker.C:

				var d = m.expCache

				if count := atomic.AddInt32(&d.second.count, 1); count >= 4 {
					// run cache second check
					go d.upSecondCache()
					atomic.SwapInt32(&d.second.count, 0)
				}
				// sort array data
				now := time.Now().UnixNano()

				d.first.lock.Lock()
				rows, max := ExpireSplit(d.first.cache, now)
				d.first.cache = max
				d.first.lock.Unlock()

				for _, row := range rows {
					go m.del(row.key)
				}

				// check memory size
				if count := atomic.AddInt32(&times, 1); count >= 100 {
					go m.checkCacheSize()
					atomic.SwapInt32(&times, 0)
				}
			}
		}
	}(m)
}

func (d *ExpireCache) putInto(data *KeyValue) {
	now := time.Now().UnixNano()

	if data.expire <= (now + d.first.cutNo*d.interval) {
		d.first.lock.Lock()
		d.first.cache = append(d.first.cache, data)
		d.first.lock.Unlock()
	} else if data.expire > (now + d.second.cutNo*d.interval) {
		d.third.lock.Lock()
		d.third.cache = append(d.third.cache, data)
		d.third.lock.Unlock()
	} else {
		d.second.lock.Lock()
		d.second.cache = append(d.second.cache, data)
		d.second.lock.Unlock()
	}
}

// When insert new data to sort
func (d *ExpireCache) upSecondCache() {
	if count := atomic.AddInt32(&d.third.count, 1); count >= 5 {
		// run cache third check
		go d.upThirdCache()
		atomic.SwapInt32(&d.third.count, 0)
	}

	var next = time.Now().UnixNano() + d.first.cutNo*d.interval
	d.second.lock.Lock()
	// data split by expire time
	min, max := ExpireSplit(d.second.cache, next)

	d.first.lock.Lock()
	d.first.cache = append(d.first.cache, min...)
	d.first.lock.Unlock()

	d.second.cache = max
	d.second.lock.Unlock()
}

// When insert new data to sort
func (d *ExpireCache) upThirdCache() {
	var next = time.Now().UnixNano() + d.second.cutNo*d.interval
	// Sort arrge
	d.third.lock.Lock()
	// data split by expire time
	min, max := ExpireSplit(d.third.cache, next)

	d.second.lock.Lock()
	d.second.cache = append(d.second.cache, min...)
	d.second.lock.Unlock()

	d.third.cache = max
	d.third.lock.Unlock()
}
