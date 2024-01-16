package dbfuncs

import (
	"database/sql"
)

func RemovePost(postID string) error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	defer db.Close()
	// delete comment likes and dislikes
	stmtCommentLikesDislikes, err := db.Prepare("DELETE FROM CommentLikes WHERE CommentId IN (SELECT id FROM 	Comments WHERE  PostId = ?)")
	if err != nil {
		return err
	}
	defer stmtCommentLikesDislikes.Close()

	_, err = stmtCommentLikesDislikes.Exec(postID)
	if err != nil {
		return err
	}

	// Delete comments
	stmtComments, err := db.Prepare("DELETE FROM Comments WHERE PostId = ?")

	if err != nil {
		return err
	}
	defer stmtComments.Close()

	_, err = stmtComments.Exec(postID)
	if err != nil {
		return err
	}

	// Delete likes
	stmtLikes, err := db.Prepare("DELETE FROM  Likes WHERE PostId= ?")
	if err != nil {
		return err
	}
	defer stmtLikes.Close()

	_, err = stmtLikes.Exec(postID)
	if err != nil {
		return err
	}

	//delete cats

	stmtPostCats, err := db.Prepare("DELETE FROM PostCat WHERE PostId= ?")
	if err != nil {
		return err
	}
	defer stmtPostCats.Close()

	_, err = stmtPostCats.Exec(postID)
	if err != nil {
		return err
	}

	// Delete the post
	stmtPost, err := db.Prepare("DELETE FROM  Posts WHERE Id = ?")
	if err != nil {
		return err
	}
	defer stmtPost.Close()

	_, err = stmtPost.Exec(postID)
	if err != nil {
		return err
	}

	return nil
}
