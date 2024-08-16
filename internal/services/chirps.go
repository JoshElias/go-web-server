package services

import "github.com/JoshElias/chirpy/internal"

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

func CreateChirp(message string) (internal.ChirpEntity, error) {
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
		Id:   id,
		Body: message,
	}
	db.Chirps[id] = newEntity
	err = conn.WriteDb(db)
	if err != nil {
		return internal.ChirpEntity{}, nil
	}
	return newEntity, nil
}
