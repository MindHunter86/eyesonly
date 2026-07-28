<script>
  import { onMount } from 'svelte';
  import SelectControl from '../components/SelectControl.svelte';
  import SecretTable from '../components/SecretTable.svelte';
  import { buildSecretUrl, createSecret, getStatus, loadSessionSecrets, secrets } from '../stores/secrets.js';
  import { notify } from '../stores/toast.js';

  export let openSecretPreview;

  let secretText = '';
  let ttl = '3600';
  let maxViews = 1;
  let hint = '';
  let burnAfterRead = true;
  let hideFromHistory = true;
  let notifyOnOpen = false;
  let notifyEmail = '';
  let lastCreated = null;
  let showDestroyToken = false;
  let isSubmitting = false;

  const ttlOptions = [
    { value: '300', label: '5 minutes' },
    { value: '3600', label: '1 hour' },
    { value: '86400', label: '24 hours' },
    { value: '604800', label: '7 days' },
  ];

  $: counts = $secrets.reduce((acc, secret) => {
    acc[getStatus(secret)] += 1;
    return acc;
  }, { active: 0, viewed: 0, burned: 0, expired: 0 });

  async function onCreateSecret() {
    const text = secretText.trim();
    if (!text) {
      notify('Secret text is required.', 'failure');
      return;
    }

    if (notifyOnOpen && !notifyEmail.trim()) {
      notify('Email is required when open notification is enabled.', 'failure');
      return;
    }

    isSubmitting = true;

    try {
      const created = await createSecret({
        secret: text,
        ttl,
        maxViews,
        hint,
        burnAfterRead,
        hideFromHistory,
        notifyOnOpen,
        notifyEmail,
      });

      lastCreated = created;
      showDestroyToken = false;
      hint = '';
      if (hideFromHistory) secretText = '';
      if (!notifyOnOpen) notifyEmail = '';
      notify('Disposable link created.');
    } catch (error) {
      notify(error.message || 'Failed to create secret.', 'failure');
    } finally {
      isSubmitting = false;
    }
  }

  onMount(() => {
    loadSessionSecrets().catch(error => {
      notify(error.message || 'Failed to load session links.', 'failure');
    });
  });

  function copyText(value) {
    navigator.clipboard.writeText(value).then(
      () => notify('Copied to clipboard.'),
      () => notify('Clipboard access was blocked.', 'failure')
    );
  }
</script>

<section class="view active" id="view-home" aria-labelledby="home-title">
  <div class="metric-grid" aria-label="Secret statistics">
    <article class="metric-card">
      <span class="metric-label">Active</span>
      <strong class="metric-value">{counts.active}</strong>
    </article>
    <article class="metric-card">
      <span class="metric-label">Viewed</span>
      <strong class="metric-value">{counts.viewed}</strong>
    </article>
    <article class="metric-card">
      <span class="metric-label">Burned</span>
      <strong class="metric-value">{counts.burned}</strong>
    </article>
    <article class="metric-card">
      <span class="metric-label">Expired</span>
      <strong class="metric-value">{counts.expired}</strong>
    </article>
  </div>

  <section class="panel create-panel" id="createPanel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">create</p>
        <h2 id="home-title">Create a disposable link</h2>
      </div>
      <span class="status-pill success">api mode</span>
    </div>

    <form class="secret-form" on:submit|preventDefault={onCreateSecret}>
      <label class="field full">
        <span>Secret</span>
        <textarea
          bind:value={secretText}
          name="secret"
          rows="7"
          placeholder="Paste password, token, note, or any short secret..."
          required
        ></textarea>
      </label>

      <div class="form-grid">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <label class="field">
          <span>TTL</span>
          <SelectControl id="ttlSelect" name="ttl" bind:value={ttl} options={ttlOptions} />
        </label>

        <label class="field">
          <span>Max views</span>
          <input bind:value={maxViews} name="maxViews" type="number" min="1" max="25" />
        </label>

        <label class="field">
          <span>Secret public description</span>
          <input
            bind:value={hint}
            name="hint"
            type="text"
            maxlength="64"
            placeholder="optional public description shown before reveal"
          />
        </label>
      </div>

      <div class="option-row">
        <div class="switch-line">
          <input id="burnAfterRead" type="checkbox" bind:checked={burnAfterRead} />
          <label for="burnAfterRead">Burn after first view</label>
        </div>
        <div class="switch-line">
          <input id="hideFromHistory" type="checkbox" bind:checked={hideFromHistory} />
          <label for="hideFromHistory">Do not store secret text in browser state</label>
        </div>
      </div>

      <div class="notify-row">
        <div class="switch-line">
          <input id="notifyOnOpen" type="checkbox" bind:checked={notifyOnOpen} />
          <label for="notifyOnOpen">Notify me when this secret is opened</label>
        </div>
        {#if notifyOnOpen}
          <label class="field notify-email">
            <span>Email for open notification</span>
            <input bind:value={notifyEmail} type="email" autocomplete="email" placeholder="name@example.com" />
          </label>
        {/if}
      </div>

      <div class="form-actions">
        <p class="muted">
          Secret body is submitted to the configured API. Backend must encrypt, store, and burn it server-side.
        </p>
        <button class="primary-button" type="submit" disabled={isSubmitting}>{isSubmitting ? 'Creating...' : 'Create link'}</button>
      </div>
    </form>
  </section>

  {#if lastCreated}
    <section class="panel result-panel">
      <div class="panel-heading">
        <div>
          <p class="eyebrow">created</p>
          <h2>Share this link once</h2>
        </div>
        <span class="status-pill">ready</span>
      </div>
      <div class="result-grid">
        <label class="field">
          <span>Secret link</span>
          <div class="copy-line">
            <input value={buildSecretUrl(lastCreated.id, lastCreated.url)} type="text" readonly />
            <button class="ghost-button" type="button" on:click={() => copyText(buildSecretUrl(lastCreated.id, lastCreated.url))}>Copy</button>
          </div>
        </label>
        <label class="field">
          <span>Destroy token</span>
          <div class="copy-line">
            <input value={lastCreated.destroyToken} type={showDestroyToken ? 'text' : 'password'} readonly />
            <button class="ghost-button" type="button" on:click={() => showDestroyToken = !showDestroyToken}>{showDestroyToken ? 'Hide' : 'Show'}</button>
            <button class="ghost-button" type="button" on:click={() => copyText(lastCreated.destroyToken)}>Copy</button>
          </div>
        </label>
        <div class="result-actions full">
          <button class="ghost-button" type="button" on:click={() => openSecretPreview(lastCreated.id)}>Preview recipient page</button>
        </div>
        <p class="muted full">
          Send the destroy token only to a trusted person. It can revoke the link before it is viewed.
        </p>
      </div>
    </section>
  {/if}

  <SecretTable {openSecretPreview} />
</section>
