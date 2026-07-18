const SESSION_STORAGE_KEY = 'burnvault.session.id';

function randomHex(byteCount) {
  const bytes = new Uint8Array(byteCount);
  crypto.getRandomValues(bytes);
  return Array.from(bytes, byte => byte.toString(16).padStart(2, '0')).join('');
}

function getStorage() {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

export function getSessionId() {
  const storage = getStorage();
  if (!storage) return 'session-unavailable';

  const existing = storage.getItem(SESSION_STORAGE_KEY);
  if (existing) return existing;

  const created = randomHex(8);
  storage.setItem(SESSION_STORAGE_KEY, created);
  return created;
}

export function resetSessionId() {
  const storage = getStorage();
  if (!storage) return 'session-unavailable';

  const created = randomHex(8);
  storage.setItem(SESSION_STORAGE_KEY, created);
  return created;
}
