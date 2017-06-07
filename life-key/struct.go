package lifekey

import "sync"

type LifeData struct {
	data map[string]*lifeKey
	mc   sync.RWMutex
}

type lifeKey struct {
	Data  interface{}
	Start int64
	Life  int64
}
