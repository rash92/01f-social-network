package dbfuncs

import "fmt"

func GetNumberOfById(userId string, table string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE UserId=?", table)
	err := database.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}

func GetNumberOfFollowersAndFollowing(flag string, ownerId string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM Follows WHERE %s=?", flag)
	err := database.QueryRow(query, ownerId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}
