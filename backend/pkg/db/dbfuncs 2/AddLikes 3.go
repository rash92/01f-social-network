package dbfuncs

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddLikes(UserID, PostId string) error {
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
