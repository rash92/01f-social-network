package main

import (
	"backend/pkg/db/dbfuncs"
	"backend/pkg/handlefuncs"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func wrapperHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlefuncs.Cors(&w, r)

		cookie, err := r.Cookie("user_token")
		if err != nil {
			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
			return
		}

		cookieValue, err := dbfuncs.ValidateCookie(cookie.Value)

		if !(err == nil && cookieValue) {
			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func main() {

	db, err := sql.Open("sqlite3", "file:./pkg/db/sqlite/sqlite.db?_foreign_keys=on")
	if err != nil {
		log.Fatal("Invalid DB config, unable to open database:", err)
	}
	defer db.Close()

	dbfuncs.Configure(db)

	// sqlite.Migrate() //this line resets the database when the server runs, commented out to persist
	http.HandleFunc("/ws", wrapperHandler(handlefuncs.HandleConnection))
	http.HandleFunc("/login", handlefuncs.HandleLogin)
	http.HandleFunc("/checksession", handlefuncs.HandleValidateCookie)
	http.HandleFunc("/add-comment", wrapperHandler(handlefuncs.HandleAddComment))
	http.HandleFunc("/group", wrapperHandler(handlefuncs.HandleGroup))
	http.HandleFunc("/react-Post-like-dislike", wrapperHandler(handlefuncs.HandlePostLikeDislike))
	http.HandleFunc("/get-group-messages", wrapperHandler(handlefuncs.HandleGetGroupMessages))
	http.HandleFunc("/newUser", handlefuncs.HandleNewUser)
	http.HandleFunc("/get-post", wrapperHandler(handlefuncs.HandleGetPost))
	http.HandleFunc("/logout", wrapperHandler(handlefuncs.HandleLogout))
	http.HandleFunc("/dashboard", wrapperHandler(handlefuncs.HandleDashboard))
	http.HandleFunc("/profile", wrapperHandler(handlefuncs.HandleGetProfile))
	http.HandleFunc("/get-users", wrapperHandler(handlefuncs.HandleGetUsers))
	http.HandleFunc("/get-messages", wrapperHandler(handlefuncs.MessagesHandler))
	http.HandleFunc("/search", wrapperHandler(handlefuncs.HandleSearchUser))
	http.HandleFunc("/toggle-privacy", wrapperHandler(handlefuncs.HandleToggleProfilePrivacy))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./pkg/db/images"))))
	http.Handle("/search-Follower", wrapperHandler(handlefuncs.HandleSearchFollower))

	fmt.Println("Starting server on http://localhost:8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
