<script>
  import { onMount } from 'svelte';
  import SelectControl from '../components/SelectControl.svelte';
  import { backendApi } from '../services/api.js';
  import { notify } from '../stores/toast.js';

  let query = '';
  let page = 1;
  let pageSize = '10';
  let rows = [];
  let total = 0;
  let loading = false;
  let error = '';
  let searchTimer = null;

  const pageSizeOptions = [
    { value: '5', label: '5' },
    { value: '10', label: '10' },
    { value: '25', label: '25' },
  ];

  $: totalPages = Math.max(1, Math.ceil(total / Number(pageSize)));

  function normalizeAdminRow(row) {
    return {
      id: row.id,
      session: row.session ?? row.session_id ?? '-',
      status: row.status || 'active',
      ttl: row.ttl || row.ttl_remaining || row.expires_in || 'server',
      views: row.views || `${row.views_used ?? 0}/${row.max_views ?? 1}`,
      preview: row.preview || row.public_description || row.content_preview || '-',
    };
  }

  async function loadRows() {
    loading = true;
    error = '';

    try {
      const data = await backendApi.listAdminSecrets({ q: query.trim(), page, pageSize: Number(pageSize) });
      rows = (data.items || []).map(normalizeAdminRow);
      total = data.total ?? rows.length;
    } catch (loadError) {
      rows = [];
      total = 0;
      error = loadError.message || 'Failed to load admin list.';
      notify(error, 'failure');
    } finally {
      loading = false;
    }
  }

  function queueSearch() {
    page = 1;
    clearTimeout(searchTimer);
    searchTimer = setTimeout(loadRows, 250);
  }

  function filterBySession(session) {
    query = session;
    page = 1;
    loadRows();
  }

  function clearSearch() {
    query = '';
    page = 1;
    loadRows();
  }

  function changePage(nextPage) {
    page = nextPage;
    loadRows();
  }

  function changePageSize(value) {
    pageSize = value;
    page = 1;
    loadRows();
  }

  function copyId(id) {
    navigator.clipboard.writeText(id).then(
      () => notify('Copied to clipboard.'),
      () => notify('Clipboard access was blocked.', 'failure')
    );
  }

  onMount(loadRows);
</script>

<section class="view page-stack active" id="view-admin">
  <section class="panel">
    <div class="panel-heading table-heading">
      <div>
        <p class="eyebrow">admin</p>
        <h2>All secret links</h2>
      </div>
      <span class="status-pill">api list</span>
    </div>
    <div class="admin-toolbar">
      <label class="field admin-search-field">
        <span>Search by ID, session, or secret content</span>
        <span class="search-control">
          <input bind:value={query} type="search" placeholder="type to request filtered list..." on:input={queueSearch} />
          {#if query}
            <button class="search-clear" type="button" aria-label="Clear search" on:click={clearSearch}>×</button>
          {/if}
        </span>
      </label>
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <label class="field admin-page-size">
        <span>Rows</span>
        <SelectControl id="adminPageSize" value={pageSize} options={pageSizeOptions} onChange={changePageSize} />
      </label>
    </div>
    <div class="table-wrap admin-table-wrap">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Session</th>
            <th>Status</th>
            <th>TTL</th>
            <th>Views</th>
            <th>Preview</th>
            <th class="actions-cell">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr class="empty-row">
              <td colspan="7">Loading filtered list from API...</td>
            </tr>
          {:else if error}
            <tr class="empty-row">
              <td colspan="7">{error}</td>
            </tr>
          {:else if rows.length === 0}
            <tr class="empty-row">
              <td colspan="7">No links match the current search.</td>
            </tr>
          {:else}
            {#each rows as row (row.id)}
              <tr>
                <td data-label="ID"><span class="id-code">{row.id}</span></td>
                <td data-label="Session"><button class="id-code session-filter-button" type="button" title={`Filter by session ${row.session}`} on:click={() => filterBySession(row.session)}>{row.session}</button></td>
                <td data-label="Status"><span class={`status-pill ${row.status === 'active' ? 'success' : row.status === 'burned' ? 'danger' : ''}`}>{row.status}</span></td>
                <td data-label="TTL">{row.ttl}</td>
                <td data-label="Views">{row.views}</td>
                <td data-label="Preview">{row.preview}</td>
                <td class="actions-cell" data-label="Actions">
                  <div class="row-actions">
                    <button class="ghost-button" type="button" on:click={() => copyId(row.id)}>Copy ID</button>
                    <a class="ghost-button" href={`#/s/${row.id}`}>Open</a>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
    <div class="admin-pagination">
      <button class="ghost-button" type="button" disabled={page <= 1 || loading} on:click={() => changePage(page - 1)}>Prev</button>
      <span class="muted">Page {page} / {totalPages} · {total} matched</span>
      <button class="ghost-button" type="button" disabled={page >= totalPages || loading} on:click={() => changePage(page + 1)}>Next</button>
    </div>
  </section>
</section>
