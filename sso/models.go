package sso

import "time"

type User struct {
	ID        *int      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           string
	UserEmail    string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type RefreshSession struct {
	ID           *int
	UserID       *int
	RefreshToken string
	ExpiresAt    time.Time
}
