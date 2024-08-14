package internal

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type DbConnection struct {
	path string
	file *os.File
	mux  *sync.RWMutex
}

type DbStructure struct {
	Chirps map[int]ChirpEntity
}

func NewDbConnection(path string) (*DbConnection, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return &DbConnection{
		path: path,
		file: f,
		mux:  &sync.RWMutex{},
	}, nil
}

func (conn *DbConnection) loadDb() (DbStructure, error) {
	byteArr, err := io.ReadAll(conn.file)
	if err != nil {
		return DbStructure{}, err
	}
	var db DbStructure
	if err := json.Unmarshal(byteArr, &db); err != nil {
		return DbStructure{}, err
	}
	return db, nil
}
