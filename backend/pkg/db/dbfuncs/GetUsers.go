package dbfuncs

func Getusers() ([]User, error) {
	rows, err := db.Query("SELECT Id,FirstName, LastName, Nickname, Avatar, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	var user []User

	for rows.Next() {
		var newUser User
		err := rows.Scan(&newUser.Id, &newUser.FirstName, &newUser.LastName, &newUser.Nickname, &newUser.Profile, &newUser.AboutMe, &newUser.PrivacySetting, &newUser.DOB, &newUser.CreatedAt)
		if err != nil {
			return []User{}, err
		}
		user = append(user, newUser)
	}

	return user, err
}
