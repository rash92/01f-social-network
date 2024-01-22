package dbfuncs

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddComment() {

}

func AddCommentOld(Comment, SessionId string, PostId string) (string, error) {
	// fmt.Println(Comment, SessionId, PostId)
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	created := time.Now()
	statement, err := database.Prepare("INSERT INTO Comments VALUES (?,?,?,?,?)")
	if err != nil {
		return "", err
	}
	var UserId uuid.UUID
	err = database.QueryRow("SELECT  userId  FROM Sessions WHERE Id=?", SessionId).Scan(&UserId)

	if err != nil {

		return "", err
	}

	statement.Exec(id, Comment, UserId, PostId, created)

	return id.String(), nil
}

func CountCommentReacts() {

}

func CountCommentReactsOld(CommentId string) (likes, dislikes int) {
	rows, err := database.Query("SELECT Liked, Disliked FROM CommentLikes WHERE CommentId=?", CommentId)
	if err == sql.ErrNoRows {
		return 0, 0
	} else if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	likes = 0
	dislikes = 0
	var l bool
	var d bool
	for rows.Next() {
		err := rows.Scan(&l, &d)
		if err != nil {
			log.Fatal(err)
		}
		if l {
			likes++
		}
		if d {
			dislikes++
		}
	}

	return
}

func DislikeComment() {

}

func DislikeCommentOld(UserID, CommentId string) {
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

func LikedComment() {

}

func LikeCommentOld(UserID, CommentId string) {
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
