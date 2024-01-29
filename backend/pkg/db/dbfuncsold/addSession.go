package dbfuncs

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// moved
func AddSession(id uuid.UUID, user, UserID string, Expires time.Time) {
	statement, _ := database.Prepare("INSERT INTO sessions VALUES (?,?,?,?)")

	statement.Exec(id, user, Expires, UserID)
}
