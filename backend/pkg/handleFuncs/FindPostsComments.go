package handlefuncs

import (
	"database/sql"
	"fmt"
	"log"
	dbfuncs "server/pkg/db/dbfuncs"

	_ "github.com/mattn/go-sqlite3"
)

func FindPostsComments(PostId string) ([]Comment, error) {
	rows, err := database.Query("SELECT Id,Body,UserId,PostId,Created FROM  Comments WHERE PostId=? ORDER BY Created DESC", PostId)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return []Comment{}, err
	} else if err != nil {
		log.Fatal(err)
		return []Comment{}, err

	}

	defer rows.Close()
	var AllComments []Comment
	var userId string
	for rows.Next() {
		var one Comment
		err = rows.Scan(&one.ID, &one.Body, &userId, &one.PostID, &one.CreatedAt)
		if err != nil {
			fmt.Println(err.Error())
			return []Comment{}, err
		}

		err = database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", userId).Scan(&one.Username)

		if err != nil {
			log.Fatal(err, "what")
			return []Comment{}, err
		}
		one.Likes, one.Dislikes = dbfuncs.CountCommentReacts(one.ID)
		AllComments = append(AllComments, one)
	}

	return AllComments, nil
}
