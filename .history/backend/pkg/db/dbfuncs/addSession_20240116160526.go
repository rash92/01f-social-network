package dbfuncs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AddSession(id uuid.UUID, user, UserID string, Expires time.Time) {
	statement, err := database.Prepare("INSERT INTO Sessions VALUES (?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error(), "line 14 dbfuncs/addsesion")
		panic(err)
		return
	}
	statement.Exec(id, user, Expires, UserID)
}