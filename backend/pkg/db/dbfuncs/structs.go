package dbfuncs

import "sync"

type Event struct {
	EventId     string `json:"EventId"`
	GroupId     string `json:"GroupId"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	CreatorId   string `json:"CreatorId"`
}

var lock sync.RWMutex
