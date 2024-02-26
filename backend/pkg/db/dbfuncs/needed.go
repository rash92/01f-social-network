package dbfuncs

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Get post from comment id done aaccording to this plan.

/* addEvent has different signature but exists, getFollowersByFollowingIds is split into get accepted vs pending followers, and both follower by following and following by follower, GetPostPrivacyLevelByCommentId and GetCreatedAtByUserId are unnecessary as can just get post and get user and access fields from them, getPostChosenFollowersByPostId renamed to getPostChosenFollowerIdsByPostId, getGroupMembersByGroupId split into getting accepted vs invited vs requested group members

getPostByCommentId split into getPostIdByCommentId and getPostById

AddPost, IsUserPrivate, AddFollow, AddNotification, AcceptFollow, RejectFollow, AddEvent, GetEventsById, NotificationSeen, AddGrouMember, GetGroupCreatorByGroupId, AddComment, GetCommentById, GetUserById all exist and are complete and have same inputs/outputs as placeholder functions

re: above: AddEvent renamed to AddGroupEvent, GetEventsById renamed to GetGroupEventsById

ToggleAttendEvent split into AddGroupEventParticipant and DeleteGroupEventParticipant

*/

func AddPost(post *Post) error {
	_, err := db.Exec("INSERT INTO posts (Id, Title, Body, CreatorId, GroupId, CreatedAt, Image, PrivacyLevel) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", post.Id, post.Title, post.Body, post.CreatorId, post.GroupId, post.CreatedAt, post.Image, post.PrivacyLevel)
	return err
}

func AddNotification(notification *Notification) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	notification.Id = id.String()
	notification.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(notification.Id, notification.Body, notification.Type, notification.CreatedAt, notification.ReceiverId, notification.SenderId, notification.Seen)
	return err
}

func GetNotificationById(id string) (Notification, error) {
	var notification Notification
	err := db.QueryRow("SELECT Id, Body, Type, CreatedAt, ReceiverId, SenderId, Seen FROM Notifications WHERE Id=?", id).Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen)
	return notification, err
}

func IsUserPrivate(userId string) (bool, error) {
	var privacySetting string
	err := db.QueryRow("SELECT PrivacySetting FROM Users WHERE Id=?", userId).Scan(&privacySetting)
	if err != nil {
		return false, err
	}
	if privacySetting == "public" {
		return false, nil
	}
	if privacySetting == "private" {
		return true, nil
	}
	return false, errors.New("privacy setting not recognized, should be either 'private' or 'public'")
}

func GetUserById(id string) (User, error) {
	var user User
	err := db.QueryRow("SELECT Id, Nickname, FirstName, LastName, Email, Password, Profile, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users WHERE Id=?", id).Scan(&user.Id, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Profile, &user.AboutMe, &user.PrivacySetting, &user.DOB, &user.CreatedAt)
	return user, err
}

func AddEvent(event *Event) (string, time.Time, error) {
	var err error
	return "", time.Now(), err
}

func GetEventsById(id string) ([]Event, error) {
	var events []Event
	var err error
	return events, err
}

func NotificationSeen(notificationId string) error {
	return nil
}

func AddGroupMember(member *GroupMember) error {
	return nil
}

func GetGroupCreatorFromGroupId(groupId string) (string, error) {
	return "", nil
}

func AddComment(comment *Comment) error {
	return nil
}

func GetCommentById(id string) (Comment, error) {
	var comment Comment
	return comment, nil
}

func GetPostChosenFollowerIdsByPostId(id string) ([]string, error) {
	var followers []string
	return followers, nil
}

func GetGroupMembersByGroupId(groupId string) ([]string, error) {
	var members []string
	return members, nil
}

func GetPostPrivacyLevelByCommentId(id string) (string, error) {
	return "", nil
}

func GetPostIdByCommentId(id string) (string, error) {
	return "", nil
}

func GetPostById(id string) (Post, error) {
	post := Post{}
	return post, nil
}

func GetCreatedAtByUserId(id string) (time.Time, error) {
	user, err := GetUserById(id)
	if err != nil {
		return time.Time{}, err
	}

	return user.CreatedAt, nil
}

func GetGroupCreatorByGroupId(id string) (string, error) {
	return "", nil
}

func AddPrivateMessage(privateMessage *PrivateMessage) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	privateMessage.Id = id.String()
	privateMessage.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO PrivateMessages VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(privateMessage.Id, privateMessage.SenderId, privateMessage.ReceiverId, privateMessage.Message, privateMessage.CreatedAt)
	return err
}

func AddGroupMessage(groupMessage *GroupMessage) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	groupMessage.Id = id.String()
	groupMessage.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO GroupMessages VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupMessage.Id, groupMessage.SenderId, groupMessage.GroupId, groupMessage.Message, groupMessage.CreatedAt)
	return err
}

func GetGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	row, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "accepted")
	if err == sql.ErrNoRows {
		return GroupMemberIds, nil
	}
	if err != nil {
		return GroupMemberIds, err
	}
	defer row.Close()
	for row.Next() {
		var GroupMemberId string
		err = row.Scan(&GroupMemberId)
		if err != nil {
			return GroupMemberIds, err
		}
		GroupMemberIds = append(GroupMemberIds, GroupMemberId)
	}
	return GroupMemberIds, nil
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
	statement, err := db.Prepare("UPDATE Follows SET Status=?, WHERE FollowerId=? AND FollowingId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec("accepted", followerId, followingId)
	return err
}

// may not use this and can just delete follow from table instead on rejected a follow request
func RejectFollow(followerId, followingId string) error {
	statement, err := db.Prepare("UPDATE Follows SET Status=?, WHERE FollowerId=? AND FollowingId=?")
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
	_, err = statement.Exec(followerId, followerId)
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
