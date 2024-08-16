package services

import (
	"errors"

	"github.com/JoshElias/chirpy/internal"
)

func CreateUser(email string, password []byte) (internal.UserEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.UserEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.UserEntity{}, err
	}
	id := len(db.Users) + 1
	newUser := internal.UserEntity{
		Id:       id,
		Email:    email,
		Password: password,
	}
	db.Users[id] = newUser
	err = conn.WriteDb(db)
	if err != nil {
		return internal.UserEntity{}, nil
	}
	return newUser, nil
}

func GetUserByEmail(email string) (internal.UserEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.UserEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.UserEntity{}, err
	}
	for _, user := range db.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return internal.UserEntity{}, errors.New("User not found")
}
