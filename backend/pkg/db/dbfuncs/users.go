package dbfuncs

import (
	"time"

	"github.com/google/uuid"
)

func AddUser(user *User) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user.Id = id.String()
	user.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO Comments VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.Id, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password, user.Profile, user.AboutMe, user.PrivacySetting, user.DOB, user.CreatedAt)

	return err
}

func AddSession(userId string) (Session, error) {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return Session{}, err
	}
	session := Session{
		Id:      id.String(),
		Expires: time.Now().Add(time.Minute * 60),
		UserId:  userId,
	}

	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?)")
	if err != nil {
		return Session{}, err
	}
	_, err = statement.Exec(session.Id, session.Expires, session.UserId)

	return session, err
}

func DeleteSession(sessionId string) error {
	statement, err := db.Prepare("DELETE FROM Sessions WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(sessionId)
	return err
}

func GetUserIdFromCookie(SessionId string) (string, error) {
	var userID string
	//if table has UserId starting lowercase, change table to be consistent.
	err := db.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", SessionId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetCreatorIdFromPostId(postId string) (string, error) {
	var creatorId string
	err := db.QueryRow("SELECT CreatorId FROM Posts WHERE Id=?", postId).Scan(&creatorId)
	if err != nil {
		return "", err
	}

	return creatorId, nil
}
