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
type User struct {
	Email      string    `json:"email"`
	NickName   string    `json:"nickname"`
	FirstName  string    ` json:"firstName"`
	LastName   string    `json:"lastName"`
	Age        string    `json:"age"`
	Gender     string    `json:"gender"`
	Password   string    `json:"password"`
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Aboutme    string  `json:"aboutme"`
	Avatar     *Image    `json:"avatar,omitempty"`
}

type Categories struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}

type Post struct {
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
type Comment struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Username  string    `json:"username"`
}
type Session struct {
	Id       uuid.UUID
	Username string
	Expires  time.Time
	UserID   string
}

type PostFontEnd struct {
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

type Message struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Message     string `json:"message"`
	Created     string `json:"created"`
	Type        string `json:"type"`
	Typing      bool   `json:"typing"`
}
