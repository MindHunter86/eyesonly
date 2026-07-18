package secrets

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/MindHunter86/eyesonly/internal/webapp/platform/crypto"
	"github.com/MindHunter86/eyesonly/internal/webapp/platform/idgen"
	"github.com/MindHunter86/eyesonly/internal/webapp/platform/text"
	"github.com/MindHunter86/eyesonly/internal/webapp/sessions"
	"github.com/MindHunter86/eyesonly/internal/webapp/shared/apperr"
)

type Service struct {
	store         Store
	crypto        *crypto.Service
	publicBaseURL string
}

func NewService(store Store, crypto *crypto.Service, publicBaseURL string) *Service {
	return &Service{store: store, crypto: crypto, publicBaseURL: strings.TrimRight(publicBaseURL, "/")}
}

func (s *Service) Create(ctx context.Context, session *sessions.Session, input CreateInput, now time.Time) (_ *Created, e error) {
	var plain string
	if plain = strings.TrimSpace(input.Secret); plain == "" {
		return nil, apperr.Validation("secret is required")
	}

	var ttl = input.TTLSeconds
	if ttl <= 0 {
		ttl = 3600
	} else if ttl > MAX_TTL {
		ttl = MAX_TTL
	}

	var maxViews = input.MaxViews
	if maxViews <= 0 {
		maxViews = 1
	} else if maxViews > MAX_VIEW {
		maxViews = MAX_VIEW
	}

	var id, destroyTokenRaw string
	if id, e = idgen.NewHex(4); e != nil {
		return

	}
	if destroyTokenRaw, e = idgen.NewHex(32); e != nil {
		return
	}
	destroyToken := "dt_" + destroyTokenRaw
	cipherText := s.crypto.Encrypt([]byte(plain))

	secret := &Secret{
		ID:                id,
		SessionHash:       session.Hash,
		PublicDescription: strings.TrimSpace(input.PublicDescription),
		ContentPreview:    text.Preview(plain, 64),
		Status:            StatusActive,
		SecretCiphertext:  string(cipherText),
		DestroyTokenHash:  crypto.SHA256Hex(destroyToken),
		CreatedAt:         now,
		ExpiresAt:         now.Add(time.Duration(ttl) * time.Second),
		BurnAfterRead:     input.BurnAfterRead,
		MaxViews:          maxViews,
		ViewsUsed:         0,
	}

	if e = s.store.Create(ctx, secret); e != nil {
		return
	}

	return &Created{
		Public:       s.toPublic(secret, now, true),
		DestroyToken: destroyToken,
	}, nil
}

func (s *Service) ListSession(ctx context.Context, session *sessions.Session, now time.Time) (_ []Public, e error) {
	var items []*Secret
	if items, e = s.store.ListBySession(ctx, session.Hash, now); e != nil {
		return
	}

	out := make([]Public, 0, len(items))
	for _, item := range items {
		out = append(out, *s.toPublic(item, now, true))
	}

	return out, nil
}

func (s *Service) GetPublic(ctx context.Context, id string, now time.Time) (_ *Public, e error) {
	var item *Secret
	if item, e = s.store.GetByID(ctx, id); e != nil {
		return
	}

	return s.toPublic(item, now, false), nil
}

func (s *Service) Reveal(ctx context.Context, id string, now time.Time) (_ *Reveal, e error) {
	var item *Secret
	if item, e = s.store.Reveal(ctx, id, now); e != nil {
		return
	}

	plain := s.crypto.Decrypt([]byte(item.SecretCiphertext))
	status := item.EffectiveStatus(now)

	return &Reveal{
		ID:        item.ID,
		Secret:    string(plain),
		Status:    status,
		Burned:    status == StatusBurned,
		ViewsUsed: item.ViewsUsed,
		MaxViews:  item.MaxViews,
	}, nil
}

func (s *Service) DestroyByToken(ctx context.Context, destroyToken string, now time.Time) (_ *DestroyResult, e error) {
	if destroyToken = strings.TrimSpace(destroyToken); destroyToken == "" {
		return nil, apperr.InvalidDestroyToken()
	}

	var item *Secret
	if item, e = s.store.BurnByDestroyTokenHash(ctx, crypto.SHA256Hex(destroyToken), now); e != nil {
		return
	}

	return toDestroyResult(item, now), nil
}

func (s *Service) DestroyBySession(ctx context.Context, session *sessions.Session, id string, now time.Time) (_ *DestroyResult, e error) {
	var item *Secret
	if item, e = s.store.BurnByIDAndSession(ctx, id, session.Hash, now); e != nil {
		return
	}

	return toDestroyResult(item, now), nil
}

func (s *Service) Stats(ctx context.Context, now time.Time) (*Stats, error) {
	return s.store.Stats(ctx, now)
}

func (s *Service) ListAdmin(ctx context.Context, filter *AdminFilter, now time.Time) (*AdminList, error) {
	return s.store.ListAdmin(ctx, filter, now)
}

func (s *Service) toPublic(secret *Secret, now time.Time, includeURL bool) *Public {
	out := &Public{
		ID:                secret.ID,
		Status:            secret.EffectiveStatus(now),
		PublicDescription: secret.PublicDescription,
		ViewsUsed:         secret.ViewsUsed,
		MaxViews:          secret.MaxViews,
		ExpiresAt:         secret.ExpiresAt.UTC().Format(time.RFC3339),
		TTL:               ttlText(secret, now),
	}

	if includeURL {
		out.URL = s.secretURL(secret.ID)
	}

	return out
}

func (s *Service) secretURL(id string) string {
	return fmt.Sprintf("%s/#/s/%s", s.publicBaseURL, url.PathEscape(id))
}

func toDestroyResult(secret *Secret, now time.Time) *DestroyResult {
	destroyedAt := now.UTC().Format(time.RFC3339)
	if secret.DestroyedAt != nil {
		destroyedAt = secret.DestroyedAt.UTC().Format(time.RFC3339)
	}

	return &DestroyResult{
		ID:          secret.ID,
		Status:      StatusBurned,
		DestroyedAt: destroyedAt,
	}
}
