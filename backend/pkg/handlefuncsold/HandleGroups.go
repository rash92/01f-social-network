package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//struct for frontend to consume

type BasicGroupInfo struct {
	Id          string    `json:"Id"`
	CreatorId   string    `json:"CreatorId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type GroupCard struct {
	BasicInfo BasicGroupInfo `json:"BasicInfo"`
	Status    string         `json:"Status"`
}

type GroupEvent struct {
	Id           string          `json:"Id"`
	GroupId      string          `json:"GroupId"`
	Title        string          `json:"Title"`
	Description  string          `json:"Description"`
	CreatorId    string          `json:"CreatorId"`
	Time         time.Time       `json:"Time"`
	Participants []BasicUserInfo `json:"Participants"`
}

type GroupMessage struct {
	Id        string    `json:"Id"`
	SenderId  string    `json:"SenderId"`
	GroupId   string    `json:"GroupId"`
	Message   string    `json:"Message"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type DetailedGroupInfo struct {
	BasicInfo        BasicGroupInfo  `json:"BasicInfo"`
	InvitedMembers   []BasicUserInfo `json:"InvitedMembers"`
	RequestedMembers []BasicUserInfo `json:"RequestedMembers"`
	Members          []BasicUserInfo `json:"Members"`
	Posts            []Post          `json:"Posts"`
	Events           []GroupEvent    `json:"Events"`
	Messages         []GroupMessage  `json:"Messages"`
}

type GroupDash struct {
	GroupCards []GroupCard `json:"GroupCards"`
}

//

func GetGroup(groupId string, userId string) (DetailedGroupInfo, error) {
	groupInfo := DetailedGroupInfo{}
	status, err := dbfuncs.GetGroupStatus(groupId, userId)
	if err == sql.ErrNoRows {
		fmt.Println("User is not a member of this group")
		return groupInfo, err
	}
	if err != nil {
		fmt.Println("Error getting group status: ", err)
		return groupInfo, err
	}
	if status != "accepted" {
		fmt.Println("user is not a member of the group")
		return groupInfo, errors.New("User is not a member of the group, status is: " + status)
	}

	group, err := dbfuncs.GetGroupByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}

	invitedMembers, err := dbfuncs.GetInvitedGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}

	invitedMembersBasicInfo, err := GetBasicUserInfoFromUsers(invitedMembers)
	if err != nil {
		return groupInfo, err
	}

	requestedMembers, err := dbfuncs.GetRequestedGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}

	requestedMembersBasicInfo, err := GetBasicUserInfoFromUsers(requestedMembers)
	if err != nil {
		return groupInfo, err
	}
	members, err := dbfuncs.GetGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}
	membersBasicInfo, err := GetBasicUserInfoFromUsers(members)
	if err != nil {
		return groupInfo, err
	}
	groupDbPosts, err := dbfuncs.GetPostsByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}
	groupPosts, err := DbPostsToFrontend(groupDbPosts)
	if err != nil {
		return groupInfo, err
	}

	groupDbEvents, err := dbfuncs.GetGroupEventsByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}
	groupEvents, err := DbGroupEventsToFrontend(groupDbEvents)
	if err != nil {
		return groupInfo, err
	}

	dbGroupMessages, err := dbfuncs.GetAllGroupMessagesByGroupId(groupId)
	if err != nil {
		return groupInfo, err
	}
	groupMessages := DbGroupMessagesToFrontend(dbGroupMessages)

	groupInfo.BasicInfo = BasicGroupInfo{
		Id:          group.Id,
		CreatorId:   group.CreatorId,
		Name:        group.Title,
		Description: group.Description,
		CreatedAt:   group.CreatedAt,
	}
	groupInfo.InvitedMembers = invitedMembersBasicInfo
	groupInfo.RequestedMembers = requestedMembersBasicInfo
	groupInfo.Members = membersBasicInfo
	groupInfo.Posts = groupPosts
	groupInfo.Events = groupEvents
	groupInfo.Messages = groupMessages
	return groupInfo, err
}

