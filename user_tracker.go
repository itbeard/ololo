package main

import (
	"time"
)

const banDuration = 5 * time.Minute

type UserTracker interface {
	AddUser(userID int64)
	IsNewUser(userID int64) bool
	RemoveUser(userID int64)
}

type SimpleUserTracker struct {
	users map[int64]time.Time
}

func NewSimpleUserTracker() *SimpleUserTracker {
	return &SimpleUserTracker{users: make(map[int64]time.Time)}
}

func (t *SimpleUserTracker) AddUser(userID int64) {
	t.users[userID] = time.Now()
}

func (t *SimpleUserTracker) IsNewUser(userID int64) bool {
	joinTime, exists := t.users[userID]
	if !exists {
		return false
	}
	if time.Since(joinTime) > banDuration {
		delete(t.users, userID)
		return false
	}
	return true
}

func (t *SimpleUserTracker) RemoveUser(userID int64) {
	delete(t.users, userID)
}
