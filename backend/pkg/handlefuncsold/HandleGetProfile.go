package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// I've included whatever structs I needed in this file. They can be replaced
// with the real ones when they're ready, or if anyone knows where they live
// now. UPDATE: We've moved some of these structs to dbfuncs, along with the
// helper functions that access the database.
// I just added this to get rid of the red line under *Image. I don't know
// what Image is really supposed to be.
// type Image []byte

// This is what will be returned by the handler.

type Profile struct {
	Owner     dbfuncs.User
	Posts     []dbfuncs.Post
	Followers []BasicUserInfo
	Following []BasicUserInfo

	// NumberOfPosts     int
	// NumberOfFollowers int
	// NumberOfFollowing int
	PendingFollowers []BasicUserInfo
	IsFollowed       bool
	IsPending        bool
}

func helper(input []string) []BasicUserInfo {
	res := []BasicUserInfo{}

	for _, followerId := range input {
		follower, err := dbfuncs.GetUserById(followerId)
		if err != nil {
			log.Panicf("error getting follower: %v", err.Error())
		}

		basicInfo := BasicUserInfo{
			Avatar:         follower.Avatar,
			Id:             follower.Id,
			Nickname:       follower.Nickname,
			FirstName:      follower.FirstName,
			LastName:       follower.LastName,
			PrivacySetting: follower.PrivacySetting,
		}
		res = append(res, basicInfo)
	}

	return res
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

	profile.Owner, err = dbfuncs.GetUserById(ownerId)

	if err != nil {
		errorMessage := fmt.Sprintf("error getting profile owner: %v", err.Error())
		fmt.Println(err.Error(), "66")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Owner.Password = []byte{}

	cookie, _ := r.Cookie("user_token")
	userId, _ = dbfuncs.GetUserIdFromCookie(cookie.Value)

	if userId == ownerId {
		usersOwnProfile = true
	}

	// Check Follows table to see if there's a row with FollowerId = userId and FollowingId = ownerId.
	profile.IsFollowed, err = dbfuncs.IsFollowing(userId, ownerId)
	if err != nil {
		fmt.Printf("failed to execute query: %v\n", err)
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	profile.IsPending, err = dbfuncs.IsPending(userId, ownerId)
	if err != nil {
		fmt.Printf("failed to execute query: %v\n", err)
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	if !usersOwnProfile && profile.Owner.PrivacySetting == "private" && !profile.IsFollowed {

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

	fmt.Println(ownerId, "")

	//  this  is where we would get the number of posts, followers, and following
	// we decided we not going to use this for now so it can be removed

	// profile.NumberOfPosts, err = dbfuncs.GetNumberOfById(ownerId, "Posts")
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("error getting number of posts: %v", err.Error())
	// 	fmt.Println(err.Error(), "90")
	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
	// 	return
	// }

	// profile.NumberOfFollowers, err = dbfuncs.GetNumberOfFollowersAndFollowing("FollowingId", ownerId)
	// if err != nil {
	// 	errorMessage := fmt.Sprintf("error getting number of followers: %v", err.Error())
	// 	fmt.Println(err.Error())
	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
	// 	return
	// }

	// profile.NumberOfFollowing, err = dbfuncs.GetNumberOfFollowersAndFollowing("FollowerId", ownerId)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	errorMessage := fmt.Sprintf("error getting number of following: %v", err.Error())
	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
	// 	return
	// }

	if usersOwnProfile {
		pendingFollowers, err := dbfuncs.GetPendingFollowerIdsByFollowingId(ownerId)
		if err != nil {

			fmt.Println(err.Error(), "105")
			errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		profile.PendingFollowers = helper(pendingFollowers)

	}

	profile.Posts, err = dbfuncs.GetPosts(userId, 1, 10, usersOwnProfile)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		fmt.Println(err.Error(), "90")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if usersOwnProfile {
		profile.Posts, err = dbfuncs.GetAllPostsByCreatorId(ownerId)
		if err != nil {
			errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}
	} else {
		profile.Posts, err = dbfuncs.GetVisiblePostsForProfile(userId, ownerId)

		if err != nil {
			errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}
		// if profile.IsFollowed {
		// 	a, err := dbfuncs.GetPostsByCreatorId(ownerId)
		// 	if err != nil {
		// 		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		// 		http.Error(w, errorMessage, http.StatusInternalServerError)
		// 		return
		// 	}
		// 	for _, post := range a {
		// 		if post.PrivacyLevel == "public" || post.PrivacyLevel == "private" {
		// 			profile.Posts = append(profile.Posts, post)
		// 		}
		// 		if post.PrivacyLevel == "superprivate" {
		// 			b, err := dbfuncs.GetPostChosenFollowerIdsByPostId(post.Id)
		// 			if err != nil {
		// 				errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		// 				http.Error(w, errorMessage, http.StatusInternalServerError)
		// 				return
		// 			}
		// 			for _, followerId := range b {
		// 				if followerId == userId {
		// 					profile.Posts = append(profile.Posts, post)
		// 				}
		// 			}
		// 		} else {
		// 			a, err := dbfuncs.GetPostsByCreatorId(ownerId)
		// 			if err != nil {
		// 				errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
		// 				http.Error(w, errorMessage, http.StatusInternalServerError)
		// 				return
		// 			}
		// 			for _, post := range a {
		// 				if post.PrivacyLevel == "public" {
		// 					profile.Posts = append(profile.Posts, post)
		// 				}
		// 			}
		// 		}
		// 	}

		// }
	}

	acceptedFollowers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting followers: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Followers = helper(acceptedFollowers)

	following, err := dbfuncs.GetAcceptedFollowingIdsByFollowerId(ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Following = helper(following)

	fmt.Println(profile, "profile")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)

}
