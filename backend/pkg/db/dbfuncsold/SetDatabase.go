package dbfuncs

import "database/sql"

var database *sql.DB

// moved
func SetDatabase(db *sql.DB) {
	database = db
}
