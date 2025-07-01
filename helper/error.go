package helper

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
