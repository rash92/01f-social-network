package dbfuncs

import (
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
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(group.Id, group.Title, group.Description, group.CreatorId, group.CreatedAt)

	return err
}

func AddGroupMember(groupmember *GroupMember) error {
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?)")
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
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEvent.Id, groupEvent.GroupId, groupEvent.Title, groupEvent.Description, groupEvent.CreatorId, groupEvent.Time)

	return err
}

func AddGroupEventParticipant(groupEventParticipant *GroupEventParticipant) error {
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupEventParticipant.EventId, groupEventParticipant.UserId, groupEventParticipant.GroupId)

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
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(groupMessage.Id, groupMessage.SenderId, groupMessage.GroupId, groupMessage.Message, groupMessage.CreatedAt)

	return err
}
