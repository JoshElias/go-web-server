package internal

import (
	"encoding/json"
	"errors"
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
	Chirps map[int]ChirpEntity `json:"chirps"`
	Users  map[int]UserEntity  `json:"users"`
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
			Chirps: make(map[int]ChirpEntity),
			Users:  make(map[int]UserEntity),
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
	db, err := conn.LoadDb()
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
	db, err := conn.LoadDb()
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

func (conn *DbConnection) CreateUser(email string, password []byte) (UserEntity, error) {
	db, err := conn.LoadDb()
	if err != nil {
		return UserEntity{}, err
	}
	id := len(db.Users) + 1
	newUser := UserEntity{
		Id:       id,
		Email:    email,
		Password: password,
	}
	db.Users[id] = newUser
	err = conn.writeDb(db)
	if err != nil {
		return UserEntity{}, nil
	}
	return newUser, nil
}

func (conn *DbConnection) GetUserByEmail(email string) (UserEntity, error) {
	db, err := conn.LoadDb()
	if err != nil {
		return UserEntity{}, err
	}
	for _, user := range db.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return UserEntity{}, errors.New("User not found")
}
