package handlefuncs

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func FindPostsComment(id string) (Comment, error) {
// 	rows, err := database.Query("SELECT Id,Body, UserId,PostId,Created FROM Comments WHERE PostId=?", id)
// 	if err == sql.ErrNoRows {
// 		fmt.Println(err)
// 		return Comment{}, err
// 	} else if err != nil {
// 		return Comment{}, err
// 	}

// 	defer rows.Close()
// 	var one Comment
// 	var userId string
// 	for rows.Next() {
// 		err = rows.Scan(&one.ID, &one.Body, &userId, &one.PostID, &one.CreatedAt)
// 		if err != nil {
// 			return Comment{}, err
// 		}

// 		err = database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", userId).Scan(&one.Username)

// 		if err != nil {
// 			// log.Fatal(err, "what")
// 			return Comment{}, err
// 		}

// 	}
// 	// one.Likes, one.Dislikes = dbfuncs.CountCommentReacts(one.ID)
// 	return one, nil
// }
