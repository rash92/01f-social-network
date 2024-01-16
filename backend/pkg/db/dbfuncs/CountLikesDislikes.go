package dbfuncs

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CountLikesDislikes(PostId string) (likes, dislikes int) {
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
