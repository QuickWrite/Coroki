package data
import (
	"context"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, encodedHash string) error
	NeedsRehash(encodedHash string) bool
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// The main service that manages the users directly
type UserService interface {
	GetUsers(ctx context.Context) ([]User, error)
}
