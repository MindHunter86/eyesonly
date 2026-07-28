package sessions

import "time"

type Session struct {
	ID   string
	Hash string
}

type Token struct {
	Token     string    `json:"token"`
	Session   string    `json:"session"`
	ExpiresAt time.Time `json:"expires_at"`
}
