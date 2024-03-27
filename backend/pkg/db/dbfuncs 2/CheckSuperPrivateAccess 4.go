package dbfuncs

import "database/sql"

func CheckSuperprivateAccess(postId string, userId string) (bool, error) {
	var exists int
	err := database.QueryRow("SELECT 1 FROM PostAccess WHERE PostId = ? AND UserId = ?", postId, userId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
