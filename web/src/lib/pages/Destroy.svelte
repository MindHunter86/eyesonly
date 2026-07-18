<script>
  import { destroyByToken } from '../stores/secrets.js';
  import { notify } from '../stores/toast.js';

  let destroyToken = '';
  let isSubmitting = false;

  async function onDestroy() {
    const token = destroyToken.trim();
    if (!token || isSubmitting) return;

    isSubmitting = true;

    try {
      const result = await destroyByToken(token);
      destroyToken = '';
      notify(result.id ? `Secret ${result.id} destroyed.` : 'Secret destroyed.');
    } catch (error) {
      notify(error.message || 'Failed to destroy secret.', 'failure');
    } finally {
      isSubmitting = false;
    }
  }
</script>

<section class="view page-stack active" id="view-destroy">
  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">destroy</p>
        <h2>Destroy a secret by token</h2>
      </div>
      <span class="status-pill danger">revocation</span>
    </div>
    <form class="secret-form" on:submit|preventDefault={onDestroy}>
      <label class="field full">
        <span>Destroy token</span>
        <input
          bind:value={destroyToken}
          type="password"
          autocomplete="off"
          placeholder="Paste destroy token..."
          required
        />
      </label>
      <div class="form-actions">
        <p class="muted">
          The token is sent to the backend and can revoke a link before it is opened.
        </p>
        <button class="primary-button danger-button" type="submit" disabled={isSubmitting}>{isSubmitting ? 'Destroying...' : 'Destroy secret'}</button>
      </div>
    </form>
  </section>
</section>
