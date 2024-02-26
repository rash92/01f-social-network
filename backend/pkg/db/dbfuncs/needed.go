package dbfuncs

import (
	"database/sql"
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

func IsUserPrivate(id string) (bool, error) {
	var privacySetting string
	err := db.QueryRow("SELECT PrivacySetting FROM Users WHERE id = $1", id).Scan(&privacySetting)
	if err != nil {
		return false, err
	}
	return privacySetting == "private", nil
}

func AddFollow(follow *Follow) error {
	_, err := db.Exec("INSERT INTO followers (FollowerId, FollowingId) VALUES ($1, $2)", follow.FollowerId, follow.FollowingId)
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
	return notification, nil
}

func AcceptFollow(followId string, followingId string) error {
	return nil
}

func RejectFollow(followId string, followingId string) error {
	return nil
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

func GetFollowersByFollowingId(id string) ([]string, error) {
	var followers []string
	return followers, nil
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

func GetUserById(id string) (User, error) {
	user := User{}
	return user, nil
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

func AddPrivateMessage(message *PrivateMessage) error {
	return nil
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
