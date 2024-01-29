package dbfuncs

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AddCat(CatName string, CatDesc string) {
	id, _ := uuid.NewRandom()
	created := time.Now()
	statement, _ := database.Prepare("INSERT INTO Categories VALUES (?,?,?,?)")
	statement.Exec(id, CatName, CatDesc, created)

}
