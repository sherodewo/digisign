package helpers

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

/*
|===============================================================================
|	This File For Generic Function Password Hashing on another file or struct
|	and use it wisely !
|===============================================================================
*/

func HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func CheckPassword(password string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plain))
	return err == nil
}
