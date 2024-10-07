package model

type Session struct {
	Token     string `json:"session"`
	ExpiresAt string `json:"expires_at"`
}
