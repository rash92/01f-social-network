package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func AddNotification(notification *Notification) (string, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	notification.Id = id.String()
	notification.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO Notifications VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return "", err
	}
	_, err = statement.Exec(notification.Id, notification.Body, notification.Type, notification.CreatedAt, notification.ReceiverId, notification.SenderId, notification.Seen)

	return notification.Id, err
}

// updates notification to seen
func NotificationSeen(notificationId string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("UPDATE Notifications SET Seen=? WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(true, notificationId)
	return err
}

func GetAllNotificationsByRecieverId(recieverId string) ([]Notification, error) {
	var notifications []Notification
	rows, err := db.Query("SELECT n.*, u.Avatar FROM Notifications n JOIN Users u ON n.SenderId = u.Id WHERE n.ReceiverId=? ORDER BY n.CreatedAt DESC", recieverId)
	if err == sql.ErrNoRows {
		return notifications, nil
	}
	if err != nil {
		return notifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen, &notification.SenderAvatar)
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

func GetNotificationById(id string) (Notification, error) {
	var notification Notification
	err := db.QueryRow("SELECT Id, Body, Type, CreatedAt, ReceiverId, SenderId, Seen FROM Notifications WHERE Id=?", id).Scan(&notification.Id, &notification.Body, &notification.Type, &notification.CreatedAt, &notification.ReceiverId, &notification.SenderId, &notification.Seen)

	return notification, err
}

//TO DO: get 10 at a time? decide if doing it through SQL or get all and do in handlefunc
