package handlefuncs

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	dbfuncs "server/pkg/db/dbfuncs"
// )

// func HandleFilter(w http.ResponseWriter, r *http.Request) {

// 	Cors(&w, r)

// 	if r.Method == http.MethodPost {
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

// 		response := map[string]interface{}{
// 			"success": true,
// 		}
// 		json.NewEncoder(w).Encode(response)

// 		// w.WriteHeader(http.StatusOK)
// 	} else {
// 		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}

// }
