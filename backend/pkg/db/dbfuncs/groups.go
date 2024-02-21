package dbfuncs

import (
	"database/sql"
	"fmt"
	"log"
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

// to do
func GetEventsById(id string) ([]Event, error) {
	var events []Event
	var err error
	return events, err
}

func GetGroupCreatorByGroupId(groupId string) (string, error) {
	return "", nil
}

// possibly split into get accepted vs get pending? whatever status means
func GetGroupMembersByGroupId(groupId string) ([]string, error) {
	var members []string
	return members, nil
}

// old version of above
func GetGroupMembers(groupId string) []string {
	var result []string
	lock.RLock()
	rows, err := database.Query("SELECT * FROM GroupMembers WHERE GroupId = ?", groupId)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		rows.Close()
		lock.RUnlock()
	}()

	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, userId)
	}

	return result
}

// to do, see if we want choice field or just entry vs not for going vs not going
func ToggleAttendEvent(eventId, userId string) error {
	// Adapt this code for dbfuncs.CreateEvent.
	// newLike, err := database.Prepare("INSERT INTO GroupEventParticipants VALUES (?,?,?)")
	// if err != nil {
	// 	return err
	// }
	toggleAttendance, err := database.Prepare("UPDATE GroupEventParticipants SET Choice=? WHERE EventId=? AND UserId=?")
	if err != nil {
		return err
	}

	row := database.QueryRow("SELECT Choice FROM GroupEventParticipants WHERE EventId=? AND UserId=?", eventId, userId)
	var choice bool
	err = row.Scan(&choice)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no matching row found for EventId %s and UserId %s", eventId, userId)
		} else {
			return err
		}
	}

	if choice {
		_, err = toggleAttendance.Exec(false, eventId, userId)
	} else {
		_, err = toggleAttendance.Exec(true, eventId, userId)
	}

	return err
}
