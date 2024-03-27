package dbfuncs

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func LikeComment(UserID, CommentId string) {
	database, _ := sql.Open("sqlite3", "../sever/forum.db")
	newLike, _ := database.Prepare("INSERT INTO CommentLikes VALUES (?,?,?,?)")
	updateLike, _ := database.Prepare("UPDATE CommentLikes SET Liked=?, Disliked=? WHERE UserId=? AND CommentId=?")
	row := database.QueryRow("SELECT Liked, Disliked FROM CommentLikes WHERE UserId=? AND CommentId=?", UserID, CommentId)
	var like bool
	var dislike bool
	err := row.Scan(&like, &dislike)
	if err == sql.ErrNoRows {

		newLike.Exec(UserID, CommentId, true, false)
	} else if err != nil {

		log.Fatal(err)

	}
	if like {

		updateLike.Exec(false, false, UserID, CommentId)
	} else {

		updateLike.Exec(true, false, UserID, CommentId)
	}
}
