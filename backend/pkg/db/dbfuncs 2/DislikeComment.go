package dbfuncs

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DislikeComment(UserID, CommentId string) {
	newDislike, _ := database.Prepare("INSERT INTO CommentLikes VALUES (?,?,?,?)")
	updateDislike, _ := database.Prepare("UPDATE CommentLikes SET Liked=?, Disliked=? WHERE UserId=? AND CommentId=?")
	row := database.QueryRow("SELECT Liked, Disliked FROM CommentLikes WHERE UserId=? AND CommentId=?", UserID, CommentId)
	var like bool
	var dislike bool
	err := row.Scan(&like, &dislike)
	if err == sql.ErrNoRows {

		newDislike.Exec(UserID, CommentId, false, true)
	} else if err != nil {

		log.Fatal(err)

	}
	if dislike {

		updateDislike.Exec(false, false, UserID, CommentId)
	} else {

		updateDislike.Exec(false, true, UserID, CommentId)
	}
}
