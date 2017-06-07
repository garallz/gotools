package lifekey

import (
	"time"
	"sync"
	"sync/atomic"
)

/*
	1.The package is use second as a lifecycle unit.
	2.Key-data can null.
	3.Key type is string.
 */
type LifeData struct {
	data map[string]*LifeKey
	mc sync.RWMutex
}

func (d *LifeData)Set(key string,  life int64) {
	d.mc.Lock()
	defer d.mc.Unlock()
	d.data[key] = &LifeKey{
		Start:time.Now().Unix(),
		Life:life,
	}
}

func (d *LifeData)Check(key string) bool {
	d.mc.RLock()
	defer d.mc.RUnlock()
	if _, ok := d.data[key]; ok {
		return true
	} else {
		return false
	}
}

func (d *LifeData)Get(key string) interface{} {
	d.mc.RLock()
	defer d.mc.RUnlock()
	if data, ok := d.data[key]; ok {
		return data.Data
	} else {
		return nil
	}
}

func (d *LifeData)SetAddData(data interface{}, key string,  life int64) {
	d.mc.Lock()
	defer d.mc.Unlock()
	d.data[key] = &LifeKey{
		Data:data,
		Start:time.Now().Unix(),
		Life:life,
	}
}

func (d *LifeData)Delete(key string) {
	d.mc.Lock()
	defer d.mc.Unlock()
	delete(d.data, key)
}

// If data not update, input nil.
func (d *LifeData)UpdateData(data interface{}, key string) {
	d.mc.Lock()
	defer d.mc.Unlock()
	d.data[key].Start = time.Now().Unix()
	if data != nil {
		d.data[key].Data = data
	}
}

func (d *LifeData)GcData(life time.Duration) {
	d.data = make(map[string]*LifeKey)
	go func () {
		var check int32 = 0
		var key string
		var data *LifeKey
		for _ = range time.NewTicker(life).C {
			if check == 0 {
				atomic.AddInt32(&check, 1)
				cacheData := d.data
				var timestamp = time.Now().Unix()
				for key, data = range cacheData {
					if data.Start + data.Life < timestamp {
						d.Delete(key)
					}
				}
				atomic.AddInt32(&check, -1)
			} else {
				continue
			}
		}
	}()
}