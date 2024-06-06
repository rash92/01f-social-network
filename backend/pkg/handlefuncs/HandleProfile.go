package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// This is what will be returned by the handler.

func UserIdsToBasicUserInfos(input []string) []BasicUserInfo {
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

	if usersOwnProfile {
		pendingFollowers, err := dbfuncs.GetPendingFollowerIdsByFollowingId(ownerId)
		if err != nil {

			fmt.Println(err.Error(), "105")
			errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		profile.PendingFollowers = UserIdsToBasicUserInfos(pendingFollowers)

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
	}

	acceptedFollowers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting followers: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Followers = UserIdsToBasicUserInfos(acceptedFollowers)

	following, err := dbfuncs.GetAcceptedFollowingIdsByFollowerId(ownerId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	profile.Following = UserIdsToBasicUserInfos(following)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)

}

func HandleToggleProfilePrivacy(w http.ResponseWriter, r *http.Request) {

	var privacy PrivcySetting

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&privacy)

	if err != nil {
		errorMessage := fmt.Sprintf("error decoding userId: %v", err.Error())
		fmt.Println(err.Error(), "60")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	fmt.Println("privacy setting", privacy)
	err = dbfuncs.UpdatePrivacySetting(privacy.UserId, privacy.Privacy)
	if err != nil {
		errorMessage := fmt.Sprintf("error updating privacy setting: %v", err.Error())
		fmt.Println(err.Error(), "66")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	response["message"] = "Successfully updated privacy setting"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
