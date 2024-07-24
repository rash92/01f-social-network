package dbfuncs

import (
	"database/sql"

	"log"
	"sync"
	"time"
)

// global database variable, so we only have to open it once and can access it etc.
// possibly we don't want it globally and open and close it as needed
var db *sql.DB

var dbLock sync.RWMutex

// opens database at beginning, should close automatically on server quit
func init() {
	var err error
	db, err = sql.Open("sqlite3", "file:./pkg/db/sqlite/sqlite.db?_foreign_keys=on")
	if err != nil {
		log.Fatal("Invalid DB config, unable to open database:", err)
	}
}

// structs based on database for entering and retrieving info
type User struct {
	Id             string
	Nickname       string
	FirstName      string
	LastName       string
	Email          string
	Password       []byte
	Avatar         string
	AboutMe        string
	PrivacySetting string
	DOB            string
	CreatedAt      time.Time
}

// avatar and nickname isn't used in frontend so could get rid of it.
type PrivateMessage struct {
	Id          string
	SenderId    string
	RecipientId string
	Message     string
	CreatedAt   time.Time
	Nickname    string
	Avatar      string
}

type Post struct {
	Id              string
	Title           string
	Body            string
	CreatorId       string
	GroupId         string
	CreatedAt       time.Time
	Image           string
	PrivacyLevel    string
	Likes           int
	Dislikes        int
	CreatorNickname string
	UserLikeDislike int
	Comments        []Comment
	Ncomment        int
}

type Comment struct {
	Id              string
	Body            string
	CreatorId       string
	PostId          string
	CreatedAt       time.Time
	Image           string
	Likes           int
	Dislikes        int
	CreatorNickname string
}

type Follow struct {
	FollowerId  string
	FollowingId string
	Status      string
}

type Group struct {
	Id          string
	Title       string
	Description string
	CreatorId   string
	CreatedAt   time.Time
}

// double check what status means
type GroupMember struct {
	GroupId string
	UserId  string
	Status  string
}

type GroupEvent struct {
	Id          string
	GroupId     string
	Title       string
	Description string
	CreatorId   string
	Time        time.Time
}

// currently different to database, will change database to replace choice with groupId
type GroupEventParticipant struct {
	EventId string
	UserId  string
	GroupId string
}

type Session struct {
	Id      string
	Expires time.Time
	UserId  string
}

type PostChosenFollower struct {
	PostId     string
	FollowerId string
}

type PostLike struct {
	UserId   string
	PostId   string
	Liked    bool
	Disliked bool
}

type CommentLike struct {
	UserId    string
	CommentId string
	Liked     bool
	Disliked  bool
}

type GroupMessage struct {
	Id        string
	SenderId  string
	Nickname string
	Avatar string
	GroupId   string
	Message   string
	CreatedAt time.Time
}

// misspelled in database - fix reciever to receiver
type Notification struct {
	Id           string
	Body         string
	Type         string
	CreatedAt    time.Time
	ReceiverId   string
	SenderId     string
	Seen         bool
	SenderAvatar string
}
