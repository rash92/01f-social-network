package dbfuncs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func AddMessage(sender_id, recipient_id, message, Type string) (uuid.UUID, time.Time, error) {

	isWrongType := Type != "PrivateMessage" && Type != "GroupMessage"

	if isWrongType || sender_id == "" || recipient_id == "" || message == "" {

		return uuid.Nil, time.Now(), fmt.Errorf("invalid type")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err, "AddMessage line 21")
		return uuid.Nil, time.Now(), err
	}
	created := time.Now()

	lock.Lock()
	statement, err := database.Prepare(fmt.Sprintf("INSERT INTO %s VALUES (?,?,?,?,?)", Type))
	if err != nil {
		fmt.Println(err, "AddMessage line 27")
		return uuid.Nil, created, err
	}
	_, err = statement.Exec(id, sender_id, recipient_id, message, created)
	fmt.Println(err, "AddMessage line 27")
	lock.Unlock()

	return id, created, nil
}
