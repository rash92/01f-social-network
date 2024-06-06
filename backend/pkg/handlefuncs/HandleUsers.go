package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		users, err := dbfuncs.GetUsers()
		fmt.Println(users, "users")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(users)

	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func HandleNewUser(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "unable to parse form: image more than 20Mb", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password, err := dbfuncs.HashPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	nickname := r.FormValue("nickname")
	dob := r.FormValue("DOB")
	aboutMe := r.FormValue("aboutMe")

	if nickname == "" {
		nickname = firstName + lastName
	}

	file, header, err := r.FormFile("image")

	if err != nil && err != http.ErrMissingFile {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var fileName string
	if file != nil {
		fileName, err = SaveImage(file, header)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
	}

	DBUser := dbfuncs.User{
		Email:          email,
		Password:       password,
		FirstName:      firstName,
		LastName:       lastName,
		Nickname:       nickname,
		DOB:            dob,
		AboutMe:        aboutMe,
		Avatar:         fileName,
		PrivacySetting: "private",
	}

	isexistemail, err := dbfuncs.CheckEmailInDB(email)

	if err != nil {

		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if isexistemail {

		http.Error(w, `{"error": "Email is already take please choose another one !"}`, http.StatusBadRequest)
		return
	}

	err = dbfuncs.AddUser(&DBUser)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	// return success response
	responseData := map[string]string{"message": "User created successfully"}
	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	// Set the "Content-Type" header to "application/json"
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response to the HTTP response
	w.Write(response)
}

func HandleSearchFollower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var searchQuery SearchFollowQuery

	err := json.NewDecoder(r.Body).Decode(&searchQuery)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := dbfuncs.SearchFollowers(searchQuery.Search, searchQuery.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error": "Error getting user from database"}`, http.StatusInternalServerError)
		return
	}

	// Send the user back to the client
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, `{"error": "Error encoding JSON"}`, http.StatusInternalServerError)
		return
	}
}

func HandleSearchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var searchQuery SearchQuery

	err := json.NewDecoder(r.Body).Decode(&searchQuery)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := dbfuncs.SearchUsers(searchQuery.Search)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error": "Error getting user from database"}`, http.StatusInternalServerError)
		return
	}

	// Send the user back to the client
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, `{"error": "Error encoding JSON"}`, http.StatusInternalServerError)
		return
	}
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

		fmt.Println("session might expire, but who really knows ...", session.Expires)

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

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,

		SameSite: http.SameSiteLaxMode,
	})
	response := map[string]interface{}{
		"success": true,
	}
	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)

}

func HandleValidateCookie(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("user_token")

		if cookie == nil || err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return

		}
		id, err := dbfuncs.GetUserIdFromCookie(cookie.Value)

		if err != nil {

			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		user, err := dbfuncs.GetUserById(id)

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"user":    user,
		}

		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)

	}
}
