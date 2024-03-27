package handlefuncs

import (
	_ "github.com/mattn/go-sqlite3"
)

func FindAllCats() ([]Categories, error) {
	rows, err := database.Query("SELECT * FROM Categories")
	if err != nil {
		return []Categories{}, err

	}
	defer rows.Close()
	var newName []Categories

	for rows.Next() {
		var oneCat Categories
		err := rows.Scan(&oneCat.Id, &oneCat.Name, &oneCat.Description, &oneCat.Created_at)
		if err != nil {
			return []Categories{}, err
		}
		newName = append(newName, oneCat)
	}

	return newName, err
}
