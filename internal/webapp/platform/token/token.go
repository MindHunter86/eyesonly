package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role"`
	SID  string `json:"sid"` // session id
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(c context.Context) *JWTManager {
	// context work
	return &JWTManager{}
}

func (m *JWTManager) Create(uid, role, sid string) (string, error) {
	now := time.Now().UTC()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Role: role,
		SID:  sid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uid,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}).SignedString(m.secret)
}

// todo - errors.New - remove
func (m *JWTManager) Parse(payload string) (_ *Claims, e error) {
	claims := &Claims{}

	var tkn *jwt.Token
	if tkn, e = jwt.ParseWithClaims(payload, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}); e != nil {
		return
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// !! TODO - OPTIMIZE!! Perfomance LEAK!!
func NewOpaqueToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// !! TODO - OPTIMIZE!! Perfomance LEAK!!
func HashOpaqueToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
