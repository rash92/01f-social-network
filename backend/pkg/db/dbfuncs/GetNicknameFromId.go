package dbfuncs

func GetNicknameFromId(userId string) string {
	var nickname string
	err := db.QueryRow("SELECT Nickname FROM users WHERE id = ?", userId).Scan(&nickname)
	if err != nil {
		return ""
	}
	return nickname
}
