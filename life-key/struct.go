package lifekey

import "sync"

type LifeData struct {
	data map[string]*LifeKey
	mc   sync.RWMutex
}

type LifeKey struct {
	Data  interface{}
	Start int64
	Life  int64
}
