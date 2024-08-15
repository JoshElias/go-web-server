package internal

import (
	"encoding/json"
	"fmt"
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
	conn.mux.Lock()
	defer conn.mux.Unlock()
	f, err := os.OpenFile(conn.path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return DbStructure{}, err
	}
	defer f.Close()
	byteArr, err := io.ReadAll(f)
	if err != nil {
		return DbStructure{}, err
	}
	if len(byteArr) == 0 {
		return DbStructure{
			Chirps: make(map[int]ChirpEntity),
		}, nil
	}
	var db DbStructure
	if err := json.Unmarshal(byteArr, &db); err != nil {
		fmt.Println(err)
		return DbStructure{}, err
	}
	return db, nil
}

func (conn *DbConnection) writeDb(db DbStructure) error {
	conn.mux.Lock()
	defer conn.mux.Unlock()
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

func (conn *DbConnection) GetChirps() ([]ChirpEntity, error) {
	db, err := conn.loadDb()
	if err != nil {
		return nil, err
	}
	chirpLen := len(db.Chirps)
	chirps := make([]ChirpEntity, chirpLen)
	idx := 0
	for _, chirp := range db.Chirps {
		chirps[idx] = chirp
		idx++
	}
	return chirps, nil
}

func (conn *DbConnection) CreateChirp(message string) (ChirpEntity, error) {
	db, err := conn.loadDb()
	if err != nil {
		return ChirpEntity{}, err
	}
	id := len(db.Chirps) + 1
	newEntity := ChirpEntity{
		Id:   id,
		Body: message,
	}
	db.Chirps[id] = newEntity
	err = conn.writeDb(db)
	if err != nil {
		return ChirpEntity{}, nil
	}
	return newEntity, nil
}
