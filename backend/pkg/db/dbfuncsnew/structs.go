package dbfuncs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var database *sql.DB

func SetDatabase(db *sql.DB) {
	database = db
}

type User struct {
	Id              uuid.UUID
	FirstName       string
	LastName        string
	Age             int
	Gender          string
	Img             string
	Created         time.Time
	LastMessageTime time.Time
}

type PrivateMessage struct {
	Id          string
	SenderId    string
	RecipientId string
	CreatedAt   time.Time
}

type Posts struct {
	Id        string
	Title     string
	Body      string
	CreatorId string
}

type Comment struct {
}

type ImageOld struct {
	Data string `json:"data"`
}
type UserOld struct {
	Email      string    `json:"email"`
	NickName   string    `json:"nickname"`
	FirstName  string    ` json:"firstName"`
	LastName   string    `json:"lastName"`
	Age        string    `json:"age"`
	Gender     string    `json:"gender"`
	Password   string    `json:"password"`
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Aboutme    string    `json:"aboutme"`
	Avatar     *Image    `json:"avatar,omitempty"`
}

type CategoriesOld struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}

type PostOld struct {
	Id         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"userid"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Categories []string  `json:"categories"`
	Created_at time.Time `json:"created_at"`
	Comments   []Comment `json:"comments"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
}
type CommentOld struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Username  string    `json:"username"`
}
type SessionOld struct {
	Id       uuid.UUID
	Username string
	Expires  time.Time
	UserID   string
}

type PostFontEndOld struct {
	Id         string    `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Categories []string  `json:"categories"`
	Created_at time.Time `json:"created_at"`
	Comments   []Comment `json:"comments"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Username   string    `json:"username"`
	Userlikes  []string  `json:"userlikes"`
}

type MessageOld struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Message     string `json:"message"`
	Created     string `json:"created"`
	Type        string `json:"type"`
	Typing      bool   `json:"typing"`
}
