package dbfuncs

import (
	"database/sql"
	"time"
)

func ValidateCookie(cookieValue string) bool {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return false
	}
	defer db.Close()
	var id string
	var expiration time.Time
	err = db.QueryRow("SELECT Id, expires  FROM Sessions WHERE Id=?", cookieValue).Scan(&id, &expiration)
	if err != nil {
		return false
	}

	return id == cookieValue && !(time.Now().After(expiration))

}
