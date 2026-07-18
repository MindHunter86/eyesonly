package auth

import (
	"context"
	"time"

	"github.com/MindHunter86/eyesonly/internal/webapp/platform/token"
	"github.com/MindHunter86/eyesonly/internal/webapp/users"
)

type Service struct {
	st    *Store
	users *users.Store
	jwtm  *token.JWTManager
}

type TokenPair struct {
	AccessToken     string           `json:"access_token"`
	RefreshToken    string           `json:"refresh_token,omitempty"`
	RefreshExpireAt time.Time        `json:"refresh_expires_at"`
	User            users.PublicUser `json:"user"`
}

// !! todo - get fiber.HEADERS
func (m *Service) Register(c context.Context, email, pass string) (_ *TokenPair, e error) {
	return
}

func (m *Service) Login(c context.Context, email, pass string) (_ *TokenPair, e error) {
	return
}

func (m *Service) Refresh(c context.Context, payload string) (_ *TokenPair, e error) {
	return
}

func (m *Service) Logout(c context.Context, payload string) error {
	return nil
}

func (m *Service) LogoutAll(c context.Context, uid string) error {
	return nil
}

// internal
// !! TODO - NEED CONFIG HERE
func (m *Service) issueTokenPair(c context.Context, u users.User) (_ *TokenPair, e error) {
	var opa string
	if opa, e = token.NewOpaqueToken(); e != nil {
		return nil, e
	}

	now := time.Now().UTC()
	rt := &RefreshToken{
		ID:        "",
		UID:       u.ID,
		TokenHash: token.HashOpaqueToken(opa),
		ExpiresAt: now.Add(0 * time.Second),
		CreatedAt: now,
	}

	if e = m.st.CreateRefreshToken(c, rt); e != nil {
		return
	}

	var at string
	if at, e = m.jwtm.Create(u.ID, u.Role, rt.ID); e != nil {
		return
	}

	return &TokenPair{
		AccessToken:     at,
		RefreshToken:    opa,
		RefreshExpireAt: rt.ExpiresAt,
		User:            *users.Public(&u),
	}, nil
}
