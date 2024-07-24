package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func AddGroup(group *Group) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	group.Id = id.String()
	group.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO Groups VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(group.Id, group.Title, group.Description, group.CreatorId, group.CreatedAt)

	return err
}

func AddGroupMember(groupMember *GroupMember) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("INSERT INTO GroupMembers VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupMember.GroupId, groupMember.UserId, groupMember.Status)

	return err
}

func UpdateGroupMember(groupMember *GroupMember) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("UPDATE GroupMembers SET Status=? WHERE GroupId=? AND UserId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupMember.Status, groupMember.GroupId, groupMember.UserId)

	return err
}

func DeleteGroupMember(groupMember *GroupMember) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("DELETE FROM GroupMembers WHERE GroupId=? AND UserId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupMember.GroupId, groupMember.UserId)

	return err
}

func AddGroupEvent(groupEvent *GroupEvent) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	groupEvent.Id = id.String()
	statement, err := db.Prepare("INSERT INTO GroupEvents VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEvent.Id, groupEvent.GroupId, groupEvent.Title, groupEvent.Description, groupEvent.CreatorId, groupEvent.Time)

	return err
}

func DeleteGroupEvent(groupEventId string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("DELETE FROM GroupEvents WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventId)

	return err
}

func AddGroupEventParticipant(groupEventParticipant *GroupEventParticipant) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("INSERT INTO GroupEventParticipants VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventParticipant.EventId, groupEventParticipant.UserId, groupEventParticipant.GroupId)

	return err
}

func DeleteGroupEventParticipant(participant *GroupEventParticipant) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("DELETE FROM GroupEventParticipants WHERE UserId=? AND EventId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(participant.UserId, participant.EventId)

	return err
}

func AddGroupMessage(groupMessage *GroupMessage) error {
	dbLock.Lock()
	defer dbLock.Unlock()
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

func GetAllGroupMessagesByGroupId(groupId string) ([]GroupMessage, error) {
	var groupMessages []GroupMessage
	query := `
	SELECT * FROM GroupMessages WHERE GroupId=?
	ORDER BY CreatedAt DESC
	`
	rows, err := db.Query(query, groupId)
	if err == sql.ErrNoRows {
		return groupMessages, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message GroupMessage
		err := rows.Scan(&message.Id, &message.SenderId, &message.GroupId, &message.Message, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		sender, err := GetUserById(message.SenderId)
		if err != nil {
			return nil, err
		}
		message.Nickname = sender.Nickname
		message.Avatar = sender.Avatar
		groupMessages = append(groupMessages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groupMessages, err

}

func GetLimitedGroupMessagesByGroupId(groupId string, numberOfMessages, offset int) ([]GroupMessage, error) {
	var groupMessages []GroupMessage
	query := `
	SELECT * FROM GroupMessages WHERE GroupId=?
	ORDER BY CreatedAt DESC
	LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, groupId, numberOfMessages, offset)
	if err == sql.ErrNoRows {
		return groupMessages, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message GroupMessage
		err := rows.Scan(&message.Id, &message.SenderId, &message.GroupId, &message.Message, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		groupMessages = append(groupMessages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groupMessages, err
}

func GetGroupEventsByGroupId(groupId string) ([]GroupEvent, error) {
	var GroupEvents []GroupEvent
	rows, err := db.Query("SELECT Id, GroupId, Title, Description, CreatorId, Time FROM GroupEvents WHERE GroupId=?", groupId)
	if err == sql.ErrNoRows {
		return GroupEvents, nil
	}
	if err != nil {
		return GroupEvents, err
	}
	defer rows.Close()
	for rows.Next() {
		var GroupEvent GroupEvent
		err = rows.Scan(&GroupEvent.Id, &GroupEvent.GroupId, &GroupEvent.Title, &GroupEvent.Description, &GroupEvent.CreatorId, &GroupEvent.Time)
		if err != nil {
			return GroupEvents, err
		}
		GroupEvents = append(GroupEvents, GroupEvent)
	}
	return GroupEvents, err
}

func GetGroupCreatorIdByGroupId(groupId string) (string, error) {
	var CreatorId string
	err := db.QueryRow("SELECT CreatorId FROM Groups WHERE Id=?", groupId).Scan(&CreatorId)

	return CreatorId, err
}

func GetGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	rows, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "accepted")
	if err == sql.ErrNoRows {
		return GroupMemberIds, nil
	}
	if err != nil {
		return GroupMemberIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var GroupMemberId string
		err = rows.Scan(&GroupMemberId)
		if err != nil {
			return GroupMemberIds, err
		}
		GroupMemberIds = append(GroupMemberIds, GroupMemberId)
	}
	return GroupMemberIds, nil
}

func GetRequestedGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	rows, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "requested")
	if err == sql.ErrNoRows {
		return GroupMemberIds, nil
	}
	if err != nil {
		return GroupMemberIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var GroupMemberId string
		err = rows.Scan(&GroupMemberId)
		if err != nil {
			return GroupMemberIds, err
		}
		GroupMemberIds = append(GroupMemberIds, GroupMemberId)
	}
	return GroupMemberIds, nil
}

func GetInvitedGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	rows, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "invited")
	if err == sql.ErrNoRows {
		return GroupMemberIds, nil
	}
	if err != nil {
		return GroupMemberIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var GroupMemberId string
		err = rows.Scan(&GroupMemberId)
		if err != nil {
			return GroupMemberIds, err
		}
		GroupMemberIds = append(GroupMemberIds, GroupMemberId)
	}
	return GroupMemberIds, nil
}

//  get all groups

func GetAllGroups() ([]Group, error) {
	var groups []Group
	rows, err := db.Query("SELECT Id, Title, Description, CreatorId, CreatedAt FROM Groups")
	if err == sql.ErrNoRows {
		return groups, nil
	}
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		err = rows.Scan(&group.Id, &group.Title, &group.Description, &group.CreatorId, &group.CreatedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

func GetGroupStatus(groupId string, userId string) (string, error) {
	var status string
	err := db.QueryRow("SELECT Status FROM GroupMembers WHERE GroupId=? AND UserId=?", groupId, userId).Scan(&status)
	return status, err
}

func GetPostsByGroupId(userId, groupId string) ([]Post, error) {
	var posts []Post
	rows, err := db.Query("SELECT * FROM Posts WHERE GroupId=?", groupId)
	if err == sql.ErrNoRows {
		return posts, nil
	}
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var groupId sql.NullString

		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &groupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return posts, err
		}
		post.GroupId = StringNull(groupId)

		post.Likes, post.Dislikes, err = CountPostReacts(post.Id)
		if err != nil {
			return posts, err
		}
		user, err := GetUserById(post.CreatorId)
		if err != nil {
			return posts, err
		}

		post.CreatorNickname = user.Nickname
		post.UserLikeDislike, err = GetUserLikeDislike(userId, post.Id)
		if err != nil {
			return posts, err
		}

		post.Comments, err = GetAllCommentsByPostId(post.Id)

		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}
	err = rows.Err()
	return posts, err
}

func GetEventParticipantIdsByEventId(eventId string) ([]string, error) {
	var participants []string
	rows, err := db.Query("SELECT UserId FROM GroupEventParticipants WHERE EventId=?", eventId)
	if err == sql.ErrNoRows {
		return participants, nil
	}
	if err != nil {
		return participants, err
	}
	defer rows.Close()
	for rows.Next() {
		var participant string
		err = rows.Scan(&participant)
		if err != nil {
			return participants, err
		}
		participants = append(participants, participant)
	}
	err = rows.Err()
	return participants, err
}

func IsUserAttendingEvent(userId string, eventId string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM GroupEventParticipants WHERE UserId=? AND EventId=?", userId, eventId).Scan(&count)
	return count > 0, err
}

func DeleteGroup(groupEventParticipantId string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("DELETE FROM Groups WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventParticipantId)

	return err
}

func GetGroupByGroupId(groupId string) (Group, error) {
	var group Group
	err := db.QueryRow("SELECT * FROM Groups WHERE Id=?", groupId).Scan(&group.Id, &group.Title, &group.Description, &group.CreatorId, &group.CreatedAt)

	return group, err
}
