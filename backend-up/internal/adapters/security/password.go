package security

import "golang.org/x/crypto/bcrypt"

// HashPassword gera o hash bcrypt da senha em texto puro
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword confere se a senha digitada bate com o hash salvo
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
