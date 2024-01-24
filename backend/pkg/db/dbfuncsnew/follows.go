package dbfuncs

func AddFollow(follow *Follow) error {
	statement, err := db.Prepare("INSERT INTO Comments VALUES (?,?)")
	if err != nil {
		return err
	}
	statement.Exec(follow.FollowerId, follow.FollowingId)
	return nil
}
