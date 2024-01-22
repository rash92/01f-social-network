package dbfuncs

import (
	"time"

	"github.com/google/uuid"
)

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
