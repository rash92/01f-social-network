package dbfuncs

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func AddSession() {

}

func AddSessionOld(id uuid.UUID, user, UserID string, Expires time.Time) {
	statement, _ := database.Prepare("INSERT INTO sessions VALUES (?,?,?,?)")

	statement.Exec(id, user, Expires, UserID)
}

func DeleteSessionColumn() {

}

func DeleteSessionColumnOld(column string, value interface{}) error {
	stmt, err := database.Prepare(fmt.Sprintf("DELETE FROM Sessions WHERE %s = ?", column))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(value)
	if err != nil {
		return err
	}
	return nil

}

func AddUser() {

}

func AddUserOld(nickName, firstName, lastName, Email, profile, aboutMe, privacy, DOB string, Password []byte) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO Users VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	statement.Exec(id, nickName, firstName, lastName, Email, Password, profile, aboutMe, privacy, DOB, created)

	return nil
}

func CheckValueInDB() {

}

func CheckValueInDBOld(w http.ResponseWriter, r *http.Request, val, name string) (string, bool, error) {
	if name != "Nickname" && name != "Email" {
		return "Invalid column name", false, nil
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s = ?", name)
	err := database.QueryRow(query, val).Scan(&count)
	if err != nil {
		return "Error querying database", false, err
	}
	fmt.Println(count > 0, name)
	return "", count > 0, nil
}

func GetAllUsersSortedByLastMessage() {

}

func GetAllUsersSortedByLastMessageOld(id string) ([]User, error) {
	query := `
	SELECT
			u.Id,
			u.FirstName,
			u.LastName,
			u.Age,
			u.Gender,
			u.profileImg,
			u.Created,
			MAX(CASE
					WHEN m.SenderId = ? AND m.RecipientId = u.Id THEN datetime(m.Created)
					WHEN m.SenderId = u.Id AND m.RecipientId = ? THEN datetime(m.Created)
					ELSE NULL
			END) AS LastMessageTime
	FROM Users AS u
	LEFT JOIN Messages AS m ON (u.Id = m.SenderId AND m.RecipientId = ?) OR (u.Id = m.RecipientId AND m.SenderId = ?)
	GROUP BY
			u.Id,
			u.FirstName,
			u.LastName,
			u.Age,
			u.Gender,
			u.profileImg,
			u.Created
`

	rows, err := database.Query(query, id, id, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersWithMessages []User
	var usersWithoutMessages []User

	for rows.Next() {
		var user User
		var lastMessageTimeString sql.NullString

		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.Gender,
			&user.Img,
			&user.Created,
			&lastMessageTimeString,
		)
		if err != nil {
			return nil, err
		}

		fmt.Println(lastMessageTimeString, user.FirstName, "usersWithoutMessages")
		if lastMessageTimeString.Valid {
			customTimeFormat := "2006-01-02 15:04:05"
			lastMessageTime, err := time.Parse(customTimeFormat, lastMessageTimeString.String)

			if err != nil {
				return nil, err

			}
			user.LastMessageTime = lastMessageTime
			usersWithMessages = append(usersWithMessages, user)
		} else {
			usersWithoutMessages = append(usersWithoutMessages, user)
		}

	}

	// Sort users by LastMessageTime in descending order
	sort.Slice(usersWithMessages, func(i, j int) bool {
		return usersWithMessages[i].LastMessageTime.After(usersWithMessages[j].LastMessageTime)
	})

	// Sort users alphabetically by FirstName and LastName for users without messages
	sort.Slice(usersWithoutMessages, func(i, j int) bool {
		nameI := strings.Title(usersWithoutMessages[i].LastName)
		nameJ := strings.Title(usersWithoutMessages[j].LastName)
		return nameI < nameJ
	})
	fmt.Println(usersWithoutMessages, "usersWithoutMessages")
	fmt.Println(usersWithMessages, "usersWithMessages")
	sortedUsers := append(usersWithMessages, usersWithoutMessages...)
	return sortedUsers, nil
}

func GetUserIdFromCokie() {

}

func GetUserIdFromCokieOld(sessionId string) (string, error) {
	var userID string
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", sessionId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetUserIdFromPostId() {

}

func GetUserIdFromPostIdOld(postId string) (string, error) {
	var userID string
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.QueryRow("SELECT UserId  FROM Posts  WHERE Id=?", postId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil

}

func ValidateCookie() {

}

func ValidateCookieOld(cookieValue string) bool {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return false
	}
	defer db.Close()
	var id string
	var expiration time.Time
	err = db.QueryRow("SELECT Id, expires  FROM Sessions WHERE Id=?", cookieValue).Scan(&id, &expiration)
	if err != nil {
		return false
	}

	return id == cookieValue && !(time.Now().After(expiration))

}
