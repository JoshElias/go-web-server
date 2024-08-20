package services

import (
	"github.com/JoshElias/go-web-server/internal"
)

func GetChirps() ([]internal.ChirpEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return nil, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return nil, err
	}
	chirpLen := len(db.Chirps)
	chirps := make([]internal.ChirpEntity, chirpLen)
	idx := 0
	for _, chirp := range db.Chirps {
		chirps[idx] = chirp
		idx++
	}
	return chirps, nil
}

func GetChirpById(id int) (internal.ChirpEntity, error) {
	if id < 1 {
		return internal.ChirpEntity{}, internal.ChirpNotFound
	}
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.ChirpEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.ChirpEntity{}, err
	}
	chirp, exists := db.Chirps[id]
	if !exists {
		return internal.ChirpEntity{}, internal.ChirpNotFound

	}
	return chirp, nil

}

func CreateChirp(authorId int, message string) (internal.ChirpEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.ChirpEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.ChirpEntity{}, err
	}
	id := len(db.Chirps) + 1
	newEntity := internal.ChirpEntity{
		Id:       id,
		Body:     message,
		AuthorId: authorId,
	}
	db.Chirps[id] = newEntity
	err = conn.WriteDb(db)
	if err != nil {
		return internal.ChirpEntity{}, nil
	}
	return newEntity, nil
}

func DeleteChirpById(id int) (bool, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return false, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return false, err
	}
	_, exists := db.Chirps[id]
	if !exists {
		return false, nil
	}
	delete(db.Chirps, id)
	return true, nil
}
