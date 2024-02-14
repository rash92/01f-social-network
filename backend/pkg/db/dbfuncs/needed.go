package dbfuncs

func IsPublic(id string) (bool, error) {
	var privacySetting string
	err := db.QueryRow("SELECT PrivacySetting FROM Users WHERE id = $1", id).Scan(&privacySetting)
	if err != nil {
		return false, err
	}
	return privacySetting == "public", nil
}

func AddFollower(follow *Follow) error {
	_, err := db.Exec("INSERT INTO followers (FollowerId, FollowingId) VALUES ($1, $2)", follow.FollowerId, follow.FollowingId)
	return err
}

// placeholder
func AddNotification(notification *Notification) error {
	var err error
	return err
}
