package internal

type UserDto struct {
	Email string `json:"email"`
}

type UserEntity struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
