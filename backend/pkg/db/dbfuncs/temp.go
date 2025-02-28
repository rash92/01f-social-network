package dbfuncs

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
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
		nickname, err := GetNicknameFromId(item)
		if err != nil {
			return nil, err
		}
		items = append(items, nickname)
	}
	sort.Strings(items)
	return items, nil
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
		var groupId sql.NullString
		err := rows.Scan(&post.Id,
			&post.Title,
			&post.Body,
			&post.CreatorId,
			&groupId,
			&post.CreatedAt,
			&post.Image,
			&post.PrivacyLevel)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		post.GroupId = StringNull(groupId)
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

func SearchFollowers(query, ownerId string) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT Id, Nickname, Avatar FROM users WHERE Id IN (SELECT FollowerId FROM Follows WHERE FollowingId = ?) AND (FirstName LIKE ? OR LastName LIKE ? OR Nickname LIKE ?)", ownerId, "%"+query+"%", "%"+query+"%", "%"+query+"%")
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

func ConvertBase64ToImage(base64String string, directoryPath string) (string, error) {
	// Split the base64 string to isolate the MIME type and the actual data

	splitData := strings.Split(base64String, ",")
	if len(splitData) != 2 {
		return "", fmt.Errorf("nvalid base64 string")
	}

	mimeType := strings.Split(splitData[0], ";")[0]
	data := splitData[1]

	// Map the MIME type to a file extension
	mimeToExtension := map[string]string{
		"data:image/jpeg": ".jpg",
		"data:image/png":  ".png",
		"data:image/gif":  ".gif",
		// Add more mappings as needed
	}
	extension, ok := mimeToExtension[mimeType]
	if !ok {
		return "", fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Generate a UUID for the file name
	fileName := uuid.New().String() + extension

	// Decode the base64 string back to bytes
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Write the bytes to a new image file
	imagePath := filepath.Join(directoryPath, fileName)
	err = os.WriteFile(imagePath, decodedData, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
