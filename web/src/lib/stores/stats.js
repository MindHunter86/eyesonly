import { writable } from 'svelte/store';
import { backendApi } from '../services/api.js';

export const stats = writable({
  created: null,
  burned: null,
  avgTtl: null,
});

export async function loadStats() {
  const data = await backendApi.getStats();
  stats.set({
    created: data.created ?? data.total_created ?? null,
    burned: data.burned ?? data.total_burned ?? null,
    avgTtl: data.avg_ttl ?? data.avgTtl ?? null,
  });
  return data;
}
