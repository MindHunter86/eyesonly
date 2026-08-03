import { clearSessionStorage, getSessionToken } from '../stores/session.js';

const API_BASE = import.meta.env.VITE_API_BASE_URL || '';

export class EyesOnlyApiError extends Error {
  constructor(message, { code = 'REQUEST_FAILED', status = 0, details = null } = {}) {
    super(message);
    this.name = 'EyesOnlyApiError';
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

function apiUrl(path) {
  return `${API_BASE}${path}`;
}

function normalizeResponse(body) {
  // Backend contract prefers { ok: true, data: ... }.
  // This fallback makes local integration easier during early backend work.
  if (body && typeof body === 'object' && 'ok' in body) return body.data;
  return body;
}

function isSessionRejected(response, error) {
  if (response.status !== 401 && response.status !== 403) return false;
  if (!error?.code) return true;
  return ['SESSION_REQUIRED', 'SESSION_EXPIRED', 'SESSION_INVALID', 'UNAUTHORIZED'].includes(error.code);
}

async function request(path, options = {}, attempt = 0) {
  const hasBody = options.body !== undefined;
  const sessionToken = await getSessionToken({ forceNew: attempt > 0 });

  const response = await fetch(apiUrl(path), {
    ...options,
    headers: {
      Accept: 'application/json',
      ...(hasBody ? { 'Content-Type': 'application/json' } : {}),
      ...(sessionToken ? { 'X-EyesOnly-Session': sessionToken } : {}),
      ...options.headers,
    },
  });

  const body = await response.json().catch(() => null);

  if (!response.ok || body?.ok === false) {
    const error = body?.error || {};

    if (attempt === 0 && isSessionRejected(response, error)) {
      clearSessionStorage({ includeLegacyData: false });
      return request(path, options, 1);
    }

    throw new EyesOnlyApiError(
      error.message || `Request failed with HTTP ${response.status}`,
      {
        code: error.code || 'REQUEST_FAILED',
        status: response.status,
        details: error.details || null,
      }
    );
  }

  return normalizeResponse(body);
}

export const backendApi = {
  createSecret(payload) {
    return request('/v1/secrets', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  },

  listSessionSecrets() {
    return request('/v1/session/secrets');
  },

  getSecret(id) {
    return request(`/v1/secrets/${encodeURIComponent(id)}`);
  },

  revealSecret(id) {
    return request(`/v1/secrets/${encodeURIComponent(id)}/reveal`, {
      method: 'POST',
    });
  },

  destroySecret(destroyToken) {
    return request('/v1/destroy', {
      method: 'POST',
      body: JSON.stringify({ destroy_token: destroyToken }),
    });
  },

  destroySessionSecret(id) {
    return request(`/v1/secrets/${encodeURIComponent(id)}/destroy`, {
      method: 'POST',
    });
  },

  getStats() {
    return request('/v1/stats');
  },

  listAdminSecrets({ q = '', page = 1, pageSize = 10 } = {}) {
    const params = new URLSearchParams({
      q,
      page: String(page),
      page_size: String(pageSize),
    });
    return request(`/v1/admin/secrets?${params}`);
  },
};
