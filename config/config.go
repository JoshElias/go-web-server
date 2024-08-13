package config

import (
	"sync"
)

type ApiConfig struct {
	Mu             sync.Mutex
	FileserverHits int
}

var (
	config ApiConfig
	once   sync.Once
)

func GetConfig() *ApiConfig {
	once.Do(func() {
		config = ApiConfig{
			Mu:             sync.Mutex{},
			FileserverHits: 0,
		}
	})
	return &config
}
