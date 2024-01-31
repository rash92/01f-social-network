package handlefuncs

import (
	"encoding/json"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == "POST" {

		var loginData LoginData

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		var storedPassword string

		err = dbfuncs.db.QueryRow("SELECT Password FROM Users WHERE Email = ?", loginData.Email).Scan(&storedPassword)

	}

}

func HandleLoginOld(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == "POST" {

		var loginData LoginData

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		var id string
		var username string
		var storedPassword string
		var imgUrl string

		// var storedPassword string
		err := database.QueryRow("SELECT Id, Password, Nickname,  Profile  FROM Users WHERE Nickname=? OR mail=?", entredData.Nickname, entredData.Email).Scan(&id, &storedPassword, &username, &imgUrl)
		if err != nil {
			http.Error(w, `{"error": "your email/nickname or password is incorrect"}`, http.StatusBadRequest)
			return
		}

		if isPasswordValid(storedPassword, entredData.Password) != nil {
			// fmt.Println(isPasswordValid([]byte(storedPassword), []byte(entredData.Password)))
			http.Error(w, `{"error": "your email or passord is incorrect"}`, http.StatusBadRequest)
			return

		}

		sessionId, _ := uuid.NewRandom()

		session := Session{
			Id:       sessionId,
			Username: username,
			Expires:  time.Now().Add(24 * time.Hour),
			UserID:   id,
		}
		//  detlete old session
		dbfuncs.DeleteSessionColumn("userId", id)
		// add new session
		dbfuncs.AddSession(session.Id, session.Username, session.UserID, session.Expires)

		http.SetCookie(w, &http.Cookie{
			Name:     "user_token",
			Value:    sessionId.String(),
			Expires:  session.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		response := map[string]interface{}{
			"success":    true,
			"username":   session.Username,
			"profileImg": imgUrl,
			"id":         session.UserID,
		}
		json.NewEncoder(w).Encode(response)

		// w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

}
