package app

import (
	"server/pkg/db/dbfuncs"
	"time"
)

//structs to be sent to front end or for intermediate steps

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
}

type RegistrationData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
	Profile     string `json:"profile"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"aboutMe"`
}

type PostFrontEnd struct {
}

//above are based on database fields, to deal with database. May need separate structs for front end e.g. version of user without password field

type PostFontEndOld struct {
	Id         string            `json:"id"`
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Categories []string          `json:"categories"`
	CreatedAt  time.Time         `json:"createdAt"`
	Comments   []dbfuncs.Comment `json:"comments"`
	Likes      int               `json:"likes"`
	Dislikes   int               `json:"dislikes"`
	Username   string            `json:"username"`
	Userlikes  []string          `json:"userlikes"`
}

// this is what the client sends to the server,
// Body will be reunmarshalled based on type into PrivateMessage, GroupMessage, or Notification etc.
// as well as ws messages unrelated to database operations
type WsMessage struct {
	Type      string    `json:"type"`
	Body      []byte    `json:"message"`
	TimeStamp time.Time `json:"time"`
}

// Consider whether we want to be sending id numbers to the client.
// We can continue for now as if we are, and modify this later if we
// decide not to. Be sure to investigate the security implications
// of sending id numbers to the client. If we decide not to send id
// numbers, we need to consider the actual threat model, and avoid
// naively attempting to fix things that don't actually address the
// real threat.  Matt is sending them and happy that it's not a problem.

// these should match the database fields, check types in database and
// make them match, also move them to db package eventually
// Do they need to match the fields in the database? For example,
// the frontend won't know the id of a message until it's been added
// to the database. At the moment, signals are being sent to the
// frontend in the form of a map from string to interface, with
// whatever key-value pairs are needed for that type of signal.
type PrivateMessage struct {
	Id          string    `json:"Id"`
	SenderId    string    `json:"SenderId"`
	RecipientId string    `json:"RecipientId"`
	Message     string    `json:"Message"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type GroupMessage struct {
	Id        string    `json:"Id"`
	SenderId  string    `json:"SenderId"`
	GroupId   string    `json:"GroupId"`
	Message   string    `json:"Message"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type RequestToFollow struct {
	SenderId    string `json:"SenderId"`
	RecipientId string `json:"RecipientId"`
}

// We can just use the Event type for this.
type ToggleAttendEvent struct {
	EventId string `json:"EventId"`
	UserId  string `json:"UserId"`
	GroupId string `json:"GroupId"`
	Type    string `json:"Type"`
}

type Notification struct {
	Id         string    `json:"Id"`
	RecieverId string    `json:"RecieverId"`
	SenderId   string    `json:"SenderId"`
	Body       string    `json:"Body"`
	Type       string    `json:"Type"`
	Seen       bool      `json:"Seen"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type Post struct {
	//fill in later
}

type Comment struct {
	//fill in later
}

//these may not involve database calls but can still be sent through websockets

// this can be resused for sending WsMessages to other users about a user who has
// e.g. registered, logged in, logged out, or changed their status
type BasicUserInfo struct {
	UserId         string `json:"UserId"`
	PrivacySetting string `json:"PrivacySetting"` //maybe?
	//fill in later if more needed
}
