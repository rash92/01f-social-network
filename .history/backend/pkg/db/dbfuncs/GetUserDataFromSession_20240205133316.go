package dbfuncs

func GetUserProfile(sessionId string) (string, string, string, error) {
	var userId string
	var profileImage string
	var nickname string
	var username string

	// Execute the SQL query
	err := database.QueryRow(`
			SELECT Sessions.userId, Users.Profile, Users.Nickname, Users.Username
			FROM Sessions
			JOIN Users ON Sessions.UserID = Users.Id
			WHERE Sessions.Id = ?
	`, sessionId).Scan(&userId, &profileImage, &nickname, &username)
	
	// Check for errors
	if err != nil {
			return "", "", "", err
	}

	return profileImage, nickname, username, nil
}