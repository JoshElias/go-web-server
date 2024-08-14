package internal

import (
	"sync"
)

type Metrics struct {
	Mu             sync.Mutex
	FileserverHits int
}

var (
	metrics Metrics
	once    sync.Once
)

func GetMetrics() *Metrics {
	once.Do(func() {
		metrics = Metrics{
			Mu:             sync.Mutex{},
			FileserverHits: 0,
		}
	})
	return &metrics
}
