package customerrors

import "errors"

var EmptyEmail error = errors.New("empty email")
var EmailSoLong error = errors.New("email is so long")
var EmptyName error = errors.New("empty name")
var NameSoLong error = errors.New("name is so long")
var EmptyLogin error = errors.New("emty login")
var LoginSoLong error = errors.New("login is so long")
var EmptyPassword error = errors.New("empty password")
var PasswordSoLong error = errors.New("password is so long")

var BadRequst error = errors.New("Bad request")
