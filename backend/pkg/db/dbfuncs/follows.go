package dbfuncs

import (
	"database/sql"
)

func IsFollowing(userId, ownerId string) (bool, error) {
	isFollowing := false
	query := `SELECT EXISTS(SELECT 1 FROM Follows WHERE FollowerId=? AND FollowingId=? AND Status=?)`
	err := db.QueryRow(query, userId, ownerId, "accepted").Scan(&isFollowing)
	if err != nil {
		return false, err
	}
	return isFollowing, nil
}

func IsPending(userId, ownerId string) (bool, error) {
	isPending := false
	query := `SELECT EXISTS(SELECT 1 FROM Follows WHERE FollowerId=? AND FollowingId=? AND Status=?)`
	err := db.QueryRow(query, userId, ownerId, "pending").Scan(&isPending)
	if err != nil {
		return false, err
	}
	return isPending, nil
}

func AddFollow(follow *Follow) error {
	statement, err := db.Prepare("INSERT INTO Follows VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(follow.FollowerId, follow.FollowingId, follow.Status)
	return err
}

// Only to be used when updating a pending follow to accepted follow - only necessary if following private user
func AcceptFollow(followerId, followingId string) error {

	statement, err := db.Prepare("UPDATE Follows SET Status=?  WHERE FollowerId=? AND FollowingId=?")
	if err != nil {
		return err
	}

	_, err = statement.Exec("accepted", followerId, followingId)

	return err
}

// may not use this and can just delete follow from table instead on rejected a follow request
func RejectFollow(followerId, followingId string) error {
	statement, err := db.Prepare("UPDATE Follows SET Status=? WHERE FollowerId=? AND FollowingId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec("rejected", followerId, followingId)
	return err
}

// Delete follow from table when unfollowing
func DeleteFollow(followerId, followingId string) error {
	statement, err := db.Prepare("DELETE FROM Follows WHERE FollowerId=? AND FollowingId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(followerId, followingId)
	return err
}

func GetAcceptedFollowerIdsByFollowingId(followingId string) ([]string, error) {
	var followerIds []string
	rows, err := db.Query("SELECT FollowerId FROM Follows WHERE FollowingId=? AND Status=?", followingId, "accepted")
	if err == sql.ErrNoRows {
		return followerIds, nil
	}
	if err != nil {
		return followerIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var followerId string
		err := rows.Scan(&followerId)
		if err != nil {
			return followerIds, err
		}
		followerIds = append(followerIds, followerId)
	}
	err = rows.Err()
	return followerIds, err
}

// Find all people you are following (accepted follows only)
func GetAcceptedFollowingIdsByFollowerId(followerId string) ([]string, error) {
	var followingIds []string
	rows, err := db.Query("SELECT FollowingId FROM Follows WHERE FollowerId=? AND Status=?", followerId, "accepted")
	if err == sql.ErrNoRows {
		return followingIds, nil
	}
	if err != nil {
		return followingIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var followingId string
		err := rows.Scan(&followingId)
		if err != nil {
			return followingIds, err
		}
		followingIds = append(followingIds, followingId)
	}
	err = rows.Err()
	return followingIds, err
}

func GetPendingFollowerIdsByFollowingId(followingId string) ([]string, error) {
	var followerIds []string
	rows, err := db.Query("SELECT FollowerId FROM Follows WHERE FollowingId=? AND Status=?", followingId, "pending")
	if err == sql.ErrNoRows {
		return followerIds, nil
	}
	if err != nil {
		return followerIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var followerId string
		err := rows.Scan(&followerId)
		if err != nil {
			return followerIds, err
		}
		followerIds = append(followerIds, followerId)
	}
	err = rows.Err()
	return followerIds, err
}

// Find all people you are following (pending follows only)
func GetPendingFollowingIdsByFollowerId(followerId string) ([]string, error) {
	var followingIds []string
	rows, err := db.Query("SELECT FollowingId FROM Follows WHERE FollowerId=? AND Status=?", followerId, "pending")
	if err == sql.ErrNoRows {
		return followingIds, nil
	}
	if err != nil {
		return followingIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var followingId string
		err := rows.Scan(&followingId)
		if err != nil {
			return followingIds, err
		}
		followingIds = append(followingIds, followingId)
	}
	err = rows.Err()
	return followingIds, err
}

//TO DO: get 10 at a time? decide if doing it through SQL or get all and do in handlefunc

// rewrite?

// func GetNumberOfFollowersAndFollowing(flag string, ownerId string) (int, error) {
// 	var count int
// 	query := fmt.Sprintf("SELECT COUNT(*) FROM Follows WHERE %s=?", flag)
// 	err := database.QueryRow(query, ownerId).Scan(&count)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	return count, nil
// }

// func GetFollowersOrFollowing(ownerId string, itemId string, offset int) ([]string, error) {
// 	items := []string{}
// 	var oppositeId string
// 	if itemId == "FollowerId" {
// 		oppositeId = "FollowingId"
// 	} else {
// 		oppositeId = "FollowerId"
// 	}
// 	query := fmt.Sprintf("SELECT %s FROM Follows WHERE %s = ? LIMIT 10 OFFSET %d", itemId, oppositeId, offset)
// 	rows, err := db.Query(query, ownerId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var item string
// 		err := rows.Scan(&item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		items = append(items, GetNicknameFromId(item))
// 	}
// 	sort.Strings(items)
// 	return items, nil
// }
