package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"time"

	"github.com/google/uuid"
)

var CharacterLimit int = 63206

type Image struct {
	Data string `json:"data"`
}

type Post struct {
	Id              string    `json:"Id"`
	UserId          string    `json:"UserId"`
	Title           string    `json:"Title"`
	Body            string    `json:"Body"`
	CreatedAt       time.Time `json:"CreatedAt"`
	Comments        []Comment `json:"Comments"`
	Likes           int       `json:"Likes"`
	Dislikes        int       `json:"Dislikes"`
	PrivacyLevel    string    `json:"PrivacyLevel"`
	CreatorId       string    `json:"CreatorId"`
	Image           string    `json:"Image,omitempty"`
	GroupId         string    `json:"GroupId"`
	ChosenFollowers []string  `json:"ChosenFollowers"`
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

type BasicGroupInfo struct {
	Id          string    `json:"Id"`
	CreatorId   string    `json:"CreatorId"`
	Name        string    `json:"Title"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type GroupCard struct {
	BasicInfo BasicGroupInfo `json:"BasicInfo"`
	Status    string         `json:"Status"`
}

type GroupEvent struct {
	Id          string    `json:"Id"`
	GroupId     string    `json:"GroupId"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	CreatorId   string    `json:"CreatorId"`
	Time        time.Time `json:"Time"`
	Going       int       `json:"Going"`
	NotGoing    int       `json:"NotGoing"`
}

type GroupEventCard struct {
	Event GroupEvent `json:"event"`
	Going bool       `json:"Going"`
}

type DetailedGroupInfo struct {
	BasicInfo        BasicGroupInfo   `json:"BasicInfo"`
	InvitedMembers   []BasicUserInfo  `json:"InvitedMembers"`
	RequestedMembers []BasicUserInfo  `json:"RequestedMembers"`
	Members          []BasicUserInfo  `json:"Members"`
	Posts            []Post           `json:"Posts"`
	EventCards       []GroupEventCard `json:"Events"`
	Messages         []GroupMessage   `json:"Messages"`
	Status           string           `json:"Status"`

	Invite []BasicUserInfo `json:"Invite"`
}

type GroupDash struct {
	GroupCards []GroupCard `json:"GroupCards"`
}

type MessageData struct {
	CurrUser  string `json:"currUser"`
	OtherUser string `json:"otherUser"`
	Type      string `json:"type"`
}

type PostQuery struct {
	PostID string `json:"post_id"`
	UserId string `json:"user_id"`
}

type Reaction struct {
	Postid string `json:"postId"`
	Query  string `json:"query"`
	UserId string `json:"id"`
}

type Profile struct {
	Owner            dbfuncs.User
	Posts            []dbfuncs.Post
	Followers        []BasicUserInfo
	Following        []BasicUserInfo
	PendingFollowers []BasicUserInfo
	IsFollowed       bool
	IsPending        bool
}

type PrivcySetting struct {
	Privacy string `json:"setting"`
	UserId  string `json:"id"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchQuery struct {
	Search string `json:"search"`
}

type SearchFollowQuery struct {
	Search string `json:"search"`
	Id     string `json:"id"`
}
