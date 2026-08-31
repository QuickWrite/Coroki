package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
)

type DbAuthenticationService struct {
	query    *db.Queries
	password *data.PasswordService
}

func NewDbAuthenticationService(query *db.Queries, password *data.PasswordService) DbAuthenticationService {
	return DbAuthenticationService{
		query:    query,
		password: password,
	}
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s DbAuthenticationService) Authenticate(ctx context.Context, email string, password string) (*data.User, error) {
	user, err := s.query.GetUserWithPasswordHash(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("get user credentials: %w", err)
	}

	if err := (*s.password).Compare(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &data.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