// //handle get profile to use as template for handle groups

// // I've included whatever structs I needed in this file. They can be replaced
// // with the real ones when they're ready, or if anyone knows where they live
// // now. UPDATE: We've moved some of these structs to dbfuncs, along with the
// // helper functions that access the database.
// // I just added this to get rid of the red line under *Image. I don't know
// // what Image is really supposed to be.
// // type Image []byte

// // This is what will be returned by the handler.

// type Profile struct {
// 	Owner     dbfuncs.User
// 	Posts     []dbfuncs.Post
// 	Followers []BasicUserInfo
// 	Following []BasicUserInfo

// 	// NumberOfPosts     int
// 	// NumberOfFollowers int
// 	// NumberOfFollowing int
// 	PendingFollowers []BasicUserInfo
// 	IsFollowed       bool
// 	IsPending        bool
// }

// CONVERSIONS SECTION

func DbGroupMessagesToFrontend(dbGroupMessages []dbfuncs.GroupMessage) []GroupMessage {
	var frontendGroupMessages []GroupMessage
	for _, dbGroupMessage := range dbGroupMessages {
		frontendGroupMessages = append(frontendGroupMessages, DbGroupMessageToFrontend(dbGroupMessage))
	}
	return frontendGroupMessages
}

func DbGroupMessageToFrontend(dbGroupMessage dbfuncs.GroupMessage) GroupMessage {
	frontendGroupMessage := GroupMessage{
		Id:        dbGroupMessage.Id,
		SenderId:  dbGroupMessage.SenderId,
		GroupId:   dbGroupMessage.GroupId,
		Message:   dbGroupMessage.Message,
		CreatedAt: dbGroupMessage.CreatedAt,
	}
	return frontendGroupMessage
}

func DbPostsToFrontend(dbPosts []dbfuncs.Post) ([]Post, error) {
	var frontendPosts []Post
	for _, dbPost := range dbPosts {
		frontendPost, err := DbPostToFrontend(dbPost)
		if err != nil {
			return frontendPosts, err
		}
		frontendPosts = append(frontendPosts, frontendPost)
	}
	return frontendPosts, nil
}

func DbPostToFrontend(dbPost dbfuncs.Post) (Post, error) {
	chosenFollowers, err := dbfuncs.GetPostChosenFollowerIdsByPostId(dbPost.Id)
	if err != nil {
		return Post{}, err
	}
	comments, err := DbCommentsToFrontend(dbPost.Comments)
	if err != nil {
		return Post{}, err
	}

	frontendPost := Post{
		Id:              dbPost.Id,
		Title:           dbPost.Title,
		Body:            dbPost.Body,
		CreatedAt:       dbPost.CreatedAt,
		Comments:        comments,
		Likes:           dbPost.Likes,
		Dislikes:        dbPost.Dislikes,
		PrivacyLevel:    dbPost.PrivacyLevel,
		CreatorId:       dbPost.CreatorId,
		Image:           dbPost.Image,
		GroupId:         dbPost.GroupId,
		ChosenFollowers: chosenFollowers,
	}
	return frontendPost, err
}

func DbCommentsToFrontend(dbComments []dbfuncs.Comment) ([]Comment, error) {
	var frontendComments []Comment
	var err error
	for _, dbComment := range dbComments {
		frontendComment, err := DbCommentToFrontend(dbComment)
		if err != nil {
			return frontendComments, err
		}
		frontendComments = append(frontendComments, frontendComment)
	}
	return frontendComments, err
}

func DbCommentToFrontend(dbComment dbfuncs.Comment) (Comment, error) {
	user, err := dbfuncs.GetUserById(dbComment.CreatorId)
	if err != nil {
		return Comment{}, err
	}
	comment := Comment{
		Id:        dbComment.Id,
		Body:      dbComment.Body,
		UserId:    dbComment.CreatorId,
		PostId:    dbComment.PostId,
		CreatedAt: dbComment.CreatedAt,
		Image:     dbComment.Image,
		Username:  user.Nickname,
	}
	return comment, err
}

