package sessions

import (
	"strings"
	"time"

	"github.com/MindHunter86/eyesonly/internal/webapp/platform/crypto"
	"github.com/MindHunter86/eyesonly/internal/webapp/platform/idgen"
	"github.com/MindHunter86/eyesonly/internal/webapp/shared/apperr"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	ttl    time.Duration
}

type sessionClaims struct {
	SessionID   string `json:"sid"`
	SessionHash string `json:"sh"`
	jwt.RegisteredClaims
}

func NewService(secret []byte, ttl time.Duration) *Service {
	return &Service{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *Service) Issue(now time.Time) (_ *Token, e error) {
	var sessionID string
	if sessionID, e = idgen.NewHex(32); e != nil {
		return
	}

	sessionHash := shortHash(sessionID)
	expiresAt := now.Add(s.ttl)

	claims := sessionClaims{
		SessionID:   sessionID,
		SessionHash: sessionHash,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	var token string
	if token, e = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret); e != nil {
		return
	}

	return &Token{
		Token:     token,
		Session:   sessionHash,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) Parse(raw string) (_ *Session, e error) {
	if raw = strings.TrimSpace(strings.TrimPrefix(raw, "Bearer ")); raw == "" {
		return nil, apperr.Unauthorized(apperr.CodeSessionRequired, "A valid session JWT is required.")
	}

	var token *jwt.Token
	if token, e = jwt.ParseWithClaims(raw, &sessionClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return &Session{}, apperr.Unauthorized(apperr.CodeSessionRequired, "A valid session JWT is required.")
		}

		return s.secret, nil
	}); e != nil || !token.Valid {
		return &Session{}, apperr.Unauthorized(apperr.CodeSessionRequired, "A valid session JWT is required.")
	}

	claims, ok := token.Claims.(*sessionClaims)
	if !ok || claims.SessionID == "" || claims.SessionHash == "" {
		return &Session{}, apperr.Unauthorized(apperr.CodeSessionRequired, "A valid session JWT is required.")
	}

	return &Session{
		ID:   claims.SessionID,
		Hash: claims.SessionHash,
	}, nil
}

func shortHash(value string) string {
	hash := crypto.SHA256Hex(value)
	if len(hash) > 12 {
		return hash[:12]
	}
	return hash
}
