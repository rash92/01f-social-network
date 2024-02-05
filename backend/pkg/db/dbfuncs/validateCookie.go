package dbfuncs

import (
	"time"
)

func ValidateCookie(cookieValue string) bool {

	var id string
	var expiration time.Time
	err := database.QueryRow("SELECT Id, expires  FROM Sessions WHERE Id=?", cookieValue).Scan(&id, &expiration)
	if err != nil {
		return false
	}

	return id == cookieValue && !(time.Now().After(expiration))

}
