package dbfuncs

import (
	"database/sql"
)

func Getusers() ([]User_getAllUsers, error) {
	rows, err := database.Query("SELECT Id,FirstName, LastName, Nickname, Profile, AboutMe, Privacy_settin, DO  FROM User")
	if err != nil {
		return []User_getAllUsers{}, err
	}
	defer rows.Close()
	var user []User_getAllUsers

	for rows.Next() {
		var newUser User_getAllUsers
		err := rows.Scan(&newUser)
		if err != nil {
			return []User_getAllUsers{}, err
		}
		user = append(user, newUser)
	}

	return user, err
}
