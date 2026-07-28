CREATE TABLE IF NOT EXISTS secrets (
  id TEXT PRIMARY KEY,
  session_hash TEXT NOT NULL,
  public_description TEXT NOT NULL DEFAULT '',
  content_preview TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  secret_ciphertext TEXT NOT NULL,
  secret_nonce TEXT NOT NULL,
  destroy_token_hash TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  destroyed_at TEXT NULL,
  burn_after_read INTEGER NOT NULL DEFAULT 1,
  max_views INTEGER NOT NULL DEFAULT 1,
  views_used INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_secrets_session_created ON secrets(session_hash, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_secrets_status ON secrets(status);
CREATE INDEX IF NOT EXISTS idx_secrets_expires_at ON secrets(expires_at);
CREATE INDEX IF NOT EXISTS idx_secrets_preview ON secrets(content_preview);
