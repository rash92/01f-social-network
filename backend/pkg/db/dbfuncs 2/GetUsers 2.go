package dbfuncs

func Getusers() ([]User_getAllUsers, error) {
	rows, err := database.Query("SELECT Id,FirstName, LastName, Nickname, Profile, AboutMe, Privacy_setting, DOB, CreatedAt FROM Users")
	if err != nil {
		return []User_getAllUsers{}, err
	}
	defer rows.Close()
	var user []User_getAllUsers

	for rows.Next() {
		var newUser User_getAllUsers
		err := rows.Scan(&newUser.Id, &newUser.FirstName, &newUser.LastName, &newUser.Nickname, &newUser.Profile, &newUser.AboutMe, &newUser.Privacy_setting, &newUser.DOB, &newUser.CreatedAt)
		if err != nil {
			return []User_getAllUsers{}, err
		}
		user = append(user, newUser)
	}

	return user, err
}
