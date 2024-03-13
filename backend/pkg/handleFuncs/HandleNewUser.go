package handlefuncs

import (
	"encoding/json"
	"server/pkg/db/dbfuncs"

	"net/http"
)

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
		fileName, err = dbfuncs.SaveImage(file, header)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
	}

	DBUser := dbfuncs.User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		DOB:       dob,
		AboutMe:   aboutMe,
		Profile:   fileName,
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
