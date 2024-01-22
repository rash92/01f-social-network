package dbfuncs

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// moved
func CheckValueInDB(w http.ResponseWriter, r *http.Request, val, name string) (string, bool, error) {
	if name != "Nickname" && name != "Email" {
		return "Invalid column name", false, nil
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s = ?", name)
	err := database.QueryRow(query, val).Scan(&count)
	if err != nil {
		return "Error querying database", false, err
	}
	fmt.Println(count > 0, name)
	return "", count > 0, nil
}
