package main

import (
	"backend/pkg/db/dbfuncs"

	handlefuncs "backend/pkg/handlefuncsold"
	"fmt"
	"net/http"
)

// func wrapperHandler(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		handlefuncs.Cors(&w, r)
// 		cookie, err := r.Cookie("user_token")
// 		if err != nil {
// 			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
// 			return
// 		}
// 		isValidCookie, err := dbfuncs.ValidateCookie(cookie.Value)
// 		if err != nil || !isValidCookie {
// 			// fmt.Println("cookie value wrapper function", cookie.Value)
// 			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
// 			return
// 		}
// 		handler(w, r)
// 	}
// }

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
			// fmt.Println("cookie value wrapper function", cookie.Value)
			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func main() {

	// dbfuncs.DeleteUserByUsername("Accepted")
	// sqlite.Migrate()
	http.HandleFunc("/ws", wrapperHandler(handlefuncs.HandleConnection))

	http.HandleFunc("/newUser", handlefuncs.HandleNewUser)
	// http.HandleFunc("/check-nickname", handlefuncs.HanndleUserNameIsDbOrEmail)
	// http.HandleFunc("/check-email", handlefuncs.HanndleUserNameIsDbOrEmail)
	http.HandleFunc("/login", handlefuncs.HandleLogin)
	http.HandleFunc("/checksession", handlefuncs.HandleValidateCookie)
	// http.HandleFunc("/add-post", handlefuncs.HandleAddPost)
	http.HandleFunc("/get-catogries", handlefuncs.HandleCatogries)
	// http.HandleFunc("/get-posts", wrapperHandler(handlefuncs.HandleGetPosts))
	// Commented out because definition change to placeholder for the sake of the
	// web sockets.
	// http.HandleFunc("/add-Comment", handlefuncs.HandleAddComment)
	http.HandleFunc("/logout", handlefuncs.HandleLogout)
	// http.HandleFunc("/react-Post-like-dislike", handlefuncs.HandlePostLikeDislike)
	// http.HandleFunc("/react-comment-like-dislike", handlefuncs.HandleCommenttLikeDislike)
	// http.HandleFunc("/removePost", handlefuncs.HandleRemovePost)
	http.HandleFunc("/dashboard", wrapperHandler(handlefuncs.HandleDashboard))
	http.HandleFunc("/profile", wrapperHandler(handlefuncs.HandleGetProfile))
	http.HandleFunc("/get-users", wrapperHandler(handlefuncs.HandleGetUsers))
	http.HandleFunc("/get-messages", handlefuncs.MessagesHandler)
	http.HandleFunc("/search", wrapperHandler(handlefuncs.HandleSearchUser))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./pkg/db/images"))))

	fmt.Println("Starting server on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
