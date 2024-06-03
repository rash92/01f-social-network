package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"database/sql"
	"fmt"
)

func FindPostsComment(id string) (Comment, error) {
	rows, err := database.Query("SELECT Id,Body, UserId,PostId,Created FROM Comments WHERE PostId=?", id)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return Comment{}, err
	} else if err != nil {
		return Comment{}, err
	}

	defer rows.Close()
	var one Comment
	var userId string
	for rows.Next() {
		err = rows.Scan(&one.Id, &one.Body, &userId, &one.PostID, &one.CreatedAt)
		if err != nil {
			return Comment{}, err
		}

		err = database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", userId).Scan(&one.CreatorNickname)

		if err != nil {
			// log.Fatal(err, "what")
			return Comment{}, err
		}

	}
	one.Likes, one.Dislikes, err = dbfuncs.CountCommentReacts(one.Id)
	return one, err
}
