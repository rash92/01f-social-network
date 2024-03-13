package dbfuncs

import (
	"database/sql"
)

func GetUserIdFromCokie(sessionId string) (string, error) {
	var userID string
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", sessionId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
