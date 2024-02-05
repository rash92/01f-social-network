package dbfuncs

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// check if pointery way of doing it is working with * and & the right way etc., or if we want to just pass in by value
func AddComment(comment *Comment) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	comment.Id = id.String()
	comment.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO Comments VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	statement.Exec(comment.Id, comment.Body, comment.CreatorId, comment.PostId, comment.CreatedAt, comment.Image)

	return nil
}

// returns likes, dislikes, error
func CountCommentReacts(CommentId string) (totalLikes, totalDislikes int, err error) {
	rows, err := db.Query("SELECT Liked, Disliked FROM CommentLikes WHERE CommentId=?", CommentId)
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err != nil {
		return
	}
	defer rows.Close()
	var like bool
	var dislike bool
	for rows.Next() {

		err = rows.Scan(&like, &dislike)
		if err != nil {
			return
		}
		if like {
			totalLikes++
		}
		if dislike {
			totalDislikes++
		}
	}
	return
}

func GetCommentLikes(CommentId string) (likeUserIds, dislikeUserIds []string, err error) {
	rows, err := db.Query("SELECT UserId, Liked, Disliked FROM CommentLikes WHERE CommentId=?", CommentId)
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err != nil {
		return
	}
	defer rows.Close()
	var userId string
	var like bool
	var dislike bool
	for rows.Next() {
		err = rows.Scan(&userId, &like, &dislike)
		if err != nil {
			return
		}
		if like {
			likeUserIds = append(likeUserIds, userId)
		}
		if dislike {
			dislikeUserIds = append(dislikeUserIds, userId)
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
