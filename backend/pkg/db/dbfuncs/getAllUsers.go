package dbfuncs

import (
	"time"

	"github.com/google/uuid"
)

type User_getAllUsers struct {
	Id              uuid.UUID
	FirstName       string
	LastName        string
	Nickname        string
	Profile         string
	AboutMe         string
	Privacy_setting string
	DOB             string
	CreatedAt       time.Time
	// LastMessageTime time.Time
}

// func GetAllUsersSortedByLastMessage(id string) ([]User_getAllUsers, error) {
// 	query := `
// 	SELECT
// 			u.Id,
// 			Nickname,
// 			u.FirstName,
// 			u.LastName,
// 			u.Profile,
// 			AboutMe,
// 			Privacy_setting,
// 			u.DOB,
// 			u.CreatedAt,
// 			MAX(CASE
// 					WHEN m.SenderId = ? AND m.RecipientId = u.Id THEN datetime(m.Created)
// 					WHEN m.SenderId = u.Id AND m.RecipientId = ? THEN datetime(m.Created)
// 					ELSE NULL
// 			END) AS LastMessageTime
// 	FROM Users AS u
// 	LEFT JOIN Messages AS m ON (u.Id = m.SenderId AND m.RecipientId = ?) OR (u.Id = m.RecipientId AND m.SenderId = ?)
// 	GROUP BY
// 	u.Id,
// 			Nickname,
// 			u.FirstName,
// 			u.LastName,
// 			u.Profile,
// 			AboutMe,
// 			Privacy_setting,
// 			u.DOB,
// 			u.CreatedAt,
// `

// 	rows, err := database.Query(query, id, id, id, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var usersWithMessages []User_getAllUsers
// 	var usersWithoutMessages []User_getAllUsers

// 	for rows.Next() {
// 		var user User_getAllUsers
// 		var lastMessageTimeString sql.NullString
// 		// u.FirstName,
// 		// u.LastName,
// 		// u.Profile,
// 		// AboutMe,
// 		// Privacy_setting
// 		// u.DOB,
// 		// u.CreatedAt,
// 		err := rows.Scan(
// 			&user.Id,
// 			&user.FirstName,
// 			&user.LastName,
// 			&user.Profile,
// 			&user.AboutMe,
// 			&user.Privacy_setting,
// 			&user.DOB,
// 			&user.CreatedAt,
// 			// &lastMessageTimeString,
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
