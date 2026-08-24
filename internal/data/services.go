package data
import (
	"context"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// The main service that manages the users directly
type UserService interface {
	GetUsers(ctx context.Context) ([]User, error)
}
