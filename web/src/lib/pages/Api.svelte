<script>
  import { backendApi } from '../services/api.js';
  import { notify } from '../stores/toast.js';

  let openResult = null;
  let createSecretText = 'BurnVault API test secret';
  let readId = '';
  let revealId = '';
  let destroyToken = '';
  let resultText = '';
  let loadingAction = '';

  function formatJson(value) {
    return JSON.stringify(value, null, 2);
  }

  async function runTry(action, request) {
    loadingAction = action;
    openResult = action;
    resultText = '';

    try {
      const data = await request();
      resultText = formatJson({ ok: true, data });
      notify(`API response rendered for ${action}.`);
    } catch (error) {
      resultText = formatJson({
        ok: false,
        error: {
          code: error.code || 'REQUEST_FAILED',
          message: error.message || 'API request failed.',
          status: error.status || 0,
        },
      });
      notify(error.message || 'API request failed.', 'failure');
    } finally {
      loadingAction = '';
    }
  }

  function tryCreate() {
    return runTry('create', () => backendApi.createSecret({
      secret: createSecretText,
      ttl_seconds: 3600,
      max_views: 1,
      burn_after_read: true,
      public_description: 'API page test',
    }));
  }

  function tryRead() {
    return runTry('read', () => backendApi.getSecret(readId.trim()));
  }

  function tryReveal() {
    return runTry('reveal', () => backendApi.revealSecret(revealId.trim()));
  }

  function tryDestroy() {
    return runTry('destroy', () => backendApi.destroySecret(destroyToken.trim()));
  }

  function tryStats() {
    return runTry('stats', () => backendApi.getStats());
  }
</script>

