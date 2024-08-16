package internal

type UserLogin struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEntity struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserView struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ChirpDto struct {
	Body string `json:"body"`
}

type ChirpEntity struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
