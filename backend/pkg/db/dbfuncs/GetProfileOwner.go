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
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func GetProfileOwner(userId string) (User_GetProfileOwner, error) {
	profile := User_GetProfileOwner{}
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", userId).Scan(&profile.Id, &profile.Nickname, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Password, &profile.Avatar, &profile.AboutMe, &profile.PrivacySetting, &profile.DOB, &profile.CreatedAt)
	if err != nil {
		return profile, err
	}
	return profile, nil
}
