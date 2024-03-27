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
	_, err = statement.Exec(user.Id, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password, user.Avatar, user.AboutMe, user.PrivacySetting, user.DOB, user.CreatedAt)

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
	err := db.QueryRow("SELECT Id, Nickname, FirstName, LastName, Email, Password,Avatar, AboutMe, PrivacySetting, DOB, CreatedAt FROM Users WHERE Id=?", id).Scan(&user.Id, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Avatar, &user.AboutMe, &user.PrivacySetting, &user.DOB, &user.CreatedAt)

	return user, err
}

//figure out whether to delete or rewrite or keep etc.

func Getusers() ([]User, error) {
	rows, err := db.Query("SELECT Id,FirstName, LastName, Nickname,Avatar, AboutMe, Privacy_setting, DOB, CreatedAt FROM Users")
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

// func GetUserDataFromSession(sessionId string) (string, string, string, error) {
// 	var userId string
// 	varAvatarImage string
// 	var nickname string

// 	// Execute the SQL query
// 	err := database.QueryRow(`
// 			SELECT Sessions.userId, Users.Profile, Users.Nickname
// 			FROM Sessions
// 			JOIN Users ON Sessions.UserID = Users.Id
// 			WHERE Sessions.Id = ?
// 	`, sessionId).Scan(&userId, &profileImage, &nickname)

// 	// Check for errors
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	return userId,AvatarImage, nickname, nil
// }

// func GetNumberOfByUserId(userId string, table string) (int, error) {
// 	var count int
// 	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE Creatorid=?", table)
// 	err := database.QueryRow(query, userId).Scan(&count)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	return count, nil
// }
