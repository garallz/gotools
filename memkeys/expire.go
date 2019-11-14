package memkeys

import (
	"sort"
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
	buffer   *Container
	first    *Container
	second   *Container
	third    *Container
}

func (d *Memory) newExp() {
	d.expCache = &ExpireCache{
		third:  &Container{cutNo: 250 * basicNumber},
		second: &Container{cutNo: 25 * basicNumber},
		buffer: &Container{},
		first:  &Container{cutNo: 5 * basicNumber},
	}
}

// InitTicker : init timer ticker
// base interval is 100ms, dafualt 100ms, [100ms * interval]
func (m *Memory) initExpire() {
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
					// check memory size
					go m.checkCacheSize()
				}

				d.first.lock.Lock()

				// append buffer arrge data
				d.buffer.lock.Lock()
				d.first.cache = append(d.first.cache, d.buffer.cache...)
				d.buffer.cache = make([]*KeyValue, 0)
				d.buffer.lock.Unlock()

				// sort array data
				sort.Sort(KVS(d.first.cache))

				now := time.Now().UnixNano()
				length := len(d.first.cache) - 1
				for i, row := range d.first.cache {
					if row.expire == 0 {
						// Not To Do
					} else if row.expire <= now {
						m.del(row.key)
					} else {
						d.first.cache = d.first.cache[i:]
						break
					}
					if i == length {
						d.first.cache = make([]*KeyValue, 0)
					}
				}
				d.first.lock.Unlock()
			}
		}
	}(m)
}

func (d *ExpireCache) putInto(data *KeyValue) {
	now := time.Now().UnixNano()

	if data.expire <= (now + d.first.cutNo*d.interval) {
		d.buffer.lock.Lock()
		d.buffer.cache = append(d.buffer.cache, data)
		d.buffer.lock.Unlock()
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
	defer d.second.lock.Unlock()
	// Sort arrge
	sort.Sort(KVS(d.second.cache))

	if length := len(d.second.cache); length == 0 {
		return
	} else if d.second.cache[0].expire > next {
		return
	} else if d.second.cache[length-1].expire < next {
		d.buffer.lock.Lock()
		d.buffer.cache = append(d.buffer.cache, d.second.cache...)
		d.buffer.lock.Unlock()
		d.second.cache = make([]*KeyValue, 0)
	} else {
		for i, row := range d.second.cache {
			if row.expire > next {
				d.buffer.lock.Lock()
				d.buffer.cache = append(d.buffer.cache, d.second.cache[:i]...)
				d.buffer.lock.Unlock()
				d.second.cache = d.second.cache[i:]
				break
			}
		}
	}
}

// When insert new data to sort
func (d *ExpireCache) upThirdCache() {
	var next = time.Now().UnixNano() + d.second.cutNo*d.interval
	// Sort arrge
	d.third.lock.Lock()
	defer d.third.lock.Unlock()
	sort.Sort(KVS(d.third.cache))

	if length := len(d.third.cache); length == 0 {
		return
	} else if d.third.cache[0].expire > next {
		return
	} else if d.third.cache[length-1].expire < next {
		d.second.lock.Lock()
		d.second.cache = append(d.second.cache, d.third.cache...)
		d.second.lock.Unlock()
		d.third.cache = make([]*KeyValue, 0)
	} else {
		for i, row := range d.third.cache {
			if row.expire > next {
				d.second.lock.Lock()
				d.second.cache = append(d.second.cache, d.third.cache[:i]...)
				d.second.lock.Unlock()
				d.third.cache = d.third.cache[i:]
				break
			}
		}
	}
}
