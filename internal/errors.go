package internal

import "errors"

var UserNotFound = errors.New("user not found")
var UserAlreadyExists = errors.New("user already exists")
var RefreshTokenNotFound = errors.New("refresh token not found")
var RefreshTokenExpired = errors.New("refresh token expired")
