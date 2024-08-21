package services

import (
	"sort"

	"github.com/JoshElias/go-web-server/internal"
)

func GetChirps(options internal.ChirpQueryOptions) ([]internal.ChirpEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return nil, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return nil, err
	}
	chirps := make([]internal.ChirpEntity, 0)
	for _, chirp := range db.Chirps {
		if options.AuthorId != 0 && chirp.AuthorId != options.AuthorId {
			continue
		}
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
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
