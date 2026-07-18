# BurnVault backend API contract

This Svelte build uses the real HTTP API instead of in-browser demo data.

## Runtime configuration

```bash
VITE_API_BASE_URL=https://api.example.com
```

If `VITE_API_BASE_URL` is empty, the frontend calls same-origin paths such as `/v1/secrets`.
That is useful when Nginx or Vite proxy forwards `/v1` to the backend.

Every frontend request includes this header:

```http
X-BurnVault-Session: <browser-session-id>
```

The value is generated once in `localStorage`. It is used only for listing links created from the current browser session.
It is not an authentication mechanism for admin APIs.

## Common response format

Success:

```json
{
  "ok": true,
  "data": {}
}
```

Error:

```json
{
  "ok": false,
  "error": {
    "code": "SECRET_NOT_FOUND",
    "message": "Secret was not found or is no longer available.",
    "details": {}
  }
}
```

The frontend also accepts a raw JSON object without `{ ok, data }` during early backend development, but the format above should be treated as canonical.

## POST /v1/secrets

Create one secret.

Request:

```json
{
  "secret": "DATABASE_URL=postgres://...",
  "ttl_seconds": 3600,
  "max_views": 1,
  "burn_after_read": true,
  "hide_from_browser_history": true,
  "public_description": "deployment handoff",
  "notify_email": "ops@example.com"
}
```

Response:

```json
{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "url": "https://example.com/#/s/9f31ab20",
    "destroy_token": "dt_3d4d8c...",
    "status": "active",
    "public_description": "deployment handoff",
    "views_used": 0,
    "max_views": 1,
    "expires_at": "2026-07-18T05:00:00Z"
  }
}
```

Notes:

- The backend should encrypt the secret body before storing it.
- The frontend does not need the encrypted body.
- `destroy_token` is shown only after creation and should not be exposed in list endpoints.

## GET /v1/session/secrets

Return links created by the current browser session.

Request header:

```http
X-BurnVault-Session: <browser-session-id>
```

Response:

```json
{
  "ok": true,
  "data": {
    "items": [
      {
        "id": "9f31ab20",
        "url": "https://example.com/#/s/9f31ab20",
        "status": "active",
        "public_description": "deployment handoff",
        "views_used": 0,
        "max_views": 1,
        "expires_at": "2026-07-18T05:00:00Z"
      }
    ]
  }
}
```

## GET /v1/secrets/{id}

Return public metadata before reveal. Do not return the secret body here.

Response:

```json
{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "status": "active",
    "public_description": "deployment handoff",
    "views_used": 0,
    "max_views": 1,
    "expires_at": "2026-07-18T05:00:00Z"
  }
}
```

## POST /v1/secrets/{id}/reveal

Reveal the secret after the recipient presses the confirmation button.

Response:

```json
{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "secret": "DATABASE_URL=postgres://...",
    "status": "burned",
    "burned": true,
    "views_used": 1,
    "max_views": 1
  }
}
```

Notes:

- This is the only public method that returns the secret body.
- If `burn_after_read` is enabled, the backend should burn the link during this request.
- Repeated calls should return an error such as `SECRET_ALREADY_BURNED`.

## POST /v1/destroy

Destroy by destroy token.

Request:

```json
{
  "destroy_token": "dt_3d4d8c..."
}
```

Response:

```json
{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "status": "burned",
    "destroyed_at": "2026-07-18T04:20:00Z"
  }
}
```

## POST /v1/secrets/{id}/destroy

Destroy a link from the current browser session list.

Request header:

```http
X-BurnVault-Session: <browser-session-id>
```

Response:

```json
{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "status": "burned",
    "destroyed_at": "2026-07-18T04:20:00Z"
  }
}
```

This endpoint exists so the `Recent secrets` table can have a Destroy button without storing destroy tokens in list data.

## GET /v1/stats

Public aggregate counters for the sidebar.

Response:

```json
{
  "ok": true,
  "data": {
    "created": 12840,
    "burned": 12611,
    "avg_ttl": "42m"
  }
}
```

## GET /v1/admin/secrets

Admin listing endpoint.

Query params:

```text
q=<search text>
page=1
page_size=10
```

`q` should search by secret id, session id, and content preview. The content preview should be a safe backend-generated excerpt or redacted value, not full secret content.

Response:

```json
{
  "ok": true,
  "data": {
    "items": [
      {
        "id": "9f31ab20",
        "session": "a81720cc",
        "status": "active",
        "ttl": "51m",
        "views": "0/1",
        "preview": "DATABASE_URL secret"
      }
    ],
    "page": 1,
    "page_size": 10,
    "total": 42
  }
}
```

## Recommended error codes

```text
VALIDATION_ERROR
SECRET_NOT_FOUND
SECRET_EXPIRED
SECRET_ALREADY_BURNED
DESTROY_TOKEN_INVALID
RATE_LIMITED
ADMIN_AUTH_REQUIRED
INTERNAL_ERROR
```
