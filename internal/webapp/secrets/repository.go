package secrets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MindHunter86/eyesonly/internal/webapp/shared/apperr"
)

type Store interface {
	Create(ctx context.Context, secret *Secret) error
	GetByID(ctx context.Context, id string) (*Secret, error)
	ListBySession(ctx context.Context, sessionHash string, now time.Time) ([]*Secret, error)
	Reveal(ctx context.Context, id string, now time.Time) (*Secret, error)
	BurnByDestroyTokenHash(ctx context.Context, tokenHash string, now time.Time) (*Secret, error)
	BurnByIDAndSession(ctx context.Context, id string, sessionHash string, now time.Time) (*Secret, error)
	Stats(ctx context.Context, now time.Time) (*Stats, error)
	ListAdmin(ctx context.Context, filter *AdminFilter, now time.Time) (*AdminList, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, secret *Secret) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO secrets (
  id, session_hash, public_description, content_preview, status,
  secret_ciphertext, secret_nonce, destroy_token_hash,
  created_at, expires_at, destroyed_at,
  burn_after_read, max_views, views_used
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		secret.ID,
		secret.SessionHash,
		secret.PublicDescription,
		secret.ContentPreview,
		secret.Status,
		secret.SecretCiphertext,
		secret.SecretNonce,
		secret.DestroyTokenHash,
		formatTime(secret.CreatedAt),
		formatTime(secret.ExpiresAt),
		formatOptionalTime(secret.DestroyedAt),
		boolInt(secret.BurnAfterRead),
		secret.MaxViews,
		secret.ViewsUsed,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Secret, error) {
	return r.getOne(ctx, `SELECT `+secretColumns()+` FROM secrets WHERE id = ?`, id)
}

func (r *Repository) ListBySession(ctx context.Context, sessionHash string, now time.Time) (_ []*Secret, e error) {
	var rows *sql.Rows
	if rows, e = r.db.QueryContext(ctx, `SELECT `+secretColumns()+` FROM secrets WHERE session_hash = ? ORDER BY created_at DESC LIMIT 100`, sessionHash); e != nil {
		return
	}
	defer rows.Close()

	return scanSecrets(rows, now)
}

func (r *Repository) Reveal(ctx context.Context, id string, now time.Time) (_ *Secret, e error) {
	var tx *sql.Tx
	if tx, e = r.db.BeginTx(ctx, nil); e != nil {
		return
	}
	defer tx.Rollback()

	var secret *Secret
	if secret, e = getOneTx(ctx, tx, `SELECT `+secretColumns()+` FROM secrets WHERE id = ?`, id); e != nil {
		return
	}

	switch secret.EffectiveStatus(now) {
	case StatusExpired:
		_ = markStatusTx(ctx, tx, id, StatusExpired, now)
		return nil, apperr.Gone(apperr.CodeSecretExpired, "Secret has expired.")
	case StatusBurned:
		return nil, apperr.Gone(apperr.CodeSecretAlreadyBurned, "Secret has already been burned.")
	}

	secret.ViewsUsed++
	if secret.BurnAfterRead || secret.ViewsUsed >= secret.MaxViews {
		secret.Status = StatusBurned
		secret.DestroyedAt = &now
	}

	_, e = tx.ExecContext(ctx, `UPDATE secrets SET views_used = ?, status = ?, destroyed_at = ? WHERE id = ?`,
		secret.ViewsUsed,
		secret.Status,
		formatOptionalTime(secret.DestroyedAt),
		secret.ID,
	)
	if e != nil {
		return
	}

	return secret, tx.Commit()
}

func (r *Repository) BurnByDestroyTokenHash(ctx context.Context, tokenHash string, now time.Time) (_ *Secret, e error) {
	var secret *Secret
	if secret, e = r.getOne(ctx, `SELECT `+secretColumns()+` FROM secrets WHERE destroy_token_hash = ?`, tokenHash); e != nil {
		var appErr *apperr.Error
		if errors.As(e, &appErr) && appErr.Code == apperr.CodeSecretNotFound {
			return nil, apperr.InvalidDestroyToken()
		}

		return
	}

	return r.burn(ctx, secret.ID, now)
}

func (r *Repository) BurnByIDAndSession(ctx context.Context, id string, sessionHash string, now time.Time) (_ *Secret, e error) {
	var secret *Secret
	if secret, e = r.getOne(ctx, `SELECT `+secretColumns()+` FROM secrets WHERE id = ? AND session_hash = ?`, id, sessionHash); e != nil {
		return
	}

	return r.burn(ctx, secret.ID, now)
}

func (r *Repository) burn(ctx context.Context, id string, now time.Time) (_ *Secret, e error) {
	var secret *Secret
	if secret, e = r.GetByID(ctx, id); e != nil {
		return
	} else if secret.Status == StatusBurned {
		return secret, nil
	}

	secret.Status = StatusBurned
	secret.DestroyedAt = &now

	_, e = r.db.ExecContext(ctx, `UPDATE secrets SET status = ?, destroyed_at = ? WHERE id = ?`, StatusBurned, formatTime(now), id)
	if e != nil {
		return
	}

	return secret, nil
}

func (r *Repository) Stats(ctx context.Context, now time.Time) (_ *Stats, e error) {
	var created int
	var burned sql.NullInt64
	var avgSeconds sql.NullFloat64

	e = r.db.QueryRowContext(ctx, `
SELECT
  COUNT(*),
  SUM(CASE WHEN status = 'burned' THEN 1 ELSE 0 END),
  AVG((julianday(expires_at) - julianday(created_at)) * 86400.0)
FROM secrets`).Scan(&created, &burned, &avgSeconds)
	if e != nil {
		return
	}

	return &Stats{
		Created: created,
		Burned:  int(burned.Int64),
		AvgTTL:  formatDurationSeconds(avgSeconds.Float64),
	}, nil
}

func (r *Repository) ListAdmin(ctx context.Context, filter *AdminFilter, now time.Time) (_ *AdminList, e error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	q := strings.TrimSpace(filter.Q)
	like := "%" + q + "%"
	where := `WHERE (? = '' OR id LIKE ? OR session_hash LIKE ? OR content_preview LIKE ? OR public_description LIKE ?)`
	args := []any{q, like, like, like, like}

	var total int
	if e = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM secrets `+where, args...).Scan(&total); e != nil {
		return
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := `SELECT ` + secretColumns() + ` FROM secrets ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.PageSize, offset)

	var rows *sql.Rows
	if rows, e = r.db.QueryContext(ctx, query, args...); e != nil {
		return
	}
	defer rows.Close()

	var secrets []*Secret
	if secrets, e = scanSecrets(rows, now); e != nil {
		return
	}

	items := make([]*AdminSecret, 0, len(secrets))
	for _, secret := range secrets {
		status := secret.EffectiveStatus(now)
		items = append(items, &AdminSecret{
			ID:      secret.ID,
			Session: secret.SessionHash,
			Status:  status,
			TTL:     ttlText(secret, now),
			Views:   fmt.Sprintf("%d/%d", secret.ViewsUsed, secret.MaxViews),
			Preview: secret.ContentPreview,
		})
	}

	return &AdminList{
		Items:    items,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Total:    total,
	}, nil
}

