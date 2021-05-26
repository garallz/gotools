package timer

import "sync"

type TimerChain struct {
	next  int64
	data  TimerStruct
	child []TimerStruct
	head  *TimerChain
	last  *TimerChain
}

var chain struct {
	data  TimerChain
	next  int64
	hlock sync.Mutex   // The chain head safe lock
	clock sync.RWMutex // The chain update lock
}

func chainPut(d TimerStruct) {
	chain.clock.RLock()
	if chain.next == 0 || chain.next > d.Next() {
		chain.hlock.Lock()

	}
}
