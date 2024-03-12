package dbfuncs

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AddUser(nickName, firstName, lastName, Email, profile, aboutMe, privacy, DOB string, Password []byte) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO Users VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	statement.Exec(id, nickName, firstName, lastName, Email, Password, profile, aboutMe, privacy, DOB, created)

	return nil
}
