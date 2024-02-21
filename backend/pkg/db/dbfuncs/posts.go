package dbfuncs

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AddPost(post *Post) error {
	//may want to use autoincrement instead of uuids?
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	post.Id = id.String()
	post.CreatedAt = time.Now()
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(post.Id, post.Title, post.Body, post.CreatorId, post.GroupId, post.CreatedAt, post.Image, post.PrivacyLevel)

	return err
}

func DeletePost(PostId string) error {
	statement, err := db.Prepare("DELETE FROM Posts WHERE PostId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(PostId)
	return err
}

func AddPostChosenFollower(postChosenFollower *PostChosenFollower) error {
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(postChosenFollower.PostId, postChosenFollower.FollowerId)

	return err
}

func CountPostReacts(PostId string) (totalLikes, totalDislikes int, err error) {
	likes, dislikes, err := GetPostLikes(PostId)
	if err != nil {
		return
	}
	totalLikes = len(likes)
	totalDislikes = len(dislikes)
	return
}
func GetPostLikes(PostId string) (likeUserIds, dislikeUserIds []string, err error) {
	rows, err := db.Query("SELECT UserId, Liked, Disliked FROM CommentLikes WHERE PostId=?", PostId)
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

func isSupportedFileType() {

}

func SaveImage() {

}

// likeOrDislike can only take values "like" or "dislike"
func LikeDislikePost(UserId, PostId, likeOrDislike string) error {
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
	err := db.QueryRow("SELECT Liked, Disliked FROM PostLikes WHERE UserId=? AND PostId=?", UserId, PostId).Scan(&liked, &disliked)

	if err == sql.ErrNoRows {
		newRow, err := db.Prepare("INSERT INTO PostLikes VALUES (?,?,?,?)")
		if err != nil {
			return err
		}
		_, err = newRow.Exec(UserId, PostId, addLike, addDislike)
		return err
	}
	if err != nil {
		return err
	}
	if (liked && addLike) || (disliked && addDislike) {
		removeRow, err := db.Prepare("DELETE FROM PostLikes WHERE UserId=? AND PostId=?")
		if err != nil {
			return err
		}
		_, err = removeRow.Exec(UserId, PostId)
		return err
	}
	if (liked && addDislike) || (disliked && addLike) {
		updateRow, err := db.Prepare("UPDATE PostLikes SET Liked=?, Disliked=? WHERE UserId=? AND PostId=?")
		if err != nil {
			return err
		}
		_, err = updateRow.Exec(addLike, addDislike, UserId, PostId)
		return err
	}

	return errors.New("problem adding like or dislike: how did you get here?")
}

func SaveImageOld(file multipart.File, header *multipart.FileHeader) (string, error) {
	// generate new uuid for image name
	uniqueId := uuid.New()
	// remove "- from imageName"
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	// extract image extension from original file filename
	fileExt := strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]
	supported := isSupportedFileType(fileExt)

	if !supported {
		// rereturn "",error message to the user that this type of file is not supported
		return "", errors.New("this file type  is not supported")
	}

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)
	// create a new file in the "uploads" folder
	dst, err := os.Create(fmt.Sprintf("pkg/db/images/%s", image))
	if err != nil {
		log.Println("unable to create file --> ", err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	return image, nil
}

//should be unnecessary to do it this long way with ON CASCADE DELETE in db tables, but if not can still do it this way
// func RemovePostOld(postID string) error {
// 	db, err := sql.Open("sqlite3", "./forum.db")
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
// 	// delete comment likes and dislikes
// 	stmtCommentLikesDislikes, err := db.Prepare("DELETE FROM CommentLikes WHERE CommentId IN (SELECT id FROM 	Comments WHERE  PostId = ?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtCommentLikesDislikes.Close()

// 	_, err = stmtCommentLikesDislikes.Exec(postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete comments
// 	stmtComments, err := db.Prepare("DELETE FROM Comments WHERE PostId = ?")

// 	if err != nil {
// 		return err
// 	}
// 	defer stmtComments.Close()

// 	_, err = stmtComments.Exec(postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete likes
// 	stmtLikes, err := db.Prepare("DELETE FROM  Likes WHERE PostId= ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtLikes.Close()

// 	_, err = stmtLikes.Exec(postID)
// 	if err != nil {
// 		return err
// 	}

// 	//delete cats

// 	stmtPostCats, err := db.Prepare("DELETE FROM PostCat WHERE PostId= ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtPostCats.Close()

// 	_, err = stmtPostCats.Exec(postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete the post
// 	stmtPost, err := db.Prepare("DELETE FROM  Posts WHERE Id = ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtPost.Close()

// 	_, err = stmtPost.Exec(postID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func isSupportedFileTypeOld(fileType string) bool {

	supportedTypes := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
	}
	return supportedTypes[strings.ToLower(fileType)]
}

func GetPostChosenFollowerIdsByPostId(id string) ([]string, error) {
	var followerIds []string
	row, err := db.Query("SELECT FollowerId FROM PostChosenFollowers WHERE PostId=?", id)
	if err == sql.ErrNoRows {
		return followerIds, nil
	}
	if err != nil {
		return followerIds, err
	}
	defer row.Close()
	for row.Next() {
		var followerId string
		err = row.Scan(&followerId)
		if err != nil {
			return followerIds, err
		}
		followerIds = append(followerIds, followerId)
	}
	return followerIds, err
}

func GetPostByCommentId(id string) (Post, error) {
	var post Post
	err := db.QueryRow("SELECT Id, Title, Body, CreatorId, GroupId, CreatedAt, Image, PrivacyLevel FROM Posts WHERE Id=?", id).Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
	return post, err
}