func DbGroupEventsToFrontend(dbEvents []dbfuncs.GroupEvent) ([]GroupEvent, error) {
	var frontendPosts []GroupEvent
	var err error
	for _, dbEvent := range dbEvents {
		event, err := DbGroupEventToFrontend(dbEvent)
		if err != nil {
			return frontendPosts, err
		}
		frontendPosts = append(frontendPosts, event)
	}
	return frontendPosts, err
}

func DbGroupEventToFrontend(dbEvent dbfuncs.GroupEvent) (GroupEvent, error) {
	participantIds, err := dbfuncs.GetEventParticipantIdsByEventId(dbEvent.Id)
	if err != nil {
		return GroupEvent{}, err
	}
	participants, err := GetBasicUserInfoFromUsers(participantIds)
	if err != nil {
		return GroupEvent{}, err
	}
	event := GroupEvent{
		Id:           dbEvent.Id,
		GroupId:      dbEvent.GroupId,
		Title:        dbEvent.Title,
		Description:  dbEvent.Description,
		CreatorId:    dbEvent.CreatorId,
		Time:         dbEvent.Time,
		Participants: participants,
	}
	return event, err
}

func GetBasicUserInfoFromUsers(input []string) ([]BasicUserInfo, error) {
	res := []BasicUserInfo{}
	var err error
	for _, userId := range input {
		user, err := dbfuncs.GetUserById(userId)
		if err != nil {
			fmt.Println("couldn't get user from user id")
			return res, err
		}

		basicInfo := BasicUserInfo{
			Avatar:         user.Avatar,
			Id:             user.Id,
			Nickname:       user.Nickname,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			PrivacySetting: user.PrivacySetting,
		}
		res = append(res, basicInfo)
	}
	return res, err
}

// func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
// 	var userId string
// 	var ownerId string
// 	var profile Profile
// 	var usersOwnProfile bool

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&ownerId)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("error decoding userId: %v", err.Error())
// 		fmt.Println(err.Error(), "60")
// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		return
// 	}

// 	profile.Owner, err = dbfuncs.GetUserById(ownerId)

// 	if err != nil {
// 		errorMessage := fmt.Sprintf("error getting profile owner: %v", err.Error())
// 		fmt.Println(err.Error(), "66")
// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		return
// 	}

// 	profile.Owner.Password = []byte{}

// 	cookie, _ := r.Cookie("user_token")
// 	userId, _ = dbfuncs.GetUserIdFromCookie(cookie.Value)

// 	if userId == ownerId {
// 		usersOwnProfile = true
// 	}

// 	// Check Follows table to see if there's a row with FollowerId = userId and FollowingId = ownerId.
// 	profile.IsFollowed, err = dbfuncs.IsFollowing(userId, ownerId)
// 	if err != nil {
// 		fmt.Printf("failed to execute query: %v\n", err)
// 		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
// 		return
// 	}

// 	profile.IsPending, err = dbfuncs.IsPending(userId, ownerId)
// 	if err != nil {
// 		fmt.Printf("failed to execute query: %v\n", err)
// 		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
// 		return
// 	}

