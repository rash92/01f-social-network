package dbfuncs

import (
	"database/sql"
	"errors"
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
	statement, err := db.Prepare("INSERT INTO Posts VALUES (?,?,?,?,?,?,?,?)")
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
	statement, err := db.Prepare("INSERT INTO PostChosenFollowers VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(postChosenFollower.PostId, postChosenFollower.FollowerId)

	return err
}

func DeletePostChosenFollower(postChosenFollower *PostChosenFollower) error {
	statement, err := db.Prepare("DELETE FROM PostChosenFollowers WHERE PostId=? AND FollowerId=?")
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
	rows, err := db.Query("SELECT UserId, Liked, Disliked FROM PostLikes WHERE PostId=?", PostId)
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

func GetPostById(id string) (Post, error) {
	var post Post
	err := db.QueryRow("SELECT Id, Title, Body, CreatorId, GroupId, CreatedAt, Image, PrivacyLevel FROM Posts WHERE Id=?", id).Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
	return post, err
}

func GetPostIdByCommentId(commentId string) (string, error) {
	var postId string
	err := db.QueryRow("SELECT PostId FROM Comments WHERE Id=?", commentId).Scan(&postId)
	return postId, err
}

func GetPostsByCreatorId(creatorId string) ([]Post, error) {
	var posts []Post
	rows, err := db.Query("SELECT * FROM Posts WHERE CreatorId=?", creatorId)
	if err == sql.ErrNoRows {
		return posts, nil
	}
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, err
}

//TO DO: get 10 at a time? decide if doing it through SQL or get all and do in handlefunc
