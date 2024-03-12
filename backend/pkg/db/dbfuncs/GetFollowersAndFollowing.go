package dbfuncs

import (
	"fmt"
	"sort"
)

func GetFollowersOrFollowing(ownerId string, itemId string, offset int) ([]string, error) {
	items := []string{}
	var oppositeId string
	if itemId == "FollowerId" {
		oppositeId = "FollowingId"
	} else {
		oppositeId = "FollowerId"
	}
	query := fmt.Sprintf("SELECT %s FROM Follows WHERE %s = ? LIMIT 10 OFFSET %d", itemId, oppositeId, offset)
	rows, err := db.Query(query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, GetNicknameFromId(item))
	}
	sort.Strings(items)
	return items, nil
}
