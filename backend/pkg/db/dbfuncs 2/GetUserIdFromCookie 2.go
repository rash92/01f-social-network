package dbfuncs

func GetUserIdFromCookie(sessionId string) (string, error) {
	var userId string
	err := database.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", sessionId).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
