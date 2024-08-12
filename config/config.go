package config

import (
	"sync"
)

type ApiConfig struct {
	Mu             sync.Mutex
	FileserverHits int
}
