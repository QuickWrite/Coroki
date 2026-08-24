package db

import (
	"context"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
)

type DbUserService struct {
	query *db.Queries
}

func NewDbUserService(query *db.Queries) DbUserService {
	return DbUserService{
		query: query,
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
