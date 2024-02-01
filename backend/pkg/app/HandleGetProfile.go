package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"
)

// I've included whatever structs I needed in this file. They can be replaced
// with the real ones when they're ready, or if anyone knows where they live
// now. Likewise helper functions, database functions etc.

// I just added this to get rid of the red line under *Image. I don't know
// what Image is really supposed to be.
type Image []byte

// This is what will be returned by the handler.
type Profile struct {
	Owner     User     `json:"user"`
	Posts     []Post   `json:"posts,omitempty"`
	Followers []string `json:"followers,omitempty"`
	Following []string `json:"following,omitempty"`
}

// Fields as in the database.
type User struct {
	Id             string    `json:"id"`
	Nickname       string    `json:"nickname"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	Email          string    `json:"email,omitempty"`
	Password       string    `json:"password,omitempty"`
	Avatar         *Image    `json:"avatar,omitempty"`
	Aboutme        string    `json:"aboutme,omitempty"`
	PrivacySetting string    `json:"privacySetting,omitempty"`
	DOB            string    `json:"age,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	var userId string
	var ownerId string
	var profile Profile
	var usersOwnProfile bool

	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error decoding userId: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}

	profile.Owner, err = GetProfileOwner(ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting profile owner: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}

	profile.Owner.Password = ""

	cookie, _ := r.Cookie("user_token")
	userId, err = GetUserIdFromCookie(cookie.Value)
	if userId == ownerId {
		usersOwnProfile = true
	}

	if !usersOwnProfile || profile.Owner.PrivacySetting == "private" {
		profile.Owner = User{
			Nickname: profile.Owner.Nickname,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(profile)
		return
	}

	profile.Posts, err = GetPosts(userId, 1, 10, usersOwnProfile)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Followers, err = GetFollowersOrFollowing(ownerId, "FollowerId", 1)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting followers: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Following, err = GetFollowersOrFollowing(ownerId, "FollowingId", 1)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func GetProfileOwner(userId string) (User, error) {
	profile := User{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&profile)
	if err != nil {
		return profile, err
	}
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
	query := fmt.Sprintf("SELECT %s FROM Follows WHERE %s = ?", itemId, oppositeId, ownerId)
	query += fmt.Sprintf(" LIMIT 10 OFFSET %d", offset)
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
	slices.Sort(items)
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
	offset := (page - 1) * batchSize
	query := `
    SELECT * FROM Posts
    WHERE CreatorId = ?
    ORDER BY Created DESC
    LIMIT ? OFFSET ?
`
	rows, err := db.Query(query, userId, batchSize, offset)
	if err != nil {
		log.Println(err)
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
			allowed, err := CheckSuperprivateAccess(post.Id, userId)
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
			morePosts, err := GetPosts(userId, page+1, shortfall, usersOwnProfile)
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
	err := db.QueryRow("SELECT 1 FROM PostAccess WHERE PostId = ? AND UserId = ?", postId, userId).Scan(&exists)
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
	err := db.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", sessionId).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
