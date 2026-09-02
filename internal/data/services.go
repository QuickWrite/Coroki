package data

import (
	"context"
	"time"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, encodedHash string) error
	NeedsRehash(encodedHash string) bool
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, email string, password string) (*User, error)
}

type Session struct {
	ID        string
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionService interface {
	Create(
		ctx context.Context,
		userID int64,
	) (*Session, error)

	GetUser(
		ctx context.Context,
		sessionID string,
	) (*User, error)

	Revoke(
		ctx context.Context,
		sessionID string,
	) error
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// The main service that manages the users directly
type UserService interface {
	GetUsers(ctx context.Context) ([]User, error)

	CreateUser(ctx context.Context, name string, email string, password string) (*User, error)
}
