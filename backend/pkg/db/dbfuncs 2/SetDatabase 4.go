package dbfuncs

import "database/sql"

var database *sql.DB

func SetDatabase(db *sql.DB) {
	database = db
}