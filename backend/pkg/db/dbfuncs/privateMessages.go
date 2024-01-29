package dbfuncs

import (
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
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	statement.Exec(privateMessage.Id, privateMessage.SenderId, privateMessage.RecipientId, privateMessage.Message, privateMessage.CreatedAt)

	return nil
}
