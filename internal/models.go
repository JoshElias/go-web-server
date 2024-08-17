package internal

type UserLoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserLoginResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserView struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type UserEntity struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type ChirpDto struct {
	Body string `json:"body"`
}

type ChirpEntity struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}
