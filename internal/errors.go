package internal

import "errors"

var UserNotFound = errors.New("user not found")
var UserAlreadyExists = errors.New("user already exists")