// 	if !usersOwnProfile && profile.Owner.PrivacySetting == "private" && !profile.IsFollowed {

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		if err := json.NewEncoder(w).Encode(profile); err != nil {
// 			// Handle JSON encoding error
// 			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
// 			fmt.Println("Failed to encode JSON:", err)
// 			return
// 		}

// 		return
// 	}

// 	fmt.Println(ownerId, "")

// 	//  this  is where we would get the number of posts, followers, and following
// 	// we decided we not going to use this for now so it can be removed

// 	// profile.NumberOfPosts, err = dbfuncs.GetNumberOfById(ownerId, "Posts")
// 	// if err != nil {
// 	// 	errorMessage := fmt.Sprintf("error getting number of posts: %v", err.Error())
// 	// 	fmt.Println(err.Error(), "90")
// 	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// profile.NumberOfFollowers, err = dbfuncs.GetNumberOfFollowersAndFollowing("FollowingId", ownerId)
// 	// if err != nil {
// 	// 	errorMessage := fmt.Sprintf("error getting number of followers: %v", err.Error())
// 	// 	fmt.Println(err.Error())
// 	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// profile.NumberOfFollowing, err = dbfuncs.GetNumberOfFollowersAndFollowing("FollowerId", ownerId)
// 	// if err != nil {
// 	// 	fmt.Println(err.Error())
// 	// 	errorMessage := fmt.Sprintf("error getting number of following: %v", err.Error())
// 	// 	http.Error(w, errorMessage, http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	if usersOwnProfile {
// 		pendingFollowers, err := dbfuncs.GetPendingFollowerIdsByFollowingId(ownerId)
// 		if err != nil {

// 			fmt.Println(err.Error(), "105")
// 			errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
// 			http.Error(w, errorMessage, http.StatusInternalServerError)
// 			return
// 		}

// 		profile.PendingFollowers = helper(pendingFollowers)

// 	}

// 	profile.Posts, err = dbfuncs.GetPosts(userId, 1, 10, usersOwnProfile)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 		fmt.Println(err.Error(), "90")
// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		return
// 	}

// 	if usersOwnProfile {
// 		profile.Posts, err = dbfuncs.GetPostsByCreatorId(ownerId)
// 		if err != nil {
// 			errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 			http.Error(w, errorMessage, http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		profile.Posts, err = dbfuncs.GetVisiblePostsForProfile(userId, ownerId)

// 		if err != nil {
// 			errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 			http.Error(w, errorMessage, http.StatusInternalServerError)
// 			return
// 		}
// 		// if profile.IsFollowed {
// 		// 	a, err := dbfuncs.GetPostsByCreatorId(ownerId)
// 		// 	if err != nil {
// 		// 		errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 		// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		// 		return
// 		// 	}
// 		// 	for _, post := range a {
// 		// 		if post.PrivacyLevel == "public" || post.PrivacyLevel == "private" {
// 		// 			profile.Posts = append(profile.Posts, post)
// 		// 		}
// 		// 		if post.PrivacyLevel == "superprivate" {
// 		// 			b, err := dbfuncs.GetPostChosenFollowerIdsByPostId(post.Id)
// 		// 			if err != nil {
// 		// 				errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 		// 				http.Error(w, errorMessage, http.StatusInternalServerError)
// 		// 				return
// 		// 			}
// 		// 			for _, followerId := range b {
// 		// 				if followerId == userId {
// 		// 					profile.Posts = append(profile.Posts, post)
// 		// 				}
// 		// 			}
// 		// 		} else {
// 		// 			a, err := dbfuncs.GetPostsByCreatorId(ownerId)
// 		// 			if err != nil {
// 		// 				errorMessage := fmt.Sprintf("error getting posts: %v", err.Error())
// 		// 				http.Error(w, errorMessage, http.StatusInternalServerError)
// 		// 				return
// 		// 			}
// 		// 			for _, post := range a {
// 		// 				if post.PrivacyLevel == "public" {
// 		// 					profile.Posts = append(profile.Posts, post)
// 		// 				}
// 		// 			}
// 		// 		}
// 		// 	}

// 		// }
// 	}

// 	acceptedFollowers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(ownerId)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("error getting followers: %v", err.Error())
// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		return
// 	}

// 	profile.Followers = helper(acceptedFollowers)

// 	following, err := dbfuncs.GetAcceptedFollowingIdsByFollowerId(ownerId)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("error getting following: %v", err.Error())
// 		http.Error(w, errorMessage, http.StatusInternalServerError)
// 		return
// 	}

// 	profile.Following = helper(following)

// 	fmt.Println(profile, "profile")
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(profile)

// }
