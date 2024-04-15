package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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
	_, err = statement.Exec(privateMessage.Id, privateMessage.SenderId, privateMessage.RecipientId, privateMessage.Message, privateMessage.CreatedAt)

	return err
}

// need to actually test it, currently returning Ids and throwing away associated most recent message time
func GetAllUserIdsSortedByLastPrivateMessage(userId string) ([]string, error) {
	var userIds []string
	row, err := db.Query(`
	SELECT UserId, Max(CreatedAt) 
		FROM
		(
			SELECT RecipientId AS UserId, CreatedAt
				FROM PrivateMessages 
				WHERE SenderId=?
			UNION
			SELECT SenderId AS UserId, CreatedAt
				FROM PrivateMessages 
				WHERE RecipientId=?
		)
    	GROUP BY UserId`,
		userId, userId)

	if err == sql.ErrNoRows {
		return userIds, nil
	}
	if err != nil {
		return userIds, err
	}
	defer row.Close()
	for row.Next() {
		var userId string
		var mostRecentTime time.Time
		err = row.Scan(&userId, &mostRecentTime)
		if err != nil {
			return userIds, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, err
}

// need to actually test it, currently returning Ids and throwing away nickname
func GetUnmessagedUserIdsSortedAlphabetically(userId string) ([]string, error) {
	var unmessagedUserIds []string
	row, err := db.Query(`
	SELECT Id, Nickname 
		FROM Users 
		WHERE Id IN
		(
			SELECT FollowerId 
				FROM Follows 
				WHERE FollowingId=? AND Status='accepted' 
			UNION 
			SELECT FollowingId 
				FROM Follows 
				WHERE FollowerId=? AND Status='accepted'
		)
	ORDER BY Nickname ASC`,
		userId, userId)

	if err == sql.ErrNoRows {
		return unmessagedUserIds, nil
	}
	if err != nil {
		return unmessagedUserIds, err
	}
	defer row.Close()
	for row.Next() {
		var unmessagedUserId string
		var unmessagedUserName string
		err = row.Scan(&unmessagedUserId, &unmessagedUserName)
		if err != nil {
			return unmessagedUserIds, err
		}
		unmessagedUserIds = append(unmessagedUserIds, unmessagedUserId)
	}
	return unmessagedUserIds, err
}

func GetAllPrivateMessagesByUserId(userId string) ([]PrivateMessage, error) {
	return []PrivateMessage{}, nil
}

func GetRecentPrivateMessagesByUserId(userId string, numberOfMessages, offset int) ([]PrivateMessage, error) {
	return []PrivateMessage{}, nil
}

//TO DO: get 10 at a time? decide if doing it through SQL or get all and do in handlefunc
