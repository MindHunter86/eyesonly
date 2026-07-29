const SESSION_TOKEN_KEY = 'eyesonly.session.jwt';
const SESSION_HASH_KEY = 'eyesonly.session.hash';
const API_BASE = import.meta.env.VITE_API_BASE_URL || '';

let pendingSessionRequest = null;

function getStorage() {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

function normalizeResponse(body) {
  if (body && typeof body === 'object' && 'ok' in body) return body.data;
  return body;
}

async function issueSessionToken() {
  const response = await fetch(`${API_BASE}/v1/session`, {
    method: 'POST',
    headers: { Accept: 'application/json' },
  });

  const body = await response.json().catch(() => null);

  if (!response.ok || body?.ok === false) {
    const error = body?.error || {};
    throw new Error(error.message || `Session request failed with HTTP ${response.status}`);
  }

  const data = normalizeResponse(body);
  if (!data?.token) {
    throw new Error('Session token was not returned by backend.');
  }

  return data;
}

export async function getSessionToken() {
  const storage = getStorage();
  if (!storage) return '';

  const existing = storage.getItem(SESSION_TOKEN_KEY);
  if (existing) return existing;

  if (!pendingSessionRequest) {
    pendingSessionRequest = issueSessionToken()
      .then(data => {
        storage.setItem(SESSION_TOKEN_KEY, data.token);
        if (data.session) storage.setItem(SESSION_HASH_KEY, data.session);
        return data.token;
      })
      .finally(() => {
        pendingSessionRequest = null;
      });
  }

  return pendingSessionRequest;
}

export function getSessionHash() {
  return getStorage()?.getItem(SESSION_HASH_KEY) || '';
}

export function resetSessionId() {
  const storage = getStorage();
  if (!storage) return '';

  storage.removeItem(SESSION_TOKEN_KEY);
  storage.removeItem(SESSION_HASH_KEY);
  pendingSessionRequest = null;
  return '';
}
