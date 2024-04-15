package dbfuncs

import (
	"fmt"
	"log"
	"sort"
)

// type BasicUserInfo struct {
// 	Avatar string `json:"Avatar"`
// 	UserId         string `json:"UserId"`
// 	FirstName      string `json:"FirstName"`
// 	LastName       string `json:"LastName"`
// 	Nickname       string `json:"Nickname"`
// 	PrivacySetting string `json:"PrivacySetting"`
// }

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

func GetNicknameFromId(userId string) string {
	var nickname string
	err := db.QueryRow("SELECT Nickname FROM users WHERE id = ?", userId).Scan(&nickname)
	if err != nil {
		return ""
	}
	return nickname
}

func GetPosts(userId string, page int, batchSize int, usersOwnProfile bool) ([]Post, error) {
	_ = page
	_ = batchSize
	// offset := (page - 1) * batchSize
	query := `
		SELECT * FROM Posts
		WHERE CreatorId = ? AND PrivacyLevel != 'superprivate'
		ORDER BY CreatedAt DESC
	`
	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println(err, "GetPost")
		return nil, err
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id,
			&post.Title,
			&post.Body,
			&post.CreatorId,
			&post.CreatedAt,
			&post.Image,
			&post.PrivacyLevel)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Skip superprivate posts if the viewer does not have access.
		// if !usersOwnProfile || post.PrivacyLevel == "superprivate" {
		// 	allowed, err := CheckSuperprivateAccess(post.Id, userId) // Changed from post.Id.String() now that dbfuncs.Post.Id is of type string.
		// 	if err != nil {
		// 		log.Println(err)
		// 		return nil, err
		// 	}
		// 	if !allowed {
		// 		continue
		// 	}
		// }
		posts = append(posts, post)
	}
	// // If some posts could not be displayed because they were superprivate,
	// // recursively call GetPosts till we have 10 posts.
	// if !usersOwnProfile {
	// 	if len(posts) < 10 {
	// 		shortfall := 10 - len(posts)
	// 		page++
	// 		morePosts, err := GetPosts(userId, page, shortfall, usersOwnProfile)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return nil, err
	// 		}
	// 		posts = append(posts, morePosts...)
	// 	}
	// }
	// Swap order to display latest posts at the bottom.
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts, nil
}

func GetNumberOfById(userId string, table string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE Creatorid=?", table)
	err := db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}
func GetNumberOfFollowersAndFollowing(flag string, ownerId string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM Follows WHERE %s=?", flag)
	err := db.QueryRow(query, ownerId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	return count, nil
}

// func DeleteUserByUsername(username string) error {
// 	stmt, err := db.Prepare("DELETE * FROM Follows")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = stmt.Exec()
// 	if err != nil {
// 		//  you will get an arror if the user is not in the database
// 		// fmt.Println("error in deleting user by username", err)
// 		return err
// 	}
// 	return nil
// }

func SearchUsers(query string) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT Id,  Nickname, Avatar  FROM users WHERE FirstName LIKE ? OR LastName LIKE ? OR Nickname LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Avatar); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
