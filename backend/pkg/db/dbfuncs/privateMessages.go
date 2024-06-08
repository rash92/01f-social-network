package dbfuncs

import (
	"database/sql"

	"time"

	"github.com/google/uuid"
)

func AddPrivateMessage(privateMessage *PrivateMessage) error {
	dbLock.Lock()
	defer dbLock.Unlock()

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
		GROUP BY UserId
		ORDER BY Max(CreatedAt) DESC`, userId, userId)

	if err == sql.ErrNoRows {
		return userIds, nil
	}
	if err != nil {
		return userIds, err
	}
	defer row.Close()
	for row.Next() {
		var userId string
		var mostRecentTime string
		err = row.Scan(&userId, &mostRecentTime)
		if err != nil {
			return userIds, err
		}
		userIds = append(userIds, userId)
	}

	err = row.Err()
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

	err = row.Err()
	return unmessagedUserIds, err
}

func GetAllPrivateMessagesByUserId(CurrUser, OtherUser string) ([]PrivateMessage, error) {
	messages := []PrivateMessage{}
	query := `
		SELECT * FROM PrivateMessages
		WHERE (senderId = ? AND RecipientId = ?) OR (SenderId = ? AND RecipientId = ?)
		ORDER BY CreatedAt DESC
	`
	rows, err := db.Query(query, CurrUser, OtherUser, OtherUser, CurrUser)
	if err != nil {

		return messages, err
	}

	for rows.Next() {
		var message PrivateMessage
		err := rows.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.Message, &message.CreatedAt)

		if err != nil {
			return nil, err

		}

		user, err := GetUserById(message.SenderId)

		if err != nil {
			return nil, err

		}
		message.Avatar = user.Avatar
		message.Nickname = user.Nickname

		messages = append(messages, message)
	}

	// Reverse the order of messages
	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - i - 1
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil

}

func GetLimitedPrivateMessages(CurrentUserId string, OtherUserId string, numberOfMessages, offset int) ([]PrivateMessage, error) {
	messages := []PrivateMessage{}
	query := `
		SELECT * FROM PrivateMessages
		WHERE (senderId = ? AND RecipientId = ?) OR (SenderId = ? AND RecipientId = ?)
		ORDER BY CreatedAt DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, CurrentUserId, OtherUserId, OtherUserId, CurrentUserId, numberOfMessages, offset)
	if err != nil {

		return messages, err
	}

	for rows.Next() {
		var message PrivateMessage
		err := rows.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.Message, &message.CreatedAt)

		if err != nil {
			return nil, err

		}

		user, err := GetUserById(message.SenderId)

		if err != nil {
			return nil, err

		}
		message.Avatar = user.Avatar
		message.Nickname = user.Nickname

		messages = append(messages, message)
	}

	// Reverse the order of messages
	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - i - 1
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
