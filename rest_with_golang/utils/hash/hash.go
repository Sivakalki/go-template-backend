package hash


import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt = hash password (not reversible)
func Encrypt(password string) (string, error) {
	// bcrypt.DefaultCost = 10, can adjust
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare = check plain password with hashed password
func Compare(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
