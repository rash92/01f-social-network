package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
	handlefuncs "server/pkg/handlefuncs"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
}

func wrapperHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlefuncs.Cors(&w, r)
		cookie, err := r.Cookie("user_token")

		if err != nil || !dbfuncs.ValidateCookie(cookie.Value) {
			// fmt.Println("cookie value wrapper function", cookie.Value)
			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func main() {
	defer db.Close()
	handlefuncs.SetDatabase(db)
	dbfuncs.SetDatabase(db)
	// DeleteUserByUsername("bilal")
	// sqlite.Magarate()
	http.HandleFunc("/ws", wrapperHandler(handlefuncs.HandleConnection))

	http.HandleFunc("/newUser", handlefuncs.HandleNewUser)
	http.HandleFunc("/check-nickname", handlefuncs.HanndleUserNameIsDbOrEmail)
	http.HandleFunc("/check-email", handlefuncs.HanndleUserNameIsDbOrEmail)
	http.HandleFunc("/login", handlefuncs.HandleLogin)
	http.HandleFunc("/checksession", handlefuncs.HandleValidateCookie)
	http.HandleFunc("/add-post", handlefuncs.HandleAddPost)
	http.HandleFunc("/get-catogries", handlefuncs.HandleCatogries)
	http.HandleFunc("/get-posts", wrapperHandler(handlefuncs.HandleGetPosts))
	http.HandleFunc("/add-Comment", handlefuncs.HandleAddComment)
	http.HandleFunc("/logout", handlefuncs.HandleLogout)
	http.HandleFunc("/react-Post-like-dislike", handlefuncs.HandlePostLikeDislike)
	http.HandleFunc("/react-comment-like-dislike", handlefuncs.HandleCommenttLikeDislike)
	http.HandleFunc("/removePost", handlefuncs.HandleRemovePost)
	http.HandleFunc("/profile", wrapperHandler(handlefuncs.HandleGetProfile))
	http.HandleFunc("/get-users", wrapperHandler(handlefuncs.HandleGetUsers))
	http.HandleFunc("/get-messages", handlefuncs.MessagesHandler)
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./pkg/db/images"))))

	fmt.Println("Starting server on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
