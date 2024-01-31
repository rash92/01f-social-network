package handlefuncs

import "golang.org/x/crypto/bcrypt"

// move to different file?
func HashPassword(password string) ([]byte, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func isPasswordValid(storedPass, enteredPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(enteredPass))
}
