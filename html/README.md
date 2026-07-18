# BurnVault frontend MVP

Static frontend prototype for one-time secret links.

This build keeps the approved `burnvault-hyperdx-terminal-checkbox-fix.zip` dark design as the baseline and adds only local UI/functionality patches.

## Included

- Home screen with disposable link creation
- Destroy token flow
- Hidden notification email field enabled by checkbox
- Contacts, News, full news article example, Donate, API, Admin render demo, UI kit, About
- Inline API Try response preview blocks
- Success/failure notify states
- Demo admin list with search, page size selector, pagination and session hashes
- Mobile multiline table rendering for Recent secrets/Admin tables

## Run locally

```bash
python3 -m http.server 4173
```

Open:

```text
http://localhost:4173
```

## Notes

This is a static demo. Backend encryption, real delivery, real notifications, API requests, and admin search are placeholders for integration.


## Latest point patch

- Added a recipient preview page that keeps the secret hidden until the user confirms reveal.
- Renamed Passphrase hint to Secret public description and removed the explanatory note under the input.
- Added admin session-id click filtering.
- Added a clear button inside the admin search field.

