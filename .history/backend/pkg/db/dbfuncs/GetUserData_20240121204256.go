package dbfuncs

import (
	"fmt"
	"net/http"
)

type user struct {
	Email string
	Password string
	FirstName string
	LastName string
}

func GetUserData(id string) user {

	err := database.QueryRow("SELECT Id, Password, Nickname,  Profile,  FROM Users WHERE Email=?", entredData.Email).Scan(&id, &storedPassword, &username, &imgUrl)
	if err != nil {
		fmt.Println(err.Error(), "error after getting data")
		http.Error(w, `{"error": "your email/nickname or password is incorrect"}`, http.StatusBadRequest)
		return
	}

}