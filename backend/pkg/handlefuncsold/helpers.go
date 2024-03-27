package handlefuncs

import "golang.org/x/crypto/bcrypt"

func HashPassord(password string) []byte {
	p, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return p
}

func isPasswordValid(h, e []byte) error {
	return bcrypt.CompareHashAndPassword(h, e)
}
