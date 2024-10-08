package model

import "time"

type Session struct {
	Token     string    `json:"session"`
	ExpiresAt time.Time `json:"expires_at"`
}
