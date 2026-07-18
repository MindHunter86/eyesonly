<script>
  import { secrets, secretsError, secretsLoading, buildSecretUrl, destroyById, formatTtl, getStatus, statusClass } from '../stores/secrets.js';
  import { notify } from '../stores/toast.js';

  export let openSecretPreview;

  let filter = 'all';

  $: visibleSecrets = $secrets.filter(secret => filter === 'all' || getStatus(secret) === filter);

  function copyLink(secret) {
    navigator.clipboard.writeText(buildSecretUrl(secret.id, secret.url)).then(
      () => notify('Copied to clipboard.'),
      () => notify('Clipboard access was blocked.', 'failure')
    );
  }

  async function onDestroy(id) {
    if (!confirm(`Destroy secret ${id}?`)) return;

    try {
      await destroyById(id);
      notify('Secret destroyed.');
    } catch (error) {
      notify(error.message || 'Failed to destroy secret.', 'failure');
    }
  }
</script>

<section class="panel">
  <div class="panel-heading table-heading">
    <div>
      <p class="eyebrow">session</p>
      <h2>Recent secrets</h2>
    </div>
    <div class="filter-group" role="tablist" aria-label="Status filter">
      {#each ['all', 'active', 'burned'] as option}
        <button
          class="filter-chip"
          class:active={filter === option}
          type="button"
          on:click={() => filter = option}
        >{option[0].toUpperCase() + option.slice(1)}</button>
      {/each}
    </div>
  </div>

  <div class="table-wrap recent-table-wrap">
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Status</th>
          <th>TTL</th>
          <th>Views</th>
          <th>Public description</th>
          <th class="actions-cell">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if $secretsLoading}
          <tr class="empty-row">
            <td colspan="6">Loading session links from API...</td>
          </tr>
        {:else if $secretsError}
          <tr class="empty-row">
            <td colspan="6">{$secretsError}</td>
          </tr>
        {:else if visibleSecrets.length === 0}
          <tr class="empty-row">
            <td colspan="6">No secrets match the current filter.</td>
          </tr>
        {:else}
          {#each visibleSecrets as secret (secret.id)}
            {@const status = getStatus(secret)}
            {@const canUse = status === 'active' || status === 'viewed'}
            <tr>
              <td data-label="ID"><span class="id-code">{secret.id}</span></td>
              <td data-label="Status"><span class={`status-pill ${statusClass(status)}`}>{status}</span></td>
              <td data-label="TTL">{formatTtl(secret)}</td>
              <td data-label="Views">{secret.views}/{secret.maxViews}</td>
              <td data-label="Public description">{secret.hint || '-'}</td>
              <td class="actions-cell" data-label="Actions">
                <div class="row-actions">
                  <button class="ghost-button" type="button" on:click={() => copyLink(secret)}>Copy</button>
                  <button class="ghost-button" type="button" disabled={!canUse} on:click={() => openSecretPreview(secret.id)}>Open</button>
                  <button class="ghost-button" type="button" disabled={!canUse} on:click={() => onDestroy(secret.id)}>Destroy</button>
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</section>
