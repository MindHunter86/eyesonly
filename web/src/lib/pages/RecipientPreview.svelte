<script>
  import { onMount } from 'svelte';
  import { backendApi } from '../services/api.js';
  import { normalizeSecret, getStatus, statusClass } from '../stores/secrets.js';
  import { notify } from '../stores/toast.js';

  export let secretId = null;

  let loadedId = null;
  let secret = null;
  let secretValue = '';
  let revealed = false;
  let loading = false;
  let revealing = false;
  let error = '';
  let mounted = false;

  $: status = secret ? getStatus(secret) : 'active';
  $: displayId = secret?.id || secretId || 'unknown';
  $: displayDescription = secret?.hint || 'No public description provided';
  $: visibleStatus = revealed ? 'revealed' : status === 'active' ? 'ready' : status;
  $: pillClass = revealed ? 'danger' : status === 'active' ? 'success' : statusClass(status);

  async function loadMetadata(id) {
    if (!id) {
      error = 'Secret id is missing.';
      return;
    }

    loading = true;
    error = '';
    revealed = false;
    secretValue = '';

    try {
      const data = await backendApi.getSecret(id);
      secret = normalizeSecret(data);
      loadedId = id;
    } catch (loadError) {
      secret = null;
      error = loadError.message || 'Failed to load secret metadata.';
    } finally {
      loading = false;
    }
  }

  async function reveal() {
    if (!secretId || revealing) return;

    revealing = true;
    error = '';

    try {
      const data = await backendApi.revealSecret(secretId);
      secretValue = data.secret || '';
      revealed = true;
      secret = normalizeSecret({ ...secret, ...data, id: secretId, status: data.burned ? 'burned' : data.status });
      notify('Secret revealed once.');
    } catch (revealError) {
      error = revealError.message || 'Failed to reveal secret.';
      notify(error, 'failure');
    } finally {
      revealing = false;
    }
  }

  onMount(() => {
    mounted = true;
    loadMetadata(secretId);
  });

  $: if (mounted && secretId && secretId !== loadedId && !loading) {
    loadMetadata(secretId);
  }
</script>

<section class="view page-stack active" id="view-secret-preview">
  <section class="panel create-panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">recipient / one-time view</p>
        <h2>Open shared secret</h2>
      </div>
      <span class={`status-pill ${pillClass}`}>{visibleStatus}</span>
    </div>
    <div class="recipient-demo-grid">
      <div class="prose">
        <p>
          This is the page a recipient sees after opening a one-time secret link.
          The secret body stays hidden until the recipient explicitly confirms reveal.
        </p>
        <p class="muted">
          Secret ID: <span class="id-code">{displayId}</span>
        </p>
        <p class="muted">
          Description: <span>{displayDescription}</span>
        </p>
      </div>

      {#if loading}
        <div class="recipient-confirm-box">
          <h3>Loading</h3>
          <p class="muted">Loading secret metadata from API...</p>
        </div>
      {:else if error}
        <div class="recipient-confirm-box">
          <h3>Unavailable</h3>
          <p class="muted">{error}</p>
        </div>
      {:else if !revealed}
        <div class="recipient-confirm-box">
          <h3>Reveal once</h3>
          <p class="muted">
            After confirmation, the backend will mark the link as viewed and burn it according to the policy.
          </p>
          <button class="primary-button" type="button" disabled={revealing} on:click={reveal}>{revealing ? 'Revealing...' : 'Reveal secret'}</button>
        </div>
      {:else}
        <div class="recipient-secret-box">
          <span class="api-label">secret content</span>
          <pre><code>{secretValue}</code></pre>
          <p class="muted">The secret was returned by the API and should not be available again if burn policy was applied.</p>
        </div>
      {/if}
    </div>
  </section>
</section>
