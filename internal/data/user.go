package data

import "context"

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore interface {
	GetUsers(ctx context.Context) ([]User, error)
}
