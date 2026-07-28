import { get, writable } from 'svelte/store';
import { backendApi } from '../services/api.js';

export const secrets = writable([]);
export const secretsLoading = writable(false);
export const secretsError = writable('');

export function now() {
  return Date.now();
}

function toDateMs(value) {
  if (!value) return null;
  const timestamp = Date.parse(value);
  return Number.isNaN(timestamp) ? null : timestamp;
}

function normalizeStatus(value) {
  if (!value) return 'active';
  if (value === 'destroyed') return 'burned';
  return value;
}

export function normalizeSecret(raw) {
  const viewsUsed = raw.views_used ?? raw.viewsUsed ?? raw.views ?? 0;
  const maxViews = raw.max_views ?? raw.maxViews ?? 1;
  const expiresAt = toDateMs(raw.expires_at ?? raw.expiresAt);

  return {
    id: raw.id,
    url: raw.url,
    destroyToken: raw.destroy_token ?? raw.destroyToken ?? '',
    status: normalizeStatus(raw.status),
    maxViews,
    views: viewsUsed,
    burnAfterRead: Boolean(raw.burn_after_read ?? raw.burnAfterRead),
    hint: raw.public_description ?? raw.publicDescription ?? raw.hint ?? '',
    notifyEmail: raw.notify_email ?? raw.notifyEmail ?? '',
    createdAt: toDateMs(raw.created_at ?? raw.createdAt) ?? now(),
    expiresAt,
    ttl: raw.ttl,
  };
}

export function buildSecretUrl(id, url = '') {
  if (url) return url;
  return `${location.origin}${location.pathname}#/s/${id}`;
}

export function getStatus(secret) {
  if (!secret) return 'active';
  if (secret.status) return normalizeStatus(secret.status);
  if (secret.expiresAt && now() > secret.expiresAt) return 'expired';
  return 'active';
}

export function statusClass(status) {
  if (status === 'active') return 'success';
  if (status === 'burned') return 'danger';
  return '';
}

export function formatTtl(secret) {
  const status = getStatus(secret);
  if (status === 'burned' || status === 'expired') return status;
  if (secret.ttl) return secret.ttl;
  if (!secret.expiresAt) return 'server';

  const seconds = Math.max(0, Math.floor((secret.expiresAt - now()) / 1000));
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);

  if (hours > 0) return `${hours}h ${minutes % 60}m`;
  return `${minutes}m ${seconds % 60}s`;
}

export function findSecret(id) {
  return get(secrets).find(secret => secret.id === id) || null;
}

export async function loadSessionSecrets() {
  secretsLoading.set(true);
  secretsError.set('');

  try {
    const data = await backendApi.listSessionSecrets();
    const items = Array.isArray(data) ? data : data.items || [];
    secrets.set(items.map(normalizeSecret));
  } catch (error) {
    secretsError.set(error.message || 'Failed to load session secrets.');
    throw error;
  } finally {
    secretsLoading.set(false);
  }
}

export async function createSecret(payload) {
  const data = await backendApi.createSecret({
    secret: payload.secret,
    ttl_seconds: Number(payload.ttl || 3600),
    max_views: Math.max(1, Number(payload.maxViews || 1)),
    burn_after_read: Boolean(payload.burnAfterRead),
    hide_from_browser_history: Boolean(payload.hideFromHistory),
    public_description: (payload.hint || '').trim(),
    notify_email: payload.notifyOnOpen ? (payload.notifyEmail || '').trim() : '',
  });

  const created = normalizeSecret(data);
  secrets.update(items => [created, ...items.filter(item => item.id !== created.id)]);
  return created;
}

export async function destroyById(id) {
  const data = await backendApi.destroySessionSecret(id);
  const updated = normalizeSecret({ ...data, id, status: data.status || 'burned' });
  secrets.update(items => items.map(item => item.id === id ? { ...item, ...updated } : item));
  return updated;
}

export async function destroyByToken(token) {
  const data = await backendApi.destroySecret(token.trim());
  const id = data.id;

  if (id) {
    const updated = normalizeSecret({ ...data, id, status: data.status || 'burned' });
    secrets.update(items => items.map(item => item.id === id ? { ...item, ...updated } : item));
  }

  return data;
}

export function clearSecrets() {
  secrets.set([]);
  secretsError.set('');
}
