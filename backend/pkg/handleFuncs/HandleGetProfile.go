package handlefuncs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

// I've included whatever structs I needed in this file. They can be replaced
// with the real ones when they're ready, or if anyone knows where they live
// now. Likewise helper functions, database functions etc.

// I just added this to get rid of the red line under *Image. I don't know
// what Image is really supposed to be.
// type Image []byte

// This is what will be returned by the handler.
type Profile struct {
	Owner     User
	Posts     []Post
	Followers []string
	Following []string
}

// Fields as in the database.
type User struct {
	Id             string    `json:"id"`
	Nickname       string    `json:"nickname"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	Email          string    `json:"email,omitempty"`
	Password       string    `json:"password,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
	AboutMe        string    `json:"aboutme,omitempty"`
	PrivacySetting string    `json:"privacySetting,omitempty"`
	DOB            string    `json:"age,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {

	var userId string
	var ownerId string
	var profile Profile
	var usersOwnProfile bool

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error decoding userId: %v", err.Error())
		fmt.Println(err.Error(), "60")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Owner, err = GetProfileOwner(ownerId)

	if err != nil {
		errorMessage := fmt.Sprintf("error getting profile owner: %v", err.Error())
		fmt.Println(err.Error(), "66")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Owner.Password = ""

	cookie, _ := r.Cookie("user_token")
	userId, _ = GetUserIdFromCookie(cookie.Value)

	if userId == ownerId {
		usersOwnProfile = true
	}
	fmt.Println(!usersOwnProfile, "handleProfile:line:78")

	// Check Follows table to see if there's a row with FollowerId = userId and FollowingId = ownerId.
	var isFollowing bool
	query := `SELECT EXISTS(SELECT 1 FROM Follows WHERE FollowerId=? AND FollowingId=?)`
	err = database.QueryRow(query, userId, ownerId).Scan(&isFollowing)
	if err != nil {
		fmt.Println("failed to execute query: %v", err)
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	if !usersOwnProfile && profile.Owner.PrivacySetting == "private" && !isFollowing {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(profile); err != nil {
			// Handle JSON encoding error
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			fmt.Println("Failed to encode JSON:", err)
			return
		}

		return
	}

	profile.Posts, err = GetPosts(userId, 1, 10, usersOwnProfile)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		fmt.Println(err.Error(), "90")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Followers, err = GetFollowersOrFollowing(ownerId, "FollowerId", 1)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting followers: %v", err.Error())
		fmt.Println(err.Error(), "97")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Following, err = GetFollowersOrFollowing(ownerId, "FollowingId", 1)
	if err != nil {
		fmt.Println(err.Error(), "105")
		errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// fmt.Println(profile, "profile")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)

}

func GetProfileOwner(userId string) (User, error) {
	profile := User{}
	err := database.QueryRow("SELECT * FROM Users WHERE id = ?", userId).Scan(&profile.Id, &profile.Nickname, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Password, &profile.Avatar, &profile.AboutMe, &profile.PrivacySetting, &profile.DOB, &profile.CreatedAt)

	if err != nil {
		return profile, err
	}

	// fmt.Println(profile, "128 GetProfileOwner")
	return profile, nil
}

func GetFollowersOrFollowing(ownerId string, itemId string, offset int) ([]string, error) {
	items := []string{}
	var oppositeId string
	if itemId == "FollowerId" {
		oppositeId = "FollowingId"
	} else {
		oppositeId = "FollowerId"
	}
	query := fmt.Sprintf("SELECT %s FROM Follows WHERE %s = ? LIMIT 10 OFFSET %d", itemId, oppositeId, offset)
	rows, err := database.Query(query, ownerId)
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
	err := database.QueryRow("SELECT Nickname FROM users WHERE id = ?", userId).Scan(&nickname)
	if err != nil {
		return ""
	}
	return nickname
}

func GetPosts(userId string, page int, batchSize int, usersOwnProfile bool) ([]Post, error) {
	offset := (page - 1) * batchSize
	query := `
    SELECT * FROM Posts
    WHERE CreatorId = ?
    ORDER BY CreatedAt DESC
    LIMIT ? OFFSET ?
`
	rows, err := database.Query(query, userId, batchSize, offset)
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
		if !usersOwnProfile || post.PrivacyLevel == "superprivate" {
			allowed, err := CheckSuperprivateAccess(post.Id.String(), userId)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if !allowed {
				continue
			}
		}

		posts = append(posts, post)
	}

	// If some posts could not be displayed because they were superprivate,
	// recursively call GetPosts till we have 10 posts.
	if !usersOwnProfile {
		if len(posts) < 10 {
			shortfall := 10 - len(posts)
			page++
			morePosts, err := GetPosts(userId, page, shortfall, usersOwnProfile)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			posts = append(posts, morePosts...)
		}
	}

	// Swap order to display latest posts at the bottom.
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	return posts, nil
}

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

func GetUserIdFromCookie(sessionId string) (string, error) {
	var userId string
	err := database.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", sessionId).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
