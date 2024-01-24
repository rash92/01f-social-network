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
	statement.Exec(post.Id, post.Title, post.Body, post.CreatorId, post.GroupId, post.CreatedAt, post.Image, post.PrivacyLevel)

	return nil
}

func AddPostChosenFollower(postChosenFollower *PostChosenFollower) error {
	statement, err := db.Prepare("INSERT INTO groups VALUES (?,?)")
	if err != nil {
		return err
	}
	statement.Exec(postChosenFollower.PostId, postChosenFollower.FollowerId)

	return nil
}

func CountLikesDislikes() {

}

func isSupportedFileType() {

}

func RemovePost() {

}

func SaveImage() {

}

func AddDislikes() {

}

func AddLikes() {

}

func FindLikeUsers() {

}

func FindPostsCatsOld(PostId string) []string {
	rows, err := database.Query("SELECT CatId FROM PostCat WHERE PostId=?", PostId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var AllCats []string

	for rows.Next() {
		var CatId uuid.UUID
		err := rows.Scan(&CatId)
		if err != nil {
			log.Fatal(err)
		}
		name := database.QueryRow("SELECT Name FROM Categories WHERE Id=?", CatId)
		var cat string
		err = name.Scan(&cat)
		if err != nil {
			log.Fatal(err)
		}
		AllCats = append(AllCats, cat)
	}
	return AllCats
}

func FindLikeUsersOld(PostId string) []string {
	rows, err := database.Query("SELECT UserId FROM Likes WHERE PostId=? AND Liked=1", PostId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var AllLikes []string

	for rows.Next() {
		var UserId uuid.UUID
		err := rows.Scan(&UserId)
		if err != nil {
			log.Fatal(err)
		}
		name := database.QueryRow("SELECT Nickname FROM Users WHERE Id=?", UserId)
		var user string
		err = name.Scan(&user)
		if err != nil {
			log.Fatal(err)
		}
		AllLikes = append(AllLikes, user)
	}
	return AllLikes
}

func AddLikesOld(UserID, PostId string) error {
	newLike, err := database.Prepare("INSERT INTO Likes VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	updateLike, err := database.Prepare("UPDATE Likes SET Liked=?, Disliked=? WHERE UserId=? AND PostId=?")
	if err != nil {
		return err
	}
	row := database.QueryRow("SELECT Liked, Disliked FROM Likes WHERE UserId=? AND PostId=?", UserID, PostId)
	var like bool
	var dislike bool
	err = row.Scan(&like, &dislike)

	if err == sql.ErrNoRows {
		newLike.Exec(UserID, PostId, true, false)
	} else if err != nil {
		log.Fatal(err)
		return err

	}
	if like {

		updateLike.Exec(false, false, UserID, PostId)
	} else {

		updateLike.Exec(true, false, UserID, PostId)
	}

	return nil
}

func AddDislikesOld(UserID, PostId string) error {
	newDislike, err := database.Prepare("INSERT INTO Likes VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	updateDislike, err := database.Prepare("UPDATE Likes SET Liked=?, Disliked=? WHERE UserId=? AND PostId=?")
	if err != nil {
		return err
	}
	row := database.QueryRow("SELECT Liked, Disliked FROM Likes WHERE UserId=? AND PostId=?", UserID, PostId)
	var like bool
	var dislike bool
	err = row.Scan(&like, &dislike)
	if err == sql.ErrNoRows {
		newDislike.Exec(UserID, PostId, false, true)
	} else if err != nil {
		return err

	}
	if dislike {

		updateDislike.Exec(false, false, UserID, PostId)
	} else {

		updateDislike.Exec(false, true, UserID, PostId)
	}

	return nil
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

func RemovePostOld(postID string) error {
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

func isSupportedFileTypeOld(fileType string) bool {

	supportedTypes := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
	}
	return supportedTypes[strings.ToLower(fileType)]
}

func CountLikesDislikesOld(PostId string) (likes, dislikes int) {
	rows, err := database.Query("SELECT Liked, Disliked FROM Likes WHERE PostId=?", PostId)
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
