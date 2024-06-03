package handlefuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var CharacterLimit int = 63206
var database *sql.DB

func SetDatabase(db *sql.DB) {
	database = db
}

type Image struct {
	Data string `json:"data"`
}

// type User struct {
// 	Email      string    `json:"email"`
// 	NickName   string    `json:"nickname"`
// 	FirstName  string    ` json:"firstName"`
// 	LastName   string    `json:"lastName"`
// 	DOB        string    `json:"age"`
// 	Gender     string    `json:"gender"`
// 	Password   string    `json:"password"`
// 	Id         uuid.UUID `json:"id"`
// 	Created_at time.Time `json:"created_at"`
// 	Aboutme    string    `json:"aboutme"`
// 	Avatar     *Image    `json:"avatar,omitempty"`
// }

type Post struct {
	Id              string    `json:"id"`
	UserId          string    `json:"userid"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	CreatedAt       time.Time `json:"createdAt"`
	Comments        []Comment `json:"comments"`
	Likes           int       `json:"likes"`
	Dislikes        int       `json:"dislikes"`
	PrivacyLevel    string    `json:"privacyLevel"`
	CreatorId       string    `json:"creatorId "`
	Image           string    `json:"avatar,omitempty"`
	GroupId         string    `json:"groupId"`
	ChosenFollowers []string  `json:"chosenFollowers"`
}
type Comment struct {
	Id              string    `json:"Id"`
	Body            string    `json:"Body"`
	CreatorId       string    `json:"CreatorId"`
	PostID          string    `json:"PostId"`
	CreatedAt       time.Time `json:"CreatedAt"`
	Likes           int       `json:"Likes"`
	Dislikes        int       `json:"Dislikes"`
	Image           string    `json:"Image"`
	CreatorNickname string    `json:"CreatorNickname"`
}
type Session struct {
	Id       uuid.UUID
	Username string
	Expires  time.Time
	UserID   string
}

type Message struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Message     string `json:"message"`
	Created     string `json:"created"`
	Type        string `json:"type"`
	Typing      bool   `json:"typing"`
}

type BasicUserInfo struct {
	Avatar         string `json:"Avatar"`
	Id             string `json:"Id"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Nickname       string `json:"Nickname"`
	PrivacySetting string `json:"PrivacySetting"`
}
