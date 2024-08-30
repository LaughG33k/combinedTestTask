package customerrors

import "errors"

var UserNotFound error = errors.New("user not found")
var UserAlreadyExists error = errors.New("user already exists")
var RefreshNotFound error = errors.New("refresh not found")
