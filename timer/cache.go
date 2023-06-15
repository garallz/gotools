package timer

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type timerData struct {
	zone  time.Duration // time zone of second
	uniq  int64         // unique id
	root  *chainNode
	last  *chainNode
	timer *time.Ticker
}

func (t *timerData) SetTimeZone(dur time.Duration) { t.zone = dur }

func (t *timerData) TimerRun() {
	for range t.timer.C {
		now := time.Now().UnixMilli()
		t.root.lock.Lock()

		for i, data := range t.root.data {
			if data == nil || data.Next() > now {
				continue
			}
			go func(tmp runFunc) {
				if tmp.Run() {
					t.last.AddFunc(tmp.Next()-now, tmp)
				}
			}(data)
			t.root.data[i] = nil
		}

		t.root.lock.Unlock()
		t.root.next.CheckFunc()
	}
}

// delete timer function by unique id
func (t *timerData) DeleteTimer(id int64) {
	var tmp = t.root
	for tmp != nil {
		for i := range tmp.data {
			if tmp.data[i] != nil && tmp.data[i].Uuid() == id {
				tmp.data[i] = nil
				return
			}
		}
		tmp = tmp.next
	}
}

type chainNode struct {
	last *chainNode
	next *chainNode
	lock sync.Mutex
	data []runFunc

	serial int
	number int32 // check func number
	nummax int32 // check func trigger number
	durmin int64 // (ms) node limit
	durmax int64 // (ms) node limit
}

// Insert upside down (last <-)
func (c *chainNode) AddFunc(sleep int64, data runFunc) {
	if c == nil || sleep <= 0 {
		log.Panicf("abnormal timing task, sleep:%d, data:%v", sleep, data)
	}
	if (c.durmax == 0 || c.durmax >= sleep) && sleep >= c.durmin {
		c.lock.Lock()
		defer c.lock.Unlock()

		for i := range c.data {
			if c.data[i] == nil {
				c.data[i] = data
				return
			}
		}
		c.data = append(c.data, data)
		return
	} else if c.last != nil && c.last.durmax >= sleep {
		c.last.AddFunc(sleep, data)
		return
	}
	log.Panicln("add function abnormal", c.durmax, c.durmin, sleep, c.last)
}

// check function (-> next)
func (c *chainNode) CheckFunc() {
	// trigger interval check function with next time
	if c != nil && atomic.AddInt32(&c.number, 1) >= c.nummax {
		if c.durmin > 0 && c.last != nil {
			var now = time.Now().UnixMilli()
			var sleep = now + c.durmin

			c.lock.Lock()
			for i, v := range c.data {
				if v != nil && v.Next() < sleep {
					go c.last.AddFunc(v.Next()-now, v)
					c.data[i] = nil
				}
			}
			c.lock.Unlock()
		}
		atomic.SwapInt32(&c.number, 0)
		if c.next != nil {
			c.next.CheckFunc()
		}
	}
}
