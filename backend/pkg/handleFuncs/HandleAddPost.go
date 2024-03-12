package handlefuncs

// func HandleAddPost(w http.ResponseWriter, r *http.Request) {
// 	Cors(&w, r)

// 	if r.Method == http.MethodPost {
// 		var newPost Post
// 		errj := json.NewDecoder(r.Body).Decode(&newPost)
// 		if errj != nil {
// 			fmt.Println("bilal", errj.Error(), r.Body)
// 			http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
// 			return
// 		}
// 		cookie, err := r.Cookie("user_token")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
// 			return
// 		}
// 		if !dbfuncs.ValidateCookie(cookie.Value) {

// 			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
// 			return
// 		}

// 		if len(newPost.Body) > CharacterLimit {

// 			http.Error(w, `{"error": "413 Payload Too Large"}`, http.StatusRequestEntityTooLarge)
// 			return
// 		}
// 		if len(newPost.Body) == 0 {

// 			http.Error(w, `{"error": "204 No Content"}`, http.StatusNoContent)
// 			return
// 		}

// 		PostTitle := newPost.Title
// 		PostBody := newPost.Body
// 		Categories := newPost.Categories

// 		id, err := dbfuncs.AddPost(cookie.Value, PostTitle, PostBody, Categories)

// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 		}

// 		data, err := FindPostByID(id)

// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 		}
// 		response := map[string]interface{}{
// 			"success": true,
// 			"post":    data,
// 		}
// 		json.NewEncoder(w).Encode(response)

// 	} else {

// 		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}
// }
