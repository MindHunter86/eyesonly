import { writable, get } from 'svelte/store';

const THEME_KEY = 'eyesonly.theme';

function readTheme() {
  if (typeof localStorage === 'undefined') return 'dark';
  const saved = localStorage.getItem(THEME_KEY);
  return saved === 'light' || saved === 'dark' ? saved : 'dark';
}

export const theme = writable(readTheme());

export function applyTheme(nextTheme) {
  theme.set(nextTheme);
  if (typeof document !== 'undefined') {
    document.documentElement.dataset.theme = nextTheme;
  }
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(THEME_KEY, nextTheme);
  }
}

export function initTheme() {
  applyTheme(readTheme());
}

export function toggleTheme() {
  applyTheme(get(theme) === 'dark' ? 'light' : 'dark');
}
