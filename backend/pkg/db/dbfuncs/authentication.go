package dbfuncs

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func IsLoginValid(email, enteredPass string) (string, error) {
	var storedPassword string
	var userId string
	fmt.Println(email, "email")

	err := db.QueryRow("SELECT Password, Id FROM Users WHERE Email = ?", email).Scan(&storedPassword, &userId)
	if err != nil {
		fmt.Println(err, "err login")
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
