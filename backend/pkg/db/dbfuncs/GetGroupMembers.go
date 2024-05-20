package dbfuncs

import "log"

func GetGroupMembers(groupId string) []string {
	var result []string
	lock.RLock()
	rows, err := database.Query("SELECT * FROM GroupMembers WHERE GroupId = ?", groupId)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		rows.Close()
		lock.RUnlock()
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
