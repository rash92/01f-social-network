package dbfuncs

import "time"

func AddEvent(event *Event) (error, string, time.Time) {
	var err error
	return err, "", time.Now()
}

func GetEventsById(id string) ([]Event, error) {
	var events []Event
	var err error
	return events, err
}
