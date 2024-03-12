package dbfuncs

import "log"

func GetGroupMembers(groupId string) []string {
	var result []string
	dbLock.RLock()
	rows, err := db.Query("SELECT * FROM GroupMembers WHERE GroupId = ?", groupId)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		rows.Close()
		dbLock.RUnlock()
	}()
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, userId)
	}
	return result
}
