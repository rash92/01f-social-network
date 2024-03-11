package handlefuncs

import (
	"encoding/json"
	"server/pkg/db/dbfuncs"

	"net/http"
)

func HandleNewUser(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodPost {

		err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
		if err != nil {
			http.Error(w, "Unable to parse form its too big  image should be less than 20Mb", http.StatusBadRequest)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		NickName := r.FormValue("nickname")
		dob := r.FormValue("DOB")
		aboutMe := r.FormValue("aboutMe")

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
		// fmt.Println(fileName)

		// insert user into SQLite database

		_, isexistNickName, err := dbfuncs.CheckValueInDB(w, r, NickName, "Nickname")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if isexistNickName {

			http.Error(w, `{"error": "NickName is already take please choose another one !"}`, http.StatusBadRequest)
			return
		}

		_, isexistemail, err := dbfuncs.CheckValueInDB(w, r, email, "Email")

		if err != nil {

			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if isexistemail {

			http.Error(w, `{"error": "Email is already take please choose another one !"}`, http.StatusBadRequest)
			return
		}

		err = dbfuncs.AddUser(NickName, firstName, lastName, email, fileName, aboutMe, "private", dob, HashPassord(password))

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

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
