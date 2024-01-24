package dbfuncs

import (
	"database/sql"
	"time"
)

// global database variable, so we only have to open it once and can access it etc.
// possibly we don't want it globally and open and close it as needed
var db *sql.DB

// structs based on database for entering and retrieving info
type User struct {
	Id             string
	NickName       string
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Profile        string
	AboutMe        string
	PrivacySetting string
	DateOfBirth    string
	CreatedAt      time.Time
}

type PrivateMessage struct {
	Id          string
	SenderId    string
	RecipientId string
	Message     string
	CreatedAt   time.Time
}

type Post struct {
	Id           string
	Title        string
	Body         string
	CreatorId    string
	GroupId      string
	CreatedAt    time.Time
	Image        string
	PrivacyLevel string
}

type Comment struct {
	Id        string
	Body      string
	CreatorId string
	PostId    string
	CreatedAt time.Time
	Image     string
}

type Follow struct {
	FollowerId  string
	FollowingId string
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
	GroupId   string
	Message   string
	CreatedAt time.Time
}

// misspelled in database - fix reciever to receiver
type Notification struct {
	Id         string
	Body       string
	Type       string
	CreatedAt  time.Time
	ReceiverId string
	SenderId   string
	Seen       bool
}

//above are based on database fields, to deal with database. May need separate structs for front end e.g. version of user without password field

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
