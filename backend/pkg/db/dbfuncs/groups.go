package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func AddGroup(group *Group) error {
	//may want to use autoincrement instead of uuids?
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

func AddGroupMember(groupmember *GroupMember) error {
	statement, err := db.Prepare("INSERT INTO GroupMembers VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupmember.GroupId, groupmember.UserId, groupmember.Status)

	return err
}

func AddGroupEvent(groupEvent *GroupEvent) error {
	//may want to use autoincrement instead of uuids?
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

func AddGroupEventParticipant(groupEventParticipant *GroupEventParticipant) error {
	statement, err := db.Prepare("INSERT INTO GroupEventParticipants VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventParticipant.EventId, groupEventParticipant.UserId, groupEventParticipant.GroupId)

	return err
}

func DeleteGroupEventParticipant(groupEventParticipantId string) error {
	statement, err := db.Prepare("DELETE FROM GroupEventParticipants WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventParticipantId)

	return err
}

func AddGroupMessage(groupMessage *GroupMessage) error {
	//may want to use autoincrement instead of uuids?
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

// to do
func GetGroupEventsById(id string) ([]GroupEvent, error) {
	var GroupEvents []GroupEvent
	row, err := db.Query("SELECT Id, GroupId, Title, Description, CreatorId, Time FROM GroupEvents WHERE Id=?", id)
	if err == sql.ErrNoRows {
		return GroupEvents, nil
	}
	if err != nil {
		return GroupEvents, err
	}
	defer row.Close()
	for row.Next() {
		var GroupEvent GroupEvent
		err = row.Scan(&GroupEvent)
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

func GetRequestedGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	row, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "requested")
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

func GetInvitedGroupMemberIdsByGroupId(groupId string) ([]string, error) {
	var GroupMemberIds []string
	row, err := db.Query("SELECT UserId FROM GroupMembers WHERE GroupId=? AND Status=?", groupId, "invited")
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
