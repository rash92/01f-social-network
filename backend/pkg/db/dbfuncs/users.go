package dbfuncs

import (
	"database/sql"
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
	_, err = statement.Exec(user.Id, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password, user.Avatar, user.AboutMe, user.PrivacySetting, user.DOB, user.CreatedAt)

	return err
}

func DeleteUser(userId string) error {
	statement, err := db.Prepare("DELETE FROM Users WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(userId)
	return err
}

func AddSession(userId string) (Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Session{}, err
	}
	session := Session{
		Id:      id.String(),
		Expires: time.Now().Add(time.Hour * 72),
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
	if privacySetting == "private" || privacySetting == "" {
	if privacySetting == "private" || privacySetting == "" {
		return true, nil
	}
	return false, errors.New("privacy setting not recognized, should be either 'private' or 'public'")
}

func GetUserById(id string) (User, error) {
	var user User
	err := db.QueryRow("SELECT Id, Nickname, FirstName, LastName, Email, Password,Avatar, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users WHERE Id=?", id).Scan(&user.Id, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Avatar, &user.AboutMe, &user.PrivacySetting, &user.DOB, &user.CreatedAt)

	return user, err
}

//figure out whether to delete or rewrite or keep etc.

func GetUsers() ([]User, error) {
	rows, err := db.Query("SELECT Id,FirstName, LastName, Nickname, Avatar, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	var user []User

	for rows.Next() {
		var newUser User
		err := rows.Scan(&newUser.Id, &newUser.FirstName, &newUser.LastName, &newUser.Nickname, &newUser.Avatar, &newUser.AboutMe, &newUser.PrivacySetting, &newUser.DOB, &newUser.CreatedAt)
		if err != nil {
			return []User{}, err
		}
		user = append(user, newUser)
	}

	return user, err
}

// togle PrivacySetting
func UpdatePrivacySetting(userId string, privacySetting string) error {
	
	statement, err := db.Prepare("UPDATE Users SET PrivacySetting=? WHERE Id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(privacySetting, userId)
	if err != nil {
		return err
	}
	return nil
}

// possibly not necessary and can do this in helper where needed
func GetNicknameFromId(userId string) (string, error) {
	user, err := GetUserById(userId)

	return user.Nickname, err
}

func GetUserBySessionId(sessionId string) (User, error) {
	userId, err := GetUserIdFromCookie(sessionId)
	if err != nil {
		return User{}, err
	}
	user, err := GetUserById(userId)
	return user, err
}

// returns slice of users with only id, nickname and avatar filled, everything else is default/ nil values
func SearchUsers(query string) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT Id,  Nickname, Avatar  FROM users WHERE FirstName LIKE ? OR LastName LIKE ? OR Nickname LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err == sql.ErrNoRows {
		return users, nil
	}
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Avatar); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	return users, err
}
