# BurnVault Svelte API port

This build keeps the approved BurnVault visual design and connects the Svelte UI to an external HTTP API.

The CSS is still the fixed design source:

```text
src/app.css
```

No scoped component CSS is used. Page components only contain Svelte state, events, API calls, and markup.

## Configure API

Use an absolute API URL:

```bash
VITE_API_BASE_URL=https://api.example.com npm run dev
```

Or leave it empty and proxy `/v1` from the same origin:

```bash
npm run dev
```

In same-origin mode, requests go to paths such as:

```text
/v1/secrets
/v1/session/secrets
/v1/admin/secrets
```

## Run

```bash
npm install
npm run dev
```

## Build

```bash
npm run build
npm run preview
```

## API code location

```text
src/lib/services/api.js       low-level fetch wrapper
src/lib/stores/secrets.js     session links, create, destroy
src/lib/stores/stats.js       sidebar public stats
src/lib/pages/Home.svelte     create link and recent links
src/lib/pages/RecipientPreview.svelte  metadata + reveal flow
src/lib/pages/Admin.svelte    paginated server-side admin search
```

## Backend contract

See:

```text
docs/backend-contract.md
```
