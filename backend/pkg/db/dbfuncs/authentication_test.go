package dbfuncs

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	password := "hello"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
		t.Fatalf("hash password function shouldn't error on 'hello'")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Error(err)
		t.Fatalf("brcrypt was unable to compare password to hash")
	}

	password = "this_password_is_too_long_for_bcrypt_because_it_is_more_than_seventy_plus_two_chars"
	_, err = HashPassword(password)
	if err == nil {
		t.Fatal("this was expected to produce an error due to too long password")
	}

}

// func TestIsLoginValid(email, enteredPass string) (string, error) {
// 	var storedPassword string
// 	var userId string
// 	err := db.QueryRow("SELECT Password, Id FROM Users WHERE Email = ?", email).Scan(&storedPassword, &userId)
// 	if err != nil {
// 		return "", err
// 	}
// 	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPass))
// 	if err != nil {
// 		return "", err
// 	}
// 	return userId, nil
// }

// func TestValidateCookie(sessionId string) (bool, error) {
// 	var id string
// 	var expiration time.Time
// 	err := db.QueryRow("SELECT Id, Expires FROM Sessions WHERE Id=?", sessionId).Scan(&id, &expiration)

// 	if err == sql.ErrNoRows {

// 		return false, nil
// 	} else if err != nil {

// 		return false, err
// 	}

// 	return id == sessionId && time.Now().Before(expiration), nil
// }

// func TestCheckEmailInDB(email string) (bool, error) {
// 	found := ""
// 	err := db.QueryRow("SELECT Email FROM Users WHERE Email=?", email).Scan(&found)
// 	if err == sql.ErrNoRows {
// 		return false, nil
// 	} else if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func TestHadAGoodGame(t *testing.T) {
// 	tests := []struct {
// 	   name     string
// 	   stats   Stats
// 	   goodGame bool
// 	   wantErr  string
// 	}{
// 	   {"sad path: invalid stats", Stats{Name: "Sam Cassell",
// 		  Minutes: 34.1,
// 		  Points: -19,
// 		  Assists: 8,
// 		  Turnovers: -4,
// 		  Rebounds: 11,
// 		  }, false, "stat lines cannot be negative",
// 	   },
// 	   {"happy path: good game", Stats{Name: "Dejounte Murray",
// 		  Minutes: 34.1,
// 		  Points: 19,
// 		  Assists: 8,
// 		  Turnovers: 4,
// 		  Rebounds: 11,
// 	   }, true, ""},
// 	}
// 	for _, tt := range tests {
// 	   isAGoodGame, err := hadAGoodGame(tt.stats)
// 	   if tt.wantErr != "" {
// 		  assert.Contains(t, err.Error(), tt.wantErr)
// 	   } else {
// 		  assert.Equal(t, tt.goodGame, isAGoodGame)
// 	   }
// 	}
//  }
