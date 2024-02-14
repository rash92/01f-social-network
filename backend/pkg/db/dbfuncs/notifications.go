package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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

// updates notification to seen
func NotificationSeen(notificationId string) error {
	statement, err := db.Prepare("UPDATE Notifications SET Seen=? WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(true, notificationId)
	return err
}

func GetAllNotificationsByRecieverId(recieverId string) ([]Notification, error) {
	var notifications []Notification
	rows, err := db.Query("SELECT * FROM Notifications WHERE RecieverId=?", recieverId)
	if err == sql.ErrNoRows {
		return notifications, nil
	}
	if err != nil {
		return notifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen)
		if err != nil {
			return notifications, err
		}
		notifications = append(notifications, notification)
	}
	err = rows.Err()
	return notifications, err
}

func GetUnseenNotificationsByRecieverId(recieverId string) ([]Notification, error) {
	var unseenNotifications []Notification
	rows, err := db.Query("SELECT * FROM Notifications WHERE RecieverId=? AND Seen=?", recieverId, false)
	if err == sql.ErrNoRows {
		return unseenNotifications, nil
	}
	if err != nil {
		return unseenNotifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err = rows.Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen)
		if err != nil {
			return unseenNotifications, err
		}
		unseenNotifications = append(unseenNotifications, notification)
	}
	err = rows.Err()
	return unseenNotifications, err
}

func GetSeenNotificationsByRecieverId(recieverId string) ([]Notification, error) {
	var seenNotifications []Notification
	rows, err := db.Query("SELECT * FROM Notifications WHERE RecieverId=? AND Seen=?", recieverId, true)
	if err == sql.ErrNoRows {
		return seenNotifications, nil
	}
	if err != nil {
		return seenNotifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err = rows.Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen)
		if err != nil {
			return seenNotifications, err
		}
		seenNotifications = append(seenNotifications, notification)
	}
	err = rows.Err()
	return seenNotifications, err
}
