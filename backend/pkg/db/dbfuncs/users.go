package dbfuncs

import (
	"database/sql"
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
	statement.Exec(user.Id, user.Nickname, user.FirstName, user.LastName, user.Email, user.Password, user.Profile, user.AboutMe, user.PrivacySetting, user.DateOfBirth, user.CreatedAt)

	return nil
}

func AddSession(session *Session, sessionMins time.Duration) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	session.Id = id.String()
	session.Expires = time.Now().Add(time.Minute * sessionMins)
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?)")
	if err != nil {
		return err
	}
	statement.Exec(session.Id, session.Expires, session.UserId)

	return nil
}

func DeleteSession(sessionId string) error {
	statement, err := db.Prepare("DELETE FROM Sessions WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(sessionId)
	if err != nil {
		return err
	}
	return nil
}

// what is this for?!
// func DeleteSessionColumnOld(column string, value interface{}) error {
// 	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM Sessions WHERE %s = ?", column))
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(value)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }

func CheckEmailInDB(email string) (bool, error) {
	found := ""
	err := db.QueryRow("SELECT Email FROM Users WHERE Email=?", email).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GetAllUsersSortedByLastMessage() {

}

// func GetAllUsersSortedByLastMessageOld(id string) ([]User, error) {
// 	query := `
// 	SELECT
// 			u.Id,
// 			u.FirstName,
// 			u.LastName,
// 			u.Age,
// 			u.Gender,
// 			u.profileImg,
// 			u.Created,
// 			MAX(CASE
// 					WHEN m.SenderId = ? AND m.RecipientId = u.Id THEN datetime(m.Created)
// 					WHEN m.SenderId = u.Id AND m.RecipientId = ? THEN datetime(m.Created)
// 					ELSE NULL
// 			END) AS LastMessageTime
// 	FROM Users AS u
// 	LEFT JOIN Messages AS m ON (u.Id = m.SenderId AND m.RecipientId = ?) OR (u.Id = m.RecipientId AND m.SenderId = ?)
// 	GROUP BY
// 			u.Id,
// 			u.FirstName,
// 			u.LastName,
// 			u.Age,
// 			u.Gender,
// 			u.profileImg,
// 			u.Created
// `

// 	rows, err := database.Query(query, id, id, id, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var usersWithMessages []User
// 	var usersWithoutMessages []User

// 	for rows.Next() {
// 		var user User
// 		var lastMessageTimeString sql.NullString

// 		err := rows.Scan(
// 			&user.Id,
// 			&user.FirstName,
// 			&user.LastName,
// 			&user.Age,
// 			&user.Gender,
// 			&user.Img,
// 			&user.Created,
// 			&lastMessageTimeString,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		fmt.Println(lastMessageTimeString, user.FirstName, "usersWithoutMessages")
// 		if lastMessageTimeString.Valid {
// 			customTimeFormat := "2006-01-02 15:04:05"
// 			lastMessageTime, err := time.Parse(customTimeFormat, lastMessageTimeString.String)

// 			if err != nil {
// 				return nil, err

// 			}
// 			user.LastMessageTime = lastMessageTime
// 			usersWithMessages = append(usersWithMessages, user)
// 		} else {
// 			usersWithoutMessages = append(usersWithoutMessages, user)
// 		}

// 	}

// 	// Sort users by LastMessageTime in descending order
// 	sort.Slice(usersWithMessages, func(i, j int) bool {
// 		return usersWithMessages[i].LastMessageTime.After(usersWithMessages[j].LastMessageTime)
// 	})

// 	// Sort users alphabetically by FirstName and LastName for users without messages
// 	sort.Slice(usersWithoutMessages, func(i, j int) bool {
// 		nameI := strings.Title(usersWithoutMessages[i].LastName)
// 		nameJ := strings.Title(usersWithoutMessages[j].LastName)
// 		return nameI < nameJ
// 	})
// 	fmt.Println(usersWithoutMessages, "usersWithoutMessages")
// 	fmt.Println(usersWithMessages, "usersWithMessages")
// 	sortedUsers := append(usersWithMessages, usersWithoutMessages...)
// 	return sortedUsers, nil
// }

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

func ValidateCookie(sessionId string) (bool, error) {
	var id string
	var expiration time.Time
	err := db.QueryRow("SELECT Id, Expires FROM Sessions WHERE Id=?", sessionId).Scan(&id, &expiration)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return id == sessionId && time.Now().Before(expiration), nil
}
