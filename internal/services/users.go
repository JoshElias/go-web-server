package services

import (
	"errors"

	"github.com/JoshElias/go-web-server/internal"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(dto internal.UserDto) (internal.UserEntity, error) {
	if exists, err := IsUniqueUserEmail(dto.Email); err != nil {
		return internal.UserEntity{}, err
	} else if exists {
		return internal.UserEntity{}, internal.UserAlreadyExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)
	if err != nil {
		return internal.UserEntity{}, err
	}
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
		Email:    dto.Email,
		Password: passHash,
	}
	db.Users[id] = newUser
	err = conn.WriteDb(db)
	if err != nil {
		return internal.UserEntity{}, nil
	}
	return newUser, nil
}

func GetUserById(id int) (internal.UserEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.UserEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.UserEntity{}, err
	}
	user, exists := db.Users[id]
	if !exists {
		return internal.UserEntity{}, internal.UserNotFound
	}
	return user, nil
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
	return internal.UserEntity{}, internal.UserNotFound
}

func UpdateUserById(id int, patch internal.UserEntity) (internal.UserEntity, error) {
	conn, err := internal.GetTestDbConnection()
	if err != nil {
		return internal.UserEntity{}, err
	}
	db, err := conn.LoadDb()
	if err != nil {
		return internal.UserEntity{}, err
	}
	user, exists := db.Users[id]
	if !exists {
		return internal.UserEntity{}, internal.UserNotFound
	}
	user.Email = patch.Email
	user.IsChirpyRed = patch.IsChirpyRed
	// passHash, err := bcrypt.GenerateFromPassword([]byte(patch.Password), 12)
	// if err != nil {
	// 	return internal.UserEntity{}, err
	// }
	// user.Password = passHash
	db.Users[id] = user
	err = conn.WriteDb(db)
	if err != nil {
		return internal.UserEntity{}, nil
	}
	return user, nil
}

func IsUniqueUserEmail(email string) (bool, error) {
	_, err := GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, internal.UserNotFound) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}
