package secrets

import "time"

const (
	StatusActive  = "active"
	StatusBurned  = "burned"
	StatusExpired = "expired"
)

type Secret struct {
	ID                string
	SessionHash       string
	PublicDescription string
	ContentPreview    string
	Status            string
	SecretCiphertext  string
	SecretNonce       string
	DestroyTokenHash  string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	DestroyedAt       *time.Time
	BurnAfterRead     bool
	MaxViews          int
	ViewsUsed         int
}

func (s Secret) EffectiveStatus(now time.Time) string {
	if s.Status == StatusBurned {
		return StatusBurned
	}
	if !s.ExpiresAt.IsZero() && now.After(s.ExpiresAt) {
		return StatusExpired
	}
	return StatusActive
}

type CreateInput struct {
	Secret                 string `json:"secret"`
	TTLSeconds             int    `json:"ttl_seconds"`
	MaxViews               int    `json:"max_views"`
	BurnAfterRead          bool   `json:"burn_after_read"`
	HideFromBrowserHistory bool   `json:"hide_from_browser_history"`
	PublicDescription      string `json:"public_description"`
	NotifyEmail            string `json:"notify_email"`
}

type Public struct {
	ID                string `json:"id"`
	URL               string `json:"url,omitempty"`
	Status            string `json:"status"`
	PublicDescription string `json:"public_description"`
	ViewsUsed         int    `json:"views_used"`
	MaxViews          int    `json:"max_views"`
	ExpiresAt         string `json:"expires_at"`
	TTL               string `json:"ttl,omitempty"`
}

type Created struct {
	*Public
	DestroyToken string `json:"destroy_token"`
}

type Reveal struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`
	Status    string `json:"status"`
	Burned    bool   `json:"burned"`
	ViewsUsed int    `json:"views_used"`
	MaxViews  int    `json:"max_views"`
}

type DestroyResult struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	DestroyedAt string `json:"destroyed_at"`
}

type Stats struct {
	Created int    `json:"created"`
	Burned  int    `json:"burned"`
	AvgTTL  string `json:"avg_ttl"`
}

type AdminSecret struct {
	ID      string `json:"id"`
	Session string `json:"session"`
	Status  string `json:"status"`
	TTL     string `json:"ttl"`
	Views   string `json:"views"`
	Preview string `json:"preview"`
}

type AdminList struct {
	Items    []*AdminSecret `json:"items"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Total    int            `json:"total"`
}

type AdminFilter struct {
	Q        string
	Page     int
	PageSize int
}

const MAX_TTL = 60 * 60 * 24 * 7
const MAX_VIEW = 25