<section class="view page-stack active" id="view-api">
  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">api</p>
        <h2>HTTP API preview</h2>
      </div>
      <span class="status-pill">live api</span>
    </div>

    <div class="prose api-intro">
      <p>
        The examples below call the configured backend API. Set <code>VITE_API_BASE_URL</code>
        or proxy <code>/v1</code> from the same origin during development.
      </p>
    </div>

    <div class="api-list" aria-label="API method examples">
      <article class="api-method">
        <div class="api-method-head">
          <span class="method-badge method-post">POST</span>
          <code>/v1/secrets</code>
        </div>
        <p>Create an encrypted one-time secret and return a public link plus a destroy token.</p>
        <div class="api-example-grid">
          <div>
            <span class="api-label">usage</span>
            <pre><code>{`{
  "secret": "DATABASE_URL=...",
  "ttl_seconds": 3600,
  "max_views": 1,
  "burn_after_read": true,
  "public_description": "deployment handoff"
}`}</code></pre>
          </div>
          <div>
            <span class="api-label">response</span>
            <pre><code>{`{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "url": "https://example.com/#/s/9f31ab20",
    "destroy_token": "dt_..."
  }
}`}</code></pre>
          </div>
        </div>
        <div class="api-try-line">
          <input bind:value={createSecretText} type="text" placeholder="secret text for test request" />
          <button class="ghost-button" type="button" disabled={loadingAction === 'create'} on:click={tryCreate}>Try</button>
        </div>
        {#if openResult === 'create'}
          <div class="api-response-output">
            <span class="api-label">try response</span>
            <pre><code>{resultText}</code></pre>
          </div>
        {/if}
      </article>

      <article class="api-method">
        <div class="api-method-head">
          <span class="method-badge method-get">GET</span>
          <code>/v1/secrets/{'{id}'}</code>
        </div>
        <p>Return public metadata before reveal. The secret body must not be exposed here.</p>
        <div class="api-example-grid">
          <div>
            <span class="api-label">usage</span>
            <pre><code>GET /v1/secrets/9f31ab20</code></pre>
          </div>
          <div>
            <span class="api-label">response</span>
            <pre><code>{`{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "status": "active",
    "public_description": "deployment handoff"
  }
}`}</code></pre>
          </div>
        </div>
        <div class="api-try-line">
          <input bind:value={readId} type="text" placeholder="secret id" />
          <button class="ghost-button" type="button" disabled={!readId.trim() || loadingAction === 'read'} on:click={tryRead}>Try</button>
        </div>
        {#if openResult === 'read'}
          <div class="api-response-output">
            <span class="api-label">try response</span>
            <pre><code>{resultText}</code></pre>
          </div>
        {/if}
      </article>

      <article class="api-method">
        <div class="api-method-head">
          <span class="method-badge method-post">POST</span>
          <code>/v1/secrets/{'{id}'}/reveal</code>
        </div>
        <p>Reveal the secret only after explicit confirmation. Backend should burn or update view counters.</p>
        <div class="api-example-grid">
          <div>
            <span class="api-label">usage</span>
            <pre><code>POST /v1/secrets/9f31ab20/reveal</code></pre>
          </div>
          <div>
            <span class="api-label">response</span>
            <pre><code>{`{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "secret": "visible-once",
    "burned": true
  }
}`}</code></pre>
          </div>
        </div>
        <div class="api-try-line">
          <input bind:value={revealId} type="text" placeholder="secret id" />
          <button class="ghost-button" type="button" disabled={!revealId.trim() || loadingAction === 'reveal'} on:click={tryReveal}>Try</button>
        </div>
        {#if openResult === 'reveal'}
          <div class="api-response-output">
            <span class="api-label">try response</span>
            <pre><code>{resultText}</code></pre>
          </div>
        {/if}
      </article>

      <article class="api-method">
        <div class="api-method-head">
          <span class="method-badge method-post">POST</span>
          <code>/v1/destroy</code>
        </div>
        <p>Destroy a secret before it is opened by using the destroy token.</p>
        <div class="api-example-grid">
          <div>
            <span class="api-label">usage</span>
            <pre><code>{`{
  "destroy_token": "dt_..."
}`}</code></pre>
          </div>
          <div>
            <span class="api-label">response</span>
            <pre><code>{`{
  "ok": true,
  "data": {
    "id": "9f31ab20",
    "status": "burned"
  }
}`}</code></pre>
          </div>
        </div>
        <div class="api-try-line">
          <input bind:value={destroyToken} type="text" placeholder="destroy token" />
          <button class="ghost-button" type="button" disabled={!destroyToken.trim() || loadingAction === 'destroy'} on:click={tryDestroy}>Try</button>
        </div>
        {#if openResult === 'destroy'}
          <div class="api-response-output">
            <span class="api-label">try response</span>
            <pre><code>{resultText}</code></pre>
          </div>
        {/if}
      </article>

      <article class="api-method">
        <div class="api-method-head">
          <span class="method-badge method-get">GET</span>
          <code>/v1/stats</code>
        </div>
        <p>Return aggregate public counters for landing-page statistics. No secret body is exposed.</p>
        <div class="api-example-grid">
          <div>
            <span class="api-label">usage</span>
            <pre><code>GET /v1/stats</code></pre>
          </div>
          <div>
            <span class="api-label">response</span>
            <pre><code>{`{
  "ok": true,
  "data": {
    "created": 12840,
    "burned": 12611,
    "avg_ttl": "42m"
  }
}`}</code></pre>
          </div>
        </div>
        <div class="api-try-line">
          <input type="text" placeholder="no input required" disabled />
          <button class="ghost-button" type="button" disabled={loadingAction === 'stats'} on:click={tryStats}>Try</button>
        </div>
        {#if openResult === 'stats'}
          <div class="api-response-output">
            <span class="api-label">try response</span>
            <pre><code>{resultText}</code></pre>
          </div>
        {/if}
      </article>
    </div>
  </section>

  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">notify</p>
        <h2>Notification states</h2>
      </div>
      <span class="status-pill">ui test</span>
    </div>
    <div class="notify-test-row">
      <button class="ghost-button" type="button" on:click={() => notify('Success notification rendered.')}>Show success</button>
      <button class="ghost-button" type="button" on:click={() => notify('Failure notification rendered.', 'failure')}>Show failure</button>
    </div>
  </section>
</section>
