package internal

type UserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEntity struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
