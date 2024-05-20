package handlefuncs

import (
	"database/sql"
	"fmt"

	dbfuncs "server/pkg/db/dbfuncs"

	_ "github.com/mattn/go-sqlite3"
)

func FindAllPosts() ([]PostFontEnd, error) {
	rows, err := database.Query("SELECT Id,Title,Body,CreatorId,Created FROM Posts ORDER BY Created DESC")
	if err == sql.ErrNoRows {
		return []PostFontEnd{}, err
	} else if err != nil {
		return []PostFontEnd{}, err
	}
	var Everything []PostFontEnd

	defer rows.Close()

	for rows.Next() {
		var One PostFontEnd
		var id string
		rows.Scan(&One.Id, &One.Title, &One.Body, &id, &One.Created_at)

		err = database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", id).Scan(&One.Username)

		if err != nil {
			fmt.Println(err.Error())
			return []PostFontEnd{}, err
		}
		One.Categories = dbfuncs.FindPostsCats(One.Id)
		One.Comments, err = FindPostsComments(One.Id)
		if err != nil {
			fmt.Println("err, a")
			return []PostFontEnd{}, err
		}
		One.Likes, One.Dislikes = dbfuncs.CountLikesDislikes(One.Id)
		One.Userlikes = dbfuncs.FindLikeUsers(One.Id)

		Everything = append(Everything, One)
	}

	return Everything, nil
}
