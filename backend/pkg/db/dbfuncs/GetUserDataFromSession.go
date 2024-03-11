package dbfuncs

func GetUserDataFromSession(sessionId string) (string, string, string, error) {
	var userId string
	var profileImage string
	var nickname string

	// Execute the SQL query
	err := database.QueryRow(`
			SELECT Sessions.userId, Users.Profile, Users.Nickname
			FROM Sessions
			JOIN Users ON Sessions.UserID = Users.Id
			WHERE Sessions.Id = ?
	`, sessionId).Scan(&userId, &profileImage, &nickname)

	// Check for errors
	if err != nil {
		return "", "", "", err
	}

	return userId, profileImage, nickname, nil
}
