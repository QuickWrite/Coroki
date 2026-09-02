package db

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/QuickWrite/Coroki/internal/data"
	db "github.com/QuickWrite/Coroki/internal/db/sqlc/gen"
)

type DbSessionService struct {
	query *db.Queries
}

func NewDbSessionService(query *db.Queries) DbSessionService {
	return DbSessionService{
		query: query,
	}
}

func (s DbSessionService) Create(ctx context.Context, userID int64) (*data.Session, error) {
	sessionID, err := newID()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(sessionID))
	expires := time.Now().Add(6 * time.Hour) // TODO: Make this configurable

	session, err := s.query.CreateSession(ctx, db.CreateSessionParams{
		TokenHash: hash[:],
		UserID:    userID,
		ExpiresAt: expires,
	})

	if err != nil {
		return nil, err
	}

	return &data.Session{
		ID:        sessionID,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s DbSessionService) GetUser(ctx context.Context, sessionID string) (*data.User, error) {
	hash := sha256.Sum256([]byte(sessionID))
	user, err := s.query.GetSessionUser(ctx, hash[:])

	if err != nil {
		return nil, err
	}

	return &data.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s DbSessionService) Revoke(ctx context.Context, sessionID string) error {
	hash := sha256.Sum256([]byte(sessionID))
	err := s.query.RevokeSession(ctx, hash[:])

	return err
}

func newID() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
