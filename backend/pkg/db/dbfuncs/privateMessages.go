package dbfuncs

import (
	"time"

	"github.com/google/uuid"
)

func AddPrivateMessage(privateMessage *PrivateMessage) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	privateMessage.Id = id.String()
	privateMessage.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	statement.Exec(privateMessage.Id, privateMessage.SenderId, privateMessage.RecipientId, privateMessage.Message, privateMessage.CreatedAt)

	return nil
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
