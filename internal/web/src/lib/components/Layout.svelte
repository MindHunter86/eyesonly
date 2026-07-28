<script>
  import { onMount } from 'svelte';
  import { theme, toggleTheme } from '../stores/theme.js';
  import { stats, loadStats } from '../stores/stats.js';
  import logoUrl from '../../assets/logo.svg';

  export let route = 'home';
  export let onNavigate;
  export let onClearSession;
  export let onScrollCreate;

  const navItems = [
    ['home', 'Home'],
    ['contacts', 'Contacts'],
    ['news', 'News'],
    ['destroy', 'Destroy'],
    ['donate', 'Donate'],
    ['api', 'API'],
    ['admin', 'Admin'],
    // ['ui', 'UI'],
    ['about', 'About'],
  ];

  function navigate(routeName) {
    onNavigate(routeName);
  }

  function formatStat(value, fallback = '—') {
    if (value === null || value === undefined || value === '') return fallback;
    return typeof value === 'number' ? value.toLocaleString() : value;
  }

  onMount(() => {
    loadStats().catch(() => {
      // Sidebar stats are non-critical. Page actions show their own errors.
    });
  });
</script>

<div class="app-shell">
  <aside class="sidebar" aria-label="Primary navigation">
    <a class="brand" href="#home" aria-label="EyesOnly home" on:click|preventDefault={() => navigate('home')}>
      <div class="brand-mark" aria-hidden="true">
        <img class="brand-logo" src={logoUrl} alt="" />
      </div>
      <div>
        <div class="brand-name">EyesOnly</div>
        <div class="brand-meta">private links</div>
      </div>
    </a>

    <nav class="nav">
      {#each navItems as [key, label]}
        <a
          class="nav-link"
          class:active={route === key}
          href={`#${key}`}
          on:click|preventDefault={() => navigate(key)}
        >{label}</a>
      {/each}
    </nav>

    <div class="sidebar-card">
      <div class="card-title">Next</div>
      <div class="roadmap-item">
        <span class="roadmap-dot"></span>
        <div>
          <strong>Access policy</strong>
          <span>role limits and passphrase rules</span>
        </div>
      </div>
      <div class="roadmap-item">
        <span class="roadmap-dot"></span>
        <div>
          <strong>Audit view</strong>
          <span>link state and destroy events</span>
        </div>
      </div>
      <div class="roadmap-note">Secret bodies are returned only after reveal.</div>
      <div class="global-stats" aria-label="Global service statistics">
        <div>
          <span>created</span>
          <strong>{formatStat($stats.created)}</strong>
        </div>
        <div>
          <span>burned</span>
          <strong>{formatStat($stats.burned)}</strong>
        </div>
        <div>
          <span>avg ttl</span>
          <strong>{formatStat($stats.avgTtl)}</strong>
        </div>
      </div>
    </div>
  </aside>

  <main class="workspace">
    <header class="topbar">
      <div>
        <p class="eyebrow">secure / temporary / disposable</p>
        <h1>One-time secret delivery</h1>
      </div>
      <div class="topbar-actions">
        <button
          class="ghost-button icon-button"
          type="button"
          aria-label={$theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
          title={$theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
          on:click={toggleTheme}
        >{$theme === 'dark' ? '☼' : '☾'}</button>
        <button class="ghost-button" type="button" on:click={onClearSession}>
          Clear session
        </button>
        <button class="primary-button" type="button" on:click={onScrollCreate}>
          New secret
        </button>
      </div>
    </header>

    <slot />

    <footer class="site-footer">
      <span>© 2026 EyesOnly</span>
      <span>Disposable secret sharing frontend MVP</span>
    </footer>
  </main>
</div>
