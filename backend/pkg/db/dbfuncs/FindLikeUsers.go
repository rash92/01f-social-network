package dbfuncs

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func FindLikeUsers(PostId string) []string {
	rows, err := database.Query("SELECT UserId FROM Likes WHERE PostId=? AND Liked=1", PostId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var AllLikes []string

	for rows.Next() {
		var UserId uuid.UUID
		err := rows.Scan(&UserId)
		if err != nil {
			log.Fatal(err)
		}
		name := database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", UserId)
		var user string
		err = name.Scan(&user)
		if err != nil {
			log.Fatal(err)
		}
		AllLikes = append(AllLikes, user)
	}
	return AllLikes
}
