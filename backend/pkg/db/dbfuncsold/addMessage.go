package dbfuncs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func AddMessage(sender_id, recipient_id, message, Type string) (uuid.UUID, time.Time, error) {

	if Type != "message" || sender_id == "" || recipient_id == "" || message == "" {

		return uuid.Nil, time.Now(), fmt.Errorf("invalid error")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err, "AddMessage line 21")
		return uuid.Nil, time.Now(), err
	}
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO Messages VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println(err, "AddMessage line 27")
		return uuid.Nil, created, err
	}

	_, err = statement.Exec(id, sender_id, recipient_id, message, created)
	fmt.Println(err, "AddMessage line 27")
	return id, created, nil
}
