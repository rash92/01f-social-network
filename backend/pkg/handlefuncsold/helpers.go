package handlefuncs

type BasicUserInfo struct {
	Avatar         string `json:"Avatar"`
	Id             string `json:"Id"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Nickname       string `json:"Nickname"`
	PrivacySetting string `json:"PrivacySetting"`
}

// import "golang.org/x/crypto/bcrypt"

// func HashPassord(password string) []byte {
// 	p, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return p
// }

// func isPasswordValid(h, e []byte) error {
// 	return bcrypt.CompareHashAndPassword(h, e)
// }
