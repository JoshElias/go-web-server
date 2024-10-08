package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type DbConnection struct {
	path string
	file *os.File
	mux  *sync.RWMutex
}

type DbStructure struct {
	Chirps        map[int]ChirpEntity     `json:"chirps"`
	Users         map[int]UserEntity      `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

func init() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		if err := DeleteTestDb(); err != nil {
			panic(err)
		}
	}
}

var TEST_DATABASE_FILENAME = "database.json"

func GetTestDbPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, TEST_DATABASE_FILENAME), nil
}

func GetTestDbConnection() (*DbConnection, error) {
	path, err := GetTestDbPath()
	if err != nil {
		return nil, err
	}
	return NewDbConnection(path)
}

func DeleteTestDb() error {
	path, err := GetTestDbPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err = os.Remove(path); err != nil {
		return err
	}
	return nil
}

func NewDbConnection(path string) (*DbConnection, error) {
	return &DbConnection{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil
}

func (conn *DbConnection) LoadDb() (DbStructure, error) {
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
			Chirps:        make(map[int]ChirpEntity),
			Users:         make(map[int]UserEntity),
			RefreshTokens: make(map[string]RefreshToken),
		}, nil
	}
	var db DbStructure
	if err := json.Unmarshal(byteArr, &db); err != nil {
		fmt.Println("ccccc")
		fmt.Println(err)
		return DbStructure{}, err
	}
	return db, nil
}

func (conn *DbConnection) WriteDb(db DbStructure) error {
	conn.mux.Lock()
	defer conn.mux.Unlock()
	f, err := os.OpenFile(conn.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
