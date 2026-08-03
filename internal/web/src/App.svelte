<script>
  import { onMount, tick } from 'svelte';
  import './app.css';

  import Layout from './lib/components/Layout.svelte';
  import Toast from './lib/components/Toast.svelte';
  import Home from './lib/pages/Home.svelte';
  import Contacts from './lib/pages/Contacts.svelte';
  import News from './lib/pages/News.svelte';
  import NewsArticle from './lib/pages/NewsArticle.svelte';
  import RecipientPreview from './lib/pages/RecipientPreview.svelte';
  import Destroy from './lib/pages/Destroy.svelte';
  import Donate from './lib/pages/Donate.svelte';
  import Api from './lib/pages/Api.svelte';
  import Admin from './lib/pages/Admin.svelte';
  // import UiKit from './lib/pages/UiKit.svelte';
  import About from './lib/pages/About.svelte';
  import { clearSecrets } from './lib/stores/secrets.js';
  import { resetSessionId } from './lib/stores/session.js';
  import { initTheme } from './lib/stores/theme.js';
  import { notify } from './lib/stores/toast.js';

  let route = 'home';
  let secretId = null;

  function parseHash() {
    const hash = location.hash.replace('#', '');
    if (hash.startsWith('/s/')) {
      secretId = hash.replace('/s/', '') || null;
      route = 'secret-preview';
      return;
    }

    secretId = null;
    route = hash || 'home';
  }

  function goTo(nextRoute) {
    route = nextRoute;
    secretId = null;
    history.replaceState(null, '', `#${nextRoute}`);
  }

  function openSecretPreview(id) {
    secretId = id;
    route = 'secret-preview';
    history.replaceState(null, '', `#/s/${id}`);
  }

  async function scrollCreate() {
    goTo('home');
    await tick();
    document.getElementById('createPanel')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  function onClearSession() {
    if (!confirm('Reset this browser session and clear locally cached links?')) return;
    resetSessionId();
    clearSecrets();
    notify('Browser session reset.');
  }

  onMount(() => {
    initTheme();
    parseHash();
    window.addEventListener('hashchange', parseHash);
    return () => window.removeEventListener('hashchange', parseHash);
  });
</script>

<Layout route={route === 'secret-preview' ? 'home' : route} onNavigate={goTo} onClearSession={onClearSession} onScrollCreate={scrollCreate}>
  {#if route === 'home'}
    <Home {openSecretPreview} />
  {:else if route === 'contacts'}
    <Contacts />
  {:else if route === 'news'}
    <News {goTo} />
  {:else if route === 'news-example'}
    <NewsArticle {goTo} />
  {:else if route === 'secret-preview'}
    <RecipientPreview {secretId} />
  {:else if route === 'destroy'}
    <Destroy />
  {:else if route === 'donate'}
    <Donate />
  {:else if route === 'api'}
    <Api />
  {:else if route === 'admin'}
    <Admin />
  <!-- {:else if route === 'ui'}
    <UiKit /> -->
  {:else if route === 'about'}
    <About />
  {:else}
    <Home {openSecretPreview} />
  {/if}
</Layout>

<Toast />
