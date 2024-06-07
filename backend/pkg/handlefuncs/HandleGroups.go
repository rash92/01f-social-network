package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//structs for frontend to consume

func HandleGroup(w http.ResponseWriter, r *http.Request) {
	var userId string
	var groupId string

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&groupId)
	if err != nil {
		errorMessage := fmt.Sprintf("error decoding groupId: %v", err.Error())
		fmt.Println(err.Error(), errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("user_token")
	if err != nil {
		errorMessage := fmt.Sprintf("error retrieving cookie: %v", err.Error())
		fmt.Println(err.Error(), errorMessage)
		http.Error(w, errorMessage, http.StatusForbidden)
		return
	}
	userId, err = dbfuncs.GetUserIdFromCookie(cookie.Value)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting user from cookie: %v", err.Error())
		fmt.Println(err.Error(), errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	group, err := GetGroup(groupId, userId)
	if err != nil {
		errorMessage := fmt.Sprintf("error getting group: %v", err.Error())
		fmt.Println(err.Error(), errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	fmt.Println(group, "group")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

func GetGroupEventCards(groupId string, userId string) ([]GroupEventCard, error) {
	groupEvents, err := GetGroupEventsByGroupId(groupId)
	if err != nil {
		return []GroupEventCard{}, err
	}

	var eventCards []GroupEventCard
	for _, event := range groupEvents {
		eventCards = append(eventCards, GetGroupEventCard(event, userId))
	}
	return eventCards, err
}

func GetGroupEventCard(event GroupEvent, userId string) GroupEventCard {
	attending, err := dbfuncs.IsUserAttendingEvent(userId, event.Id)
	if err != nil {
		fmt.Println(err)
		return GroupEventCard{event, false}
	}
	if attending {
		return GroupEventCard{event, true}
	}
	return GroupEventCard{event, false}
}

func GetGroupEventsByGroupId(groupId string) ([]GroupEvent, error) {
	groupDbEvents, err := dbfuncs.GetGroupEventsByGroupId(groupId)
	if err != nil {
		return []GroupEvent{}, err
	}
	groupEvents, err := DbGroupEventsToFrontend(groupDbEvents)
	if err != nil {
		return []GroupEvent{}, err
	}
	return groupEvents, err
}

func GetGroupDash(userId string) (GroupDash, error) {
	allGroups, err := dbfuncs.GetAllGroups()
	var groupDash GroupDash
	if err != nil {
		return GroupDash{}, err
	}

	for _, group := range allGroups {

		groupCard, err := GetGroupCard(group.Id, userId)
		if err != nil {
			fmt.Println(err, "err")
			return GroupDash{}, err
		}

		//should be fine?
		groupDash.GroupCards = append(groupDash.GroupCards, groupCard)
	}

	return groupDash, err
}

// gets basic group info, and based on the viewer's user Id finds out their status wrt that group
// (accepted, invited, requested, or not in the database, which is recorded as none)
func GetGroupCard(groupId string, userId string) (GroupCard, error) {
	basicInfo, err := GetBasicGroupInfo(groupId)
	if err != nil {
		fmt.Println(err)
		return GroupCard{}, err
	}
	status, err := dbfuncs.GetGroupStatus(groupId, userId)
	if err == sql.ErrNoRows {
		status = "none"
	}
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("eer", err)
		return GroupCard{}, err
	}
	return GroupCard{basicInfo, status}, nil
}

func GetGroup(groupId string, userId string) (DetailedGroupInfo, error) {
	basicInfo, err := GetBasicGroupInfo(groupId)

	if err != nil {
		log.Fatalln(err, "error  GetBasicGroupInfo")
		return DetailedGroupInfo{}, err
	}

	status, err := dbfuncs.GetGroupStatus(groupId, userId)

	if err == sql.ErrNoRows {
		fmt.Println("User is not a member of this group")
		return DetailedGroupInfo{
			BasicInfo: basicInfo,
			Status:    "none",
		}, nil
	}

	if err != nil {
		fmt.Println("Error getting group status: ", err)
		return DetailedGroupInfo{}, err
	}

	if status != "accepted" {
		fmt.Println("user is not a member of the group")
		return DetailedGroupInfo{
			BasicInfo: basicInfo,
			Status:    status,
		}, nil
	}

	invitedMembers, err := dbfuncs.GetInvitedGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	invitedMembersBasicInfo, err := GetBasicUserInfoFromUsers(invitedMembers)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	requestedMembers, err := dbfuncs.GetRequestedGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	requestedMembersBasicInfo, err := GetBasicUserInfoFromUsers(requestedMembers)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	members, err := dbfuncs.GetGroupMemberIdsByGroupId(groupId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	membersBasicInfo, err := GetBasicUserInfoFromUsers(members)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	groupDbPosts, err := dbfuncs.GetPostsByGroupId(groupId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	groupPosts, err := DbPostsToFrontend(groupDbPosts)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	groupEventCards, err := GetGroupEventCards(groupId, userId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	dbGroupMessages, err := dbfuncs.GetAllGroupMessagesByGroupId(groupId)
	if err != nil {
		return DetailedGroupInfo{}, err
	}

	groupMessages := DbGroupMessagesToFrontend(dbGroupMessages)

	toBeInvited, err := WhoCanIInviteToThisGroup(groupId, userId)

	groupInfo := DetailedGroupInfo{
		BasicInfo:        basicInfo,
		InvitedMembers:   invitedMembersBasicInfo,
		RequestedMembers: requestedMembersBasicInfo,
		Members:          membersBasicInfo,
		Posts:            groupPosts,
		EventCards:       groupEventCards,
		Messages:         groupMessages,
		Status:           status,
		Invite:           toBeInvited,
	}

	return groupInfo, err
}

func GetBasicGroupInfo(groupId string) (BasicGroupInfo, error) {

	group, err := dbfuncs.GetGroupByGroupId(groupId)
	if err != nil {
		return BasicGroupInfo{}, err
	}
	basicInfo := BasicGroupInfo{
		Id:          group.Id,
		CreatorId:   group.CreatorId,
		Name:        group.Title,
		Description: group.Description,
		CreatedAt:   group.CreatedAt,
	}

	return basicInfo, err
}

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
		Id:              dbComment.Id,
		Body:            dbComment.Body,
		CreatorId:       dbComment.CreatorId,
		PostID:          dbComment.PostId,
		CreatedAt:       dbComment.CreatedAt,
		Image:           dbComment.Image,
		CreatorNickname: user.Nickname,
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
	going := len(participantIds)
	groupMembers, err := dbfuncs.GetGroupMemberIdsByGroupId(dbEvent.GroupId)
	if err != nil {
		return GroupEvent{}, err
	}
	notGoing := len(groupMembers) - going
	event := GroupEvent{
		Id:          dbEvent.Id,
		GroupId:     dbEvent.GroupId,
		Title:       dbEvent.Title,
		Description: dbEvent.Description,
		CreatorId:   dbEvent.CreatorId,
		Time:        dbEvent.Time,
		Going:       going,
		NotGoing:    notGoing,
	}
	return event, err
}

func GetBasicUserInfoFromUsers(input []string) ([]BasicUserInfo, error) {
	res := []BasicUserInfo{}
	var err error
	for _, userId := range input {
		basicInfo, err := GetBasicUserInfoById(userId)
		if err != nil {
			return res, err
		}
		res = append(res, basicInfo)
	}
	return res, err
}

func GetBasicUserInfoById(userId string) (BasicUserInfo, error) {
	user, err := dbfuncs.GetUserById(userId)
	if err != nil {
		fmt.Println("couldn't get user from user id")
		return BasicUserInfo{}, err
	}

	basicInfo := BasicUserInfo{
		Avatar:         user.Avatar,
		Id:             user.Id,
		Nickname:       user.Nickname,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		PrivacySetting: user.PrivacySetting,
	}
	return basicInfo, err
}

func WhoCanIInviteToThisGroup(groupId, userId string) ([]BasicUserInfo, error) {
	var users []BasicUserInfo
	included := make(map[BasicUserInfo]struct{})
	followers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(userId)
	if err != nil {
		log.Println("error getting followers from database", err)
		return nil, err
	}
	for _, follower := range followers {
		status, err := dbfuncs.GetGroupStatus(groupId, follower)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			return users, err
		}
		if status == "invited" || status == "accepted" || status == "requested" {
			continue
		}
		basicInfo, err := GetBasicUserInfoById(follower)
		if err != nil {
			log.Println(err)
			return []BasicUserInfo{}, err
		}
		_, ok := included[basicInfo]
		if ok {
			continue
		}
		included[basicInfo] = struct{}{}
		users = append(users, basicInfo)
	}

	following, err := dbfuncs.GetAcceptedFollowingIdsByFollowerId(userId)
	if err != nil {
		log.Println("error getting following from database", err)
		return []BasicUserInfo{}, err
	}
	for _, following := range following {
		status, err := dbfuncs.GetGroupStatus(groupId, following)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			return []BasicUserInfo{}, err
		}
		if status == "invited" || status == "accepted" || status == "requested" {
			fmt.Println("following", following)
			continue
		}
		basicInfo, err := GetBasicUserInfoById(following)
		if err != nil {
			log.Println(err)
			return []BasicUserInfo{}, err
		}
		_, ok := included[basicInfo]
		if ok {
			continue
		}
		included[basicInfo] = struct{}{}
		users = append(users, basicInfo)
	}

	return users, nil
}
