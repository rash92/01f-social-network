package dbfuncs

import (
	"time"
)

// Assume everything here is a placeholder unless it's clear that it's what you want.

// notify peter: already have but different signature
func AddEvent(event *Event) (string, time.Time, error) {
	var err error
	return "", time.Now(), err
}

func GetFollowersByFollowingId(id string) ([]string, error) {
	var followers []string
	return followers, nil
}

// check why we need it
func GetPostPrivacyLevelByCommentId(id string) (string, error) {
	return "", nil
}

func GetCreatedAtByUserId(id string) (time.Time, error) {
	user, err := GetUserById(id)
	if err != nil {
		return time.Time{}, err
	}

	return user.CreatedAt, nil
}

//funcs moved to: users, notifications, groups, comments, posts
//still move funcs in solo files toggleattendevent and getgroupmembers
