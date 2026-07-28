<script>
  export let goTo;

  const items = [
    {
      date: '2026-07-08',
      title: 'Destroy tokens added',
      text: 'Newly created demo links now include a separate destroy token. A trusted receiver can revoke the secret without seeing its contents.',
    },
    {
      date: '2026-07-01',
      title: 'Theme preview',
      text: 'The interface now supports a dark terminal palette and a light theme tuned from the same green accent color.',
    },
    {
      date: '2026-06-24',
      title: 'Session metadata table',
      text: 'Local session metadata shows status, TTL, view count, and available actions for each link.',
    },
  ];
</script>

<section class="view page-stack active" id="view-news">
  <section class="panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">updates</p>
        <h2>Product news</h2>
      </div>
      <span class="status-pill">changelog</span>
    </div>
    <div class="timeline">
      {#each items as item}
        <article class="timeline-item news-card">
          <a class="timeline-link" href="#news-example" on:click|preventDefault={() => goTo('news-example')}>
            <time>{item.date}</time>
            <h3>{item.title}</h3>
            <p>{item.text}</p>
          </a>
        </article>
      {/each}
    </div>
  </section>
</section>
