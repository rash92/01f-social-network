package dbfuncs

func AddFollow(follow *Follow) error {
	statement, err := db.Prepare("INSERT INTO Comments VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(follow.FollowerId, follow.FollowingId, follow.Status)
	return err
}
