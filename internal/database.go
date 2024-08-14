package internal

import "sync"

type DB struct {
	path string
	mux  *sync.RWMutex
}
