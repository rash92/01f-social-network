package handlefuncs

// import (
// 	"database/sql"
// 	"fmt"
// 	dbfuncs "server/pkg/db/dbfuncs"
// )

// func FindPostByID(id string) (PostFontEnd, error) {
// 	var One PostFontEnd
// 	var userId string

// 	err := database.QueryRow("SELECT Id,Title,Body,UserId,Created FROM Posts WHERE Id=?", id).Scan(&One.Id, &One.Title, &One.Body, &userId, &One.Created_at)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return PostFontEnd{}, fmt.Errorf("Post with ID %s not found", id)
// 		}
// 		return PostFontEnd{}, err
// 	}

// 	err = database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", userId).Scan(&One.Username)
// 	if err != nil {

// 		fmt.Println(err.Error())
// 		return PostFontEnd{}, err
// 	}

// 	// One.Categories = dbfuncs.FindPostsCats(One.Id)

// 	One.Comments, err = FindPostsComments(One.Id)
// 	if err != nil {
// 		fmt.Println("err, a")
// 		return PostFontEnd{}, err
// 	}

// 	// One.Likes, One.Dislikes = dbfuncs.CountLikesDislikes(One.Id)

// 	return One, nil
// }