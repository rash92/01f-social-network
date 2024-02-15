package dbfuncs

import (
	"time"

	"github.com/google/uuid"
)

// Assume everything here is a placeholder unless it's clear that it's what you want.

func IsUserPrivate(id string) (bool, error) {
	var privacySetting string
	err := db.QueryRow("SELECT PrivacySetting FROM Users WHERE id = $1", id).Scan(&privacySetting)
	if err != nil {
		return false, err
	}
	return privacySetting == "private", nil
}

func AddFollow(follow *Follow) error {
	_, err := db.Exec("INSERT INTO followers (FollowerId, FollowingId) VALUES ($1, $2)", follow.FollowerId, follow.FollowingId)
	return err
}

func AddNotification(notification *Notification) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	notification.Id = id.String()
	notification.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(notification.Id, notification.Body, notification.Type, notification.CreatedAt, notification.ReceiverId, notification.SenderId, notification.Seen)
	return err
}

func GetNotificationById(id string) (Notification, error) {
	var notification Notification
	return notification, nil
}

func AcceptFollow(followId string, followingId string) error {
	return nil
}

func RejectFollow(followId string, followingId string) error {
	return nil
}
