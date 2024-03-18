package app

import (
	dbfuncs "backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == "POST" {

		//should NOT be able to find a cookie at this point, either err == nil and found a cookie or some other error happend, either way something went wrong.
		//happy path is for err == ErrNoCookie
		cookie, err := r.Cookie("userToken")
		if err == nil {
			err = dbfuncs.DeleteSession(cookie.Value)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		} else if err.Error() != "ErrNoCookie" {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		var loginData LoginData

		err = json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		userId, err := dbfuncs.IsLoginValid(loginData.Email, loginData.Password)
		if err != nil {
			http.Error(w, `{"error": "email or password invalid"}`, http.StatusBadRequest)
			return
		}

		session, err := dbfuncs.AddSession(userId)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "userToken",
			Value:    session.Id,
			Expires:  session.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		response := struct {
			UserId string `json:"userId"`
		}{
			UserId: userId,
		}
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

}
