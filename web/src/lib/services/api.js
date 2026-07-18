import { getSessionId } from '../stores/session.js';

const API_BASE = import.meta.env.VITE_API_BASE_URL || '';

export class BurnVaultApiError extends Error {
  constructor(message, { code = 'REQUEST_FAILED', status = 0, details = null } = {}) {
    super(message);
    this.name = 'BurnVaultApiError';
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

async function request(path, options = {}) {
  const hasBody = options.body !== undefined;
  const response = await fetch(apiUrl(path), {
    ...options,
    headers: {
      Accept: 'application/json',
      ...(hasBody ? { 'Content-Type': 'application/json' } : {}),
      'X-BurnVault-Session': getSessionId(),
      ...options.headers,
    },
  });

  const body = await response.json().catch(() => null);

  if (!response.ok || body?.ok === false) {
    const error = body?.error || {};
    throw new BurnVaultApiError(
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
