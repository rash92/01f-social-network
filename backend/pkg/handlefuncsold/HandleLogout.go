package handlefuncs

import (
	"encoding/json"
	"net/http"
	"time"
)

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