func (r *Repository) getOne(ctx context.Context, query string, args ...any) (*Secret, error) {
	return getOneDB(ctx, r.db, query, args...)
}

func getOneDB(ctx context.Context, db *sql.DB, query string, args ...any) (*Secret, error) {
	row := db.QueryRowContext(ctx, query, args...)
	return scanSecret(row)
}

func getOneTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*Secret, error) {
	row := tx.QueryRowContext(ctx, query, args...)
	return scanSecret(row)
}

// !!! TODO : BUGCHECK!!!
// !!! TODO : BUGCHECK!!!
// !!! TODO : BUGCHECK!!!
func scanSecrets(rows *sql.Rows, _ time.Time) ([]*Secret, error) {
	items := make([]*Secret, 0)

	for rows.Next() {
		secret, err := scanSecret(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, secret)
	}

	return items, rows.Err()
}

type rowScanner interface{ Scan(dest ...any) error }

func scanSecret(row rowScanner) (_ *Secret, e error) {
	var secret = new(Secret)
	var createdAt string
	var expiresAt string
	var destroyedAt sql.NullString
	var burnAfterRead int

	e = row.Scan(
		&secret.ID,
		&secret.SessionHash,
		&secret.PublicDescription,
		&secret.ContentPreview,
		&secret.Status,
		&secret.SecretCiphertext,
		&secret.SecretNonce,
		&secret.DestroyTokenHash,
		&createdAt,
		&expiresAt,
		&destroyedAt,
		&burnAfterRead,
		&secret.MaxViews,
		&secret.ViewsUsed,
	)
	if errors.Is(e, sql.ErrNoRows) {
		return nil, apperr.NotFound("Secret was not found or is no longer available.")
	} else if e != nil {
		return
	}

	secret.CreatedAt, secret.ExpiresAt =
		parseTime(createdAt), parseTime(expiresAt)

	if destroyedAt.Valid && destroyedAt.String != "" {
		parsed := parseTime(destroyedAt.String)
		secret.DestroyedAt = &parsed
	}

	secret.BurnAfterRead = burnAfterRead == 1

	return secret, nil
}

func secretColumns() string {
	return `id, session_hash, public_description, content_preview, status, secret_ciphertext, secret_nonce, destroy_token_hash, created_at, expires_at, destroyed_at, burn_after_read, max_views, views_used`
}

func markStatusTx(ctx context.Context, tx *sql.Tx, id string, status string, now time.Time) error {
	_, e := tx.ExecContext(ctx, `UPDATE secrets SET status = ?, destroyed_at = ? WHERE id = ?`, status, formatTime(now), id)
	return e
}

func boolInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func formatOptionalTime(value *time.Time) any {
	if value == nil {
		return nil
	}

	return formatTime(*value)
}

func parseTime(value string) time.Time {
	parsed, e := time.Parse(time.RFC3339, value)
	if e != nil {
		return time.Time{}
	}
	return parsed
}

func ttlText(secret *Secret, now time.Time) string {
	var status string
	if status = secret.EffectiveStatus(now); status != StatusActive {
		return status
	}

	var remain time.Duration
	if remain = secret.ExpiresAt.Sub(now); remain < 0 {
		remain = 0
	}

	return formatDuration(remain)
}

func formatDurationSeconds(seconds float64) string {
	if seconds <= 0 {
		return "0m"
	}

	return formatDuration(time.Duration(seconds) * time.Second)
}

func formatDuration(value time.Duration) string {
	if value >= 24*time.Hour {
		return fmt.Sprintf("%dd", int(value.Hours()/24))
	} else if value >= time.Hour {
		return fmt.Sprintf("%dh", int(value.Hours()))
	}

	var minutes int
	if minutes = int(value.Minutes()); minutes <= 0 {
		return "<1m"
	}

	return fmt.Sprintf("%dm", minutes)
}
