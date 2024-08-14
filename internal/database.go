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
	Chirps map[int]ChirpEntity `json:"chirps"`
}

func NewDbConnection(path string) (*DbConnection, error) {
	return &DbConnection{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil
}

func (conn *DbConnection) loadDb() (DbStructure, error) {
	f, err := os.OpenFile(conn.path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return DbStructure{}, err
	}
	defer f.Close()
	byteArr, err := io.ReadAll(f)
	if err != nil {
		return DbStructure{}, err
	}
	var db DbStructure
	if err := json.Unmarshal(byteArr, &db); err != nil {
		return DbStructure{}, err
	}
	return db, nil
}

func (conn *DbConnection) writeDb(db DbStructure) error {
	f, err := os.OpenFile(conn.path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	bytes, err := json.Marshal(db)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}
