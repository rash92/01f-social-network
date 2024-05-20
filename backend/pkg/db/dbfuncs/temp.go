package dbfuncs

import (
	"fmt"

)

// type BasicUserInfo struct {
// 	Avatar string `json:"Avatar"`
// 	UserId         string `json:"UserId"`
// 	FirstName      string `json:"FirstName"`
// 	LastName       string `json:"LastName"`
// 	Nickname       string `json:"Nickname"`
// 	PrivacySetting string `json:"PrivacySetting"`
// }



// get number of x not generalised to avoid sql injection
func GetNumberOfById(userId string, table string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE Creatorid=?", table)
	err := db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}

// function seperated to avoid sql injection
func GetNumberOfFollowersAndFollowing(flag string, ownerId string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM Follows WHERE %s=?", flag)
	err := db.QueryRow(query, ownerId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}



