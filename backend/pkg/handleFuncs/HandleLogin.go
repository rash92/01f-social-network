package handlefuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodPost {

		var entredData LoginData

		err := json.NewDecoder(r.Body).Decode(&entredData)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		id, err := dbfuncs.IsLoginValid(entredData.Email, entredData.Password)

		if err != nil {
			fmt.Println(err.Error(), "error after getting data")
			http.Error(w, `{"error": "your email/nickname or password is incorrect"}`, http.StatusBadRequest)
			return

		}

		// var (
		// 	id             string
		// 	username       string
		// 	storedPassword string
		// 	imgUrl         string
		// )
		// // var storedPassword string
		// err := dbfuncs.QueryRow("SELECT Id, Password, Nickname,  Profile  FROM Users WHERE Email=?", entredData.Email).Scan(&id, &storedPassword, &username, &imgUrl)
		// if err != nil {
		// 	fmt.Println(err.Error(), "error after getting data")
		// 	http.Error(w, `{"error": "your email/nickname or password is incorrect"}`, http.StatusBadRequest)
		// 	return
		// }

		// if isPasswordValid([]byte(storedPassword), []byte(entredData.Password)) != nil {
		// 	// fmt.Println(isPasswordValid([]byte(storedPassword), []byte(entredData.Password)))
		// 	http.Error(w, `{"error": "your email or passord is incorrect"}`, http.StatusBadRequest)
		// 	return

		// }

		user, err := dbfuncs.GetUserById(id)

		if err != nil {
			log.Println(err.Error(), "error after getting user")
			http.Error(w, `{"error": "something went wrong please try again"}`, http.StatusInternalServerError)
			return
		}

		session, err := dbfuncs.AddSession(id)
		if err != nil {
			log.Println(err.Error(), "error after adding session")
			http.Error(w, `{"error": "something went wrong please try again"}`, http.StatusInternalServerError)
			return

		}

		http.SetCookie(w, &http.Cookie{
			Name:     "user_token",
			Value:    session.Id,
			Expires:  session.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		user.Password = []byte{}
		user.Email = ""

		response := map[string]interface{}{
			"success": true,
			"user":    user,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

}
