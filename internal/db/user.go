package db

import (
	"context"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
)

type DbUserService struct {
	query    *db.Queries
	password *data.PasswordService
}

func NewDbUserService(query *db.Queries, password *data.PasswordService) DbUserService {
	return DbUserService{
		query:    query,
		password: password,
	}
}

func (s DbUserService) GetUsers(ctx context.Context) ([]data.User, error) {
	dbUsers, err := s.query.GetUsers(ctx)

	if err != nil {
		return nil, err
	}

	users := make([]data.User, len(dbUsers))
	for i, v := range dbUsers {
		users[i] = data.User{
			ID:    v.ID,
			Name:  v.Name,
			Email: v.Email,
		}
	}

	return users, nil
}

func (s DbUserService) CreateUser(ctx context.Context, name string, email string, password string) (*data.User, error) {
	hash, err := (*s.password).Hash(password)

	if err != nil {
		return nil, err
	}

	user, err := s.query.CreateUserWithPasswordHash(ctx, db.CreateUserWithPasswordHashParams{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	})

	if err != nil {
		return nil, err
	}

	return &data.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
