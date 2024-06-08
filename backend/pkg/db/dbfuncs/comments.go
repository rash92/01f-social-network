package dbfuncs

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// check if pointery way of doing it is working with * and & the right way etc., or if we want to just pass in by value
func AddComment(comment *Comment) (string, error) {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	comment.Id = id.String()
	comment.CreatedAt = time.Now()

	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("INSERT INTO Comments VALUES (?,?,?,?,?,?)")
	if err != nil {
		return "", err
	}
	_, err = statement.Exec(comment.Id, comment.Body, comment.CreatorId, comment.PostId, comment.CreatedAt, comment.Image)
	fmt.Println(comment.Id)
	return comment.Id, err
}

func DeleteComment(commentId string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	statement, err := db.Prepare("DELETE FROM Comments WHERE CommentId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(commentId)
	return err
}

// returns likes, dislikes, error
func CountCommentReacts(CommentId string) (totalLikes, totalDislikes int, err error) {
	likes, dislikes, err := GetCommentLikes(CommentId)
	if err != nil {
		return
	}
	totalLikes = len(likes)
	totalDislikes = len(dislikes)
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

// likeOrDislike can only take values "like" or "dislike"
func LikeDislikeComment(UserId, CommentId, likeOrDislike string) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	addLike := false
	addDislike := false
	if likeOrDislike == "like" {
		addLike = true
	} else if likeOrDislike == "dislike" {
		addDislike = true
	} else {
		return errors.New("like or dislike are the only options for parameter likeOrDislike")
	}

	var liked bool
	var disliked bool
	err := db.QueryRow("SELECT Liked, Disliked FROM CommentLikes WHERE UserId=? AND CommentId=?", UserId, CommentId).Scan(&liked, &disliked)

	if err == sql.ErrNoRows {
		newRow, err := db.Prepare("INSERT INTO CommentLikes VALUES (?,?,?,?)")
		if err != nil {
			return err
		}
		_, err = newRow.Exec(UserId, CommentId, addLike, addDislike)
		return err
	}
	if err != nil {
		return err
	}
	if (liked && addLike) || (disliked && addDislike) {
		removeRow, err := db.Prepare("DELETE FROM CommentLikes WHERE UserId=? AND CommentId=?")
		if err != nil {
			return err
		}
		_, err = removeRow.Exec(UserId, CommentId)
		return err
	}
	if (liked && addDislike) || (disliked && addLike) {
		updateRow, err := db.Prepare("UPDATE CommentLikes SET Liked=?, Disliked=? WHERE UserId=? AND CommentId=?")
		if err != nil {
			return err
		}
		_, err = updateRow.Exec(addLike, addDislike, UserId, CommentId)
		return err
	}

	return errors.New("problem adding like or dislike: how did you get here?")
}

func GetCommentById(id string) (Comment, error) {
	var comment Comment
	err := db.QueryRow("SELECT Id, Body, CreatorId, PostId, CreatedAt, Image FROM Comments WHERE Id=?", id).Scan(&comment.Id, &comment.Body, &comment.CreatorId, &comment.PostId, &comment.CreatedAt, &comment.Image)

	return comment, err
}

func GetAllCommentsByPostId(postId string) ([]Comment, error) {
	rows, err := db.Query("SELECT Id, Body, CreatorId, PostId, CreatedAt, Image FROM Comments WHERE PostId=?", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Id, &comment.Body, &comment.CreatorId, &comment.PostId, &comment.CreatedAt, &comment.Image)
		if err != nil {
			return nil, err
		}

		user, err := GetUserById(comment.CreatorId)
		if err != nil {
			return nil, err
		}

		comment.CreatorNickname = user.Nickname
		comments = append(comments, comment)
	}
	return comments, nil
}

// get number of comments by by a postId
// func GetNumberOfCommentsByPostId(postId string) (int, error) {
// 	var count int
// 	err := db.QueryRow("SELECT COUNT(*) FROM Comments WHERE PostId=?", postId).Scan(&count)
// 	return count, err
// }

//TO DO: get 10 at a time? decide if doing it through SQL or get all and do in handlefunc
