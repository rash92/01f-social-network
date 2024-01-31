package dbfuncs

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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

	if err != nil {
		return err
	}

	return nil
}
