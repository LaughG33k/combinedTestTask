package pkg

import (
	customerrors "github.com/LaughG33k/authServiceTestTask/iternal/errors"
	"github.com/LaughG33k/authServiceTestTask/iternal/model"
)

func ValidateLoginPassword(login, password string) error {
	if login == "" {
		return customerrors.EmptyLogin
	}

	if len(login) > 30 {
		return customerrors.LoginSoLong
	}

	if password == "" {
		return customerrors.EmptyPassword
	}

	if len(password) > 30 {
		return customerrors.PasswordSoLong
	}

	return nil
}

func ValidateRegistration(v model.RegistrationModel) error {

	if err := ValidateLoginPassword(v.Login, v.Password); err != nil {
		return err
	}

	if v.Name == "" {
		return customerrors.EmptyName
	}

	if len(v.Name) > 30 {
		return customerrors.NameSoLong
	}

	if v.Email == "" {
		return customerrors.EmptyEmail
	}

	if len(v.Email) > 256 {
		return customerrors.EmailSoLong
	}

	return nil
}
