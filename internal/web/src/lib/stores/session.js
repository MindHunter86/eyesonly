const SESSION_TOKEN_KEY = 'eyesonly.session.jwt';
const SESSION_HASH_KEY = 'eyesonly.session.hash';
const API_BASE = import.meta.env.VITE_API_BASE_URL || '';
const JWT_EXPIRY_SKEW_MS = 30000;

const LEGACY_STORAGE_KEYS = [
  SESSION_TOKEN_KEY,
  SESSION_HASH_KEY,
  'eyesonly.session.id',
  'eyesonly.sessionId',
  'eyesonly.session.token',
  'eyesonly.client.session',
  'eyesonly.clientSession',
];

const LEGACY_SESSION_STORAGE_KEYS = [
  'eyesonly.session.secrets',
];

let pendingSessionRequest = null;

function getStorage() {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

function getSessionStorage() {
  try {
    return window.sessionStorage;
  } catch {
    return null;
  }
}

function normalizeResponse(body) {
  if (body && typeof body === 'object' && 'ok' in body) return body.data;
  return body;
}

function decodeJwtPayload(token) {
  const payload = token?.split?.('.')[1];
  if (!payload) return null;

  try {
    const padded = payload.replace(/-/g, '+').replace(/_/g, '/').padEnd(Math.ceil(payload.length / 4) * 4, '=');
    return JSON.parse(atob(padded));
  } catch {
    return null;
  }
}

function isTokenExpired(token) {
  const payload = decodeJwtPayload(token);
  if (!payload?.exp) return true;
  return payload.exp * 1000 <= Date.now() + JWT_EXPIRY_SKEW_MS;
}

function saveSession(data) {
  const storage = getStorage();
  if (!storage) return data.token;

  storage.setItem(SESSION_TOKEN_KEY, data.token);

  const sessionHash = data.session || data.session_hash || data.sessionHash || '';
  if (sessionHash) {
    storage.setItem(SESSION_HASH_KEY, sessionHash);
  } else {
    storage.removeItem(SESSION_HASH_KEY);
  }

  return data.token;
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

export async function getSessionToken({ forceNew = false } = {}) {
  const storage = getStorage();
  if (!storage) return '';

  const existing = storage.getItem(SESSION_TOKEN_KEY);
  if (!forceNew && existing && !isTokenExpired(existing)) return existing;

  if (forceNew || existing) {
    clearSessionStorage({ includeLegacyData: false });
  }

  if (!pendingSessionRequest) {
    pendingSessionRequest = issueSessionToken()
      .then(saveSession)
      .finally(() => {
        pendingSessionRequest = null;
      });
  }

  return pendingSessionRequest;
}

export function getSessionHash() {
  const token = getStorage()?.getItem(SESSION_TOKEN_KEY) || '';
  if (token && isTokenExpired(token)) {
    clearSessionStorage({ includeLegacyData: false });
    return '';
  }

  return getStorage()?.getItem(SESSION_HASH_KEY) || '';
}

export function clearSessionStorage({ includeLegacyData = true } = {}) {
  const storage = getStorage();
  if (storage) {
    for (const key of LEGACY_STORAGE_KEYS) {
      storage.removeItem(key);
    }
  }

  if (includeLegacyData) {
    const sessionStorage = getSessionStorage();
    if (sessionStorage) {
      for (const key of LEGACY_SESSION_STORAGE_KEYS) {
        sessionStorage.removeItem(key);
      }
    }
  }

  pendingSessionRequest = null;
}

export function resetSessionId() {
  clearSessionStorage();
  return '';
}
