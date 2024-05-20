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

func GetAllPostsByCreatorId(creatorId string) ([]Post, error) {
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

	err = rows.Err()

	return posts, err
}

func GetVisiblePosts(userId string) ([]Post, error) {
	query := `
	SELECT * FROM Posts
	WHERE 
			(PrivacyLevel = 'public') OR 
			(PrivacyLevel = 'private' AND CreatorId IN (SELECT FollowingId FROM Follows WHERE FollowerId = ?)) OR 
			(PrivacyLevel = 'superprivate' AND Id IN (SELECT PostId FROM PostChosenFollowers WHERE FollowerId = ?)) OR
			CreatorId = ?
	`
	rows, err := db.Query(query, userId, userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return nil, err
		}

		post.Dislikes, post.Likes, err = CountPostReacts(post.Id)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(post.CreatorId)
		if err != nil {
			return nil, err
		}
		post.CreatorNickname = user.Nickname
		post.UserLikeDislike, err = GetUserLikeDislike(userId, post.Id)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetVisiblePostsForProfile(userId, profileOwnerId string) ([]Post, error) {
	query := `
	SELECT * FROM Posts
	WHERE 
			(CreatorId = ?) AND
			((PrivacyLevel = 'public') OR 
			(PrivacyLevel = 'private' AND CreatorId IN (SELECT FollowingId FROM Follows WHERE FollowerId = ?)) OR 
			(PrivacyLevel = 'superprivate' AND Id IN (SELECT PostId FROM PostChosenFollowers WHERE FollowerId = ?)))
	`
	rows, err := db.Query(query, profileOwnerId, userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetUserLikeDislike(userId, postId string) (int, error) {
	var like int
	err := db.QueryRow("SELECT Liked FROM PostLikes WHERE UserId=? AND PostId=?", userId, postId).Scan(&like)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return like, nil
}

func GetNumberOfPostsByUserId(userId string) (int, error) {
	posts, err := GetAllPostsByCreatorId(userId)

	return len(posts), err
}

func GetProfilePosts(viewerId string, creatorId string) ([]Post, error) {
	var posts []Post
	// check if creator has public profile, if yes move on to checking privacy level of posts. If no, check if viewer is following them.
	creatorPrivate, err := IsUserPrivate(creatorId)
	if err != nil {
		return posts, err
	}
	isFollowing, err := IsFollowing(viewerId, creatorId)
	if err != nil {
		return posts, err
	}
	if creatorPrivate && !isFollowing {
		return posts, nil
	}

	row, err := db.Query(`
		SELECT * FROM Posts 
			WHERE CreatorId=? AND PrivacyLevel='public'
		UNION
		SELECT * FROM Posts
			WHERE CreatorId=? AND PrivacyLevel='private' AND CreatorId IN
			(
				SELECT FollowingId
					FROM Follows
					WHERE FollowerId=? AND Status='accepted'
			)
		UNION
		SELECT * FROM Posts
			WHERE CreatorId=? AND PrivacyLevel='superPrivate' AND Id IN
			(
				SELECT PostId FROM PostChosenFollowers
					WHERE FollowerId=?
			)
		ORDER BY CreatedAt DESC`,
		creatorId, creatorId, viewerId, creatorId, viewerId)

	if err == sql.ErrNoRows {
		return posts, nil
	}
	if err != nil {
		return posts, err
	}
	defer row.Close()

	for row.Next() {
		var post Post
		err = row.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	err = row.Err()

	return posts, err
}

func GetDashboardPosts(userId string) ([]Post, error) {
	var posts []Post

	row, err := db.Query(`
		SELECT * FROM Posts 
			WHERE PrivacyLevel='public'
		UNION
		SELECT * FROM Posts
			WHERE PrivacyLevel='private' AND CreatorId IN
			(
				SELECT FollowingId
					FROM Follows
					WHERE FollowerId=? AND Status='accepted'
			)
		UNION
		SELECT * FROM Posts
			WHERE PrivacyLevel='superPrivate' AND Id IN
			(
				SELECT PostId FROM PostChosenFollowers
					WHERE FollowerId=?
			)
		ORDER BY CreatedAt DESC`,
		userId, userId)

	if err == sql.ErrNoRows {
		return posts, nil
	}
	if err != nil {
		return posts, err
	}

	defer row.Close()
	for row.Next() {
		var post Post
		err = row.Scan(&post.Id, &post.Title, &post.Body, &post.CreatorId, &post.GroupId, &post.CreatedAt, &post.Image, &post.PrivacyLevel)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	err = row.Err()

	return posts, err
}
