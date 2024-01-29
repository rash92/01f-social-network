package dbfuncs

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func AddDislikes(UserID, PostId string) error {
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
