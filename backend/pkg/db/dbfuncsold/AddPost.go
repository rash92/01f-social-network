package dbfuncs

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func AddPost(cookieVal,  PostTitle, PostBody string, categories []string)( string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	created := time.Now()
	statement, err := database.Prepare("INSERT INTO Posts VALUES (?,?,?,?,?)")

	if err != nil {
		return "", err
	}
	var UserId uuid.UUID
	err = database.QueryRow("SELECT  userId FROM Sessions WHERE  Id=?", cookieVal).Scan(&UserId)
	if err != nil {
	return  "",err
	}
	
	statement.Exec(id, PostTitle, PostBody, UserId, created)

	statement, err= database.Prepare("INSERT INTO PostCat VALUES (?,?)")
	if err != nil {
		return  "",err
		}
      
	for _, v := range categories {
		row := database.QueryRow("SELECT Id FROM Categories WHERE Name = ?", v)

		var CatId uuid.UUID
		err := row.Scan(&CatId)
		if err != nil {
			fmt.Println("error of adding linking post with cats")
			log.Fatal(err)
		

		}
		statement.Exec(id, CatId)
	}

	return id.String(), nil
}

// INSERT INTO Posts (Id, Title, Body, UserId, Created) VALUES
//     ( 'testing', 'test', 'testy',  (SELECT Id from Users WHERE Username='Citric'),"14/02/2023" );
