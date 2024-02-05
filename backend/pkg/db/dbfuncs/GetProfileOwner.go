package dbfuncs

func GetProfileOwner(userId string) (User, error) {
	profile := User{}
	err := database.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&profile)
	if err != nil {
		return profile, err
	}
	return profile, nil
}
