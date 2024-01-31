package handlefuncs

import (
	dbfuncs "server/pkg/db/dbfuncs"
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
	Created_at time.Time         `json:"created_at"`
	Comments   []dbfuncs.Comment `json:"comments"`
	Likes      int               `json:"likes"`
	Dislikes   int               `json:"dislikes"`
	Username   string            `json:"username"`
	Userlikes  []string          `json:"userlikes"`
}
