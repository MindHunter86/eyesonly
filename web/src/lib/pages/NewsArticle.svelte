<script>
  export let goTo;
</script>

<section class="view page-stack active" id="view-news-example">
  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow"><a class="context-link" href="#news" on:click|preventDefault={() => goTo('news')}>news</a> / <a class="context-link" href="#news-example" on:click|preventDefault={() => goTo('news-example')}>example</a></p>
        <h2>Destroy tokens added</h2>
      </div>
      <span class="status-pill">article</span>
    </div>
    <article class="prose article-body">
      <p class="article-meta">Published 2026-07-08 · Product update · 4 min read</p>
      <p class="article-lead">
        Destroy tokens give senders and trusted recipients a second control path: a link can be
        revoked without exposing the secret body and without waiting for TTL expiration.
      </p>
      <h3>Why this exists</h3>
      <p>
        One-time links are useful only while the sender still trusts the delivery channel. If a link
        is pasted into the wrong chat, copied into a ticket, or sent to the wrong person, the sender
        needs a quick revocation path.
      </p>
      <blockquote>
        The destroy token is not the secret. It is a separate revocation credential and must be
        delivered only to someone trusted to invalidate the link.
      </blockquote>
      <h3>Recommended flow</h3>
      <ol>
        <li>Create the secret with the shortest acceptable TTL.</li>
        <li>Send the public link through the normal channel.</li>
        <li>Keep the destroy token private or share it with a trusted operator.</li>
        <li>Destroy the link if the delivery context changes.</li>
      </ol>
      <h3>Example response</h3>
      <pre><code>{`{
  "id": "9f31ab20",
  "url": "/s/9f31ab20",
  "destroyToken": "dt_demo_7b9c...",
  "expiresIn": 3600
}`}</code></pre>
      <table class="article-table">
        <thead>
          <tr>
            <th>State</th>
            <th>Meaning</th>
            <th>User action</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>active</td>
            <td>The secret can still be opened.</td>
            <td>Share or destroy.</td>
          </tr>
          <tr>
            <td>burned</td>
            <td>The secret was viewed or manually destroyed.</td>
            <td>No action required.</td>
          </tr>
        </tbody>
      </table>
    </article>
  </section>
</section>
