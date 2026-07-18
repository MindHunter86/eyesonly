package auth

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

type RefreshToken struct {
	ID        string
	UID       string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	CreatedAt time.Time
}

func (m *Store) CreateRefreshToken(c context.Context, t *RefreshToken) (e error) {
	return
}

func (m *Store) FundRefreshTokenByHash(c context.Context, hash string) (e error) {
	return
}

func (m *Store) RevokeRefreshToken(c context.Context, hash string) (e error) {
	return

}

func (m *Store) RevokeAllForUser(c context.Context, uid string) (e error) {
	return
}
