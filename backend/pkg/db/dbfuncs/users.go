package dbfuncs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

func AddUser(user *User) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user.Id = id.String()
	user.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO Users VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.Id, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password, user.Profile, user.AboutMe, user.PrivacySetting, user.DOB, user.CreatedAt)

	return err
}

func AddSession(userId string) (Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Session{}, err
	}
	session := Session{
		Id:      id.String(),
		Expires: time.Now().Add(time.Minute * 60),
		UserId:  userId,
	}

	statement, err := db.Prepare("INSERT INTO Sessions VALUES (?,?,?)")
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

	return userID, err
}

func GetCreatorIdFromPostId(postId string) (string, error) {
	var creatorId string
	err := db.QueryRow("SELECT CreatorId FROM Posts WHERE Id=?", postId).Scan(&creatorId)

	return creatorId, err
}

func IsUserPrivate(userId string) (bool, error) {
	var privacySetting string
	err := db.QueryRow("SELECT PrivacySetting FROM Users WHERE Id=?", userId).Scan(&privacySetting)
	if err != nil {
		return false, err
	}
	if privacySetting == "public" {
		return false, nil
	}
	if privacySetting == "private" {
		return true, nil
	}
	return false, errors.New("privacy setting not recognized, should be either 'private' or 'public'")
}

func GetUserById(id string) (User, error) {
	var user User
	err := db.QueryRow("SELECT Id, Nickname, FirstName, LastName, Email, Password, Profile, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users WHERE Id=?", id).Scan(&user.Id, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Profile, &user.AboutMe, &user.PrivacySetting, &user.DOB, &user.CreatedAt)

	return user, err
}
