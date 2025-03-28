package dbfuncs

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func IsLoginValid(email, enteredPass string) (string, error) {
	var storedPassword string
	var userId string
	err := db.QueryRow("SELECT Password, Id FROM Users WHERE Email = ?", email).Scan(&storedPassword, &userId)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPass))
	if err != nil {
		return "", err
	}
	return userId, nil
}

func ValidateCookie(sessionId string) (bool, error) {
	var id string
	var expiration time.Time
	err := db.QueryRow("SELECT Id, Expires FROM Sessions WHERE Id=?", sessionId).Scan(&id, &expiration)

	if err == sql.ErrNoRows {

		return false, nil
	} else if err != nil {

		return false, err
	}

	return id == sessionId && time.Now().Before(expiration), nil
}

func CheckEmailInDB(email string) (bool, error) {
	found := ""
	err := db.QueryRow("SELECT Email FROM Users WHERE Email=?", email).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
