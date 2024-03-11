package dbfuncs

import "time"

// Fields as in the database.
type User_GetProfileOwner struct {
	Id             string    `json:"id"`
	Nickname       string    `json:"nickname"`
	FirstName      string    `json:"firstName,omitempty"`
	LastName       string    `json:"lastName,omitempty"`
	Email          string    `json:"email,omitempty"`
	Password       string    `json:"password,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
	AboutMe        string    `json:"aboutme,omitempty"`
	PrivacySetting string    `json:"privacySetting,omitempty"`
	DOB            string    `json:"age,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func GetProfileOwner(userId string) (User_GetProfileOwner, error) {
	profile := User_GetProfileOwner{}
	err := database.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&profile)
	if err != nil {
		return profile, err
	}
	return profile, nil
}
