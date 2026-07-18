<script>
  import SelectControl from '../components/SelectControl.svelte';

  let uiSelect = '1 hour';
</script>

<section class="view page-stack active" id="view-ui">
  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">ui / temporary</p>
        <h2>Project UI kit</h2>
      </div>
      <span class="status-pill">builder</span>
    </div>
    <div class="ui-kit-grid">
      <section class="ui-demo-block">
        <p class="eyebrow">buttons</p>
        <div class="notify-test-row">
          <button class="primary-button" type="button">Primary</button>
          <button class="ghost-button" type="button">Ghost</button>
          <button class="filter-chip active" type="button">Filter</button>
          <button class="primary-button danger-button" type="button">Danger</button>
        </div>
      </section>
      <section class="ui-demo-block">
        <p class="eyebrow">statuses</p>
        <div class="notify-test-row">
          <span class="status-pill success">active</span>
          <span class="status-pill">pending</span>
          <span class="status-pill danger">burned</span>
        </div>
      </section>
      <section class="ui-demo-block">
        <p class="eyebrow">fields</p>
        <div class="form-grid ui-form-demo">
          <label class="field">
            <span>Input</span>
            <input type="text" placeholder="secret id" />
          </label>
          <!-- svelte-ignore a11y-label-has-associated-control -->
          <label class="field">
            <span>Select</span>
            <SelectControl bind:value={uiSelect} options={[{ value: '1 hour', label: '1 hour' }, { value: '24 hours', label: '24 hours' }]} />
          </label>
          <label class="field">
            <span>Checkbox</span>
            <span class="switch-line">
              <input id="uiDemoCheckbox" type="checkbox" checked />
              <label for="uiDemoCheckbox">Clickable label</label>
            </span>
          </label>
        </div>
      </section>
      <section class="ui-demo-block">
        <p class="eyebrow">text</p>
        <div class="prose">
          <h3>Section heading</h3>
          <p>Muted body text, inline <code>code</code>, and normal paragraph spacing for content pages.</p>
          <blockquote>Important project notes can use the same accent surface.</blockquote>
        </div>
      </section>
      <section class="ui-demo-block full">
        <p class="eyebrow">table</p>
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>ID</th><th>Status</th><th>TTL</th><th class="actions-cell">Actions</th></tr>
            </thead>
            <tbody>
              <tr>
                <td data-label="ID"><span class="id-code">demo01</span></td>
                <td data-label="Status"><span class="status-pill success">active</span></td>
                <td data-label="TTL">42m</td>
                <td class="actions-cell" data-label="Actions"><div class="row-actions"><button class="ghost-button" type="button">Copy</button></div></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </section>
</section>
