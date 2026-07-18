const STORAGE_KEY = "burnvault.session.secrets";
const THEME_KEY = "burnvault.theme";

const state = {
  secrets: loadSecrets(),
  filter: "all",
  adminPage: 1,
  adminPageSize: 10,
  adminQuery: "",
  recipientDemoId: null,
};

const adminDemoRows = [
  { id: "a81f23bd", session: "s8d31c", status: "active", ttl: "51m", views: "0/1", preview: "DATABASE_URL secret" },
  { id: "b7029fca", session: "s8d31c", status: "burned", ttl: "burned", views: "1/1", preview: "temporary root password" },
  { id: "c19a0e44", session: "k4aa90", status: "active", ttl: "4h 12m", views: "0/2", preview: "deploy token for staging" },
  { id: "d554cbb0", session: "m77f2e", status: "expired", ttl: "expired", views: "0/1", preview: "vpn recovery code" },
  { id: "e0aa93c1", session: "k4aa90", status: "viewed", ttl: "11m", views: "1/3", preview: "support one-time note" },
  { id: "f442bb19", session: "n02de1", status: "active", ttl: "22h", views: "0/1", preview: "github webhook secret" },
  { id: "g96d133a", session: "n02de1", status: "burned", ttl: "burned", views: "1/1", preview: "invoice portal token" },
  { id: "h10c7aa2", session: "z55cba", status: "active", ttl: "6d", views: "0/5", preview: "shared maintenance note" },
  { id: "i2b901ef", session: "z55cba", status: "expired", ttl: "expired", views: "0/1", preview: "old ssh bootstrap key" },
  { id: "j77ef0a8", session: "p19af3", status: "active", ttl: "2h 3m", views: "0/1", preview: "billing handoff token" },
  { id: "k081a2dd", session: "p19af3", status: "viewed", ttl: "8m", views: "1/2", preview: "read-once customer note" },
  { id: "l6c0ad91", session: "r90d18", status: "active", ttl: "31m", views: "0/1", preview: "database migration passphrase" },
];

const elements = {
  form: document.querySelector("#secretForm"),
  secretText: document.querySelector("#secretText"),
  ttlSelect: document.querySelector("#ttlSelect"),
  maxViews: document.querySelector("#maxViews"),
  hint: document.querySelector("#hint"),
  burnAfterRead: document.querySelector("#burnAfterRead"),
  hideFromHistory: document.querySelector("#hideFromHistory"),
  notifyOnOpen: document.querySelector("#notifyOnOpen"),
  notifyEmailField: document.querySelector("#notifyEmailField"),
  notifyEmail: document.querySelector("#notifyEmail"),
  table: document.querySelector("#secretsTable"),
  toast: document.querySelector("#toast"),
  resultPanel: document.querySelector("#resultPanel"),
  createdLink: document.querySelector("#createdLink"),
  destroyToken: document.querySelector("#destroyToken"),
  copyLastButton: document.querySelector("#copyLastButton"),
  toggleDestroyTokenButton: document.querySelector("#toggleDestroyTokenButton"),
  copyDestroyTokenButton: document.querySelector("#copyDestroyTokenButton"),
  previewRecipientButton: document.querySelector("#previewRecipientButton"),
  recipientPreviewId: document.querySelector("#recipientPreviewId"),
  recipientPreviewHint: document.querySelector("#recipientPreviewHint"),
  recipientPreviewHintWrap: document.querySelector("#recipientPreviewHintWrap"),
  recipientPreviewStatus: document.querySelector("#recipientPreviewStatus"),
  recipientPreviewLocked: document.querySelector("#recipientPreviewLocked"),
  recipientPreviewSecret: document.querySelector("#recipientPreviewSecret"),
  revealRecipientSecretButton: document.querySelector("#revealRecipientSecretButton"),
  clearSessionButton: document.querySelector("#clearSessionButton"),
  themeToggleButton: document.querySelector("#themeToggleButton"),
  destroyForm: document.querySelector("#destroyForm"),
  destroyTokenInput: document.querySelector("#destroyTokenInput"),
  adminSearch: document.querySelector("#adminSearch"),
  adminSearchClear: document.querySelector("#adminSearchClear"),
  adminPageSize: document.querySelector("#adminPageSize"),
  adminLinksTable: document.querySelector("#adminLinksTable"),
  adminPageInfo: document.querySelector("#adminPageInfo"),
  adminPrevPage: document.querySelector("#adminPrevPage"),
  adminNextPage: document.querySelector("#adminNextPage"),
  counts: {
    active: document.querySelector("#activeCount"),
    viewed: document.querySelector("#viewedCount"),
    burned: document.querySelector("#burnedCount"),
    expired: document.querySelector("#expiredCount"),
  },
};

function loadSecrets() {
  try {
    return JSON.parse(sessionStorage.getItem(STORAGE_KEY) || "[]");
  } catch {
    return [];
  }
}

function saveSecrets() {
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(state.secrets));
}

function now() {
  return Date.now();
}

function createId() {
  const bytes = new Uint8Array(5);
  crypto.getRandomValues(bytes);
  return Array.from(bytes, byte => byte.toString(16).padStart(2, "0")).join("");
}

function createToken() {
  const bytes = new Uint8Array(18);
  crypto.getRandomValues(bytes);
  return Array.from(bytes, byte => byte.toString(16).padStart(2, "0")).join("");
}

function buildSecretUrl(id) {
  return `${location.origin}${location.pathname}#/s/${id}`;
}

function getStatus(secret) {
  if (secret.destroyedAt) {
    return "burned";
  }

  if (now() > secret.expiresAt) {
    return "expired";
  }

  if (secret.views > 0) {
    return "viewed";
  }

  return "active";
}

function formatTtl(secret) {
  const status = getStatus(secret);
  if (status === "burned" || status === "expired") {
    return status;
  }

  const seconds = Math.max(0, Math.floor((secret.expiresAt - now()) / 1000));
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);

  if (hours > 0) {
    return `${hours}h ${minutes % 60}m`;
  }

  return `${minutes}m ${seconds % 60}s`;
}

function showToast(message, type = "success") {
  elements.toast.textContent = message;
  elements.toast.classList.toggle("failure", type === "failure");
  elements.toast.classList.toggle("success", type !== "failure");
  elements.toast.classList.add("show");
  window.clearTimeout(showToast.timer);
  showToast.timer = window.setTimeout(() => {
    elements.toast.classList.remove("show");
  }, 2200);
}

function updateCounts() {
  const counts = {
    active: 0,
    viewed: 0,
    burned: 0,
    expired: 0,
  };

  for (const secret of state.secrets) {
    counts[getStatus(secret)] += 1;
  }

  elements.counts.active.textContent = counts.active;
  elements.counts.viewed.textContent = counts.viewed;
  elements.counts.burned.textContent = counts.burned;
  elements.counts.expired.textContent = counts.expired;
}

function renderSecrets() {
  updateCounts();

  const filtered = state.secrets.filter(secret => {
    const status = getStatus(secret);
    return state.filter === "all" || status === state.filter;
  });

  if (filtered.length === 0) {
    elements.table.innerHTML = `
      <tr class="empty-row">
        <td colspan="6">No secrets match the current filter.</td>
      </tr>
    `;
    return;
  }

  elements.table.innerHTML = filtered
    .map(secret => {
      const status = getStatus(secret);
      const pillClass = status === "active" ? "success" : status === "burned" ? "danger" : "";
      const canUse = status === "active" || status === "viewed";

      return `
        <tr>
          <td data-label="ID"><span class="id-code">${secret.id}</span></td>
          <td data-label="Status"><span class="status-pill ${pillClass}">${status}</span></td>
          <td data-label="TTL">${formatTtl(secret)}</td>
          <td data-label="Views">${secret.views}/${secret.maxViews}</td>
          <td data-label="Public description">${secret.hint || "-"}</td>
          <td class="actions-cell" data-label="Actions">
            <div class="row-actions">
              <button class="ghost-button" type="button" data-copy="${secret.id}">Copy</button>
              <button class="ghost-button" type="button" data-view-secret="${secret.id}" ${canUse ? "" : "disabled"}>Simulate view</button>
              <button class="ghost-button" type="button" data-destroy="${secret.id}" ${canUse ? "" : "disabled"}>Destroy</button>
            </div>
          </td>
        </tr>
      `;
    })
    .join("");
}

function renderAdminLinks() {
  if (!elements.adminLinksTable) {
    return;
  }

  const query = state.adminQuery.trim().toLowerCase();
  const filtered = adminDemoRows.filter(row => {
    return !query || row.id.includes(query) || row.preview.toLowerCase().includes(query) || row.session.includes(query);
  });

  const pageSize = state.adminPageSize;
  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  state.adminPage = Math.min(Math.max(1, state.adminPage), totalPages);
  const start = (state.adminPage - 1) * pageSize;
  const pageRows = filtered.slice(start, start + pageSize);

  if (pageRows.length === 0) {
    elements.adminLinksTable.innerHTML = `
      <tr class="empty-row">
        <td colspan="7">No demo links match the current search.</td>
      </tr>
    `;
  } else {
    elements.adminLinksTable.innerHTML = pageRows.map(row => {
      const pillClass = row.status === "active" ? "success" : row.status === "burned" ? "danger" : "";
      return `
        <tr>
          <td data-label="ID"><span class="id-code">${row.id}</span></td>
          <td data-label="Session"><button class="id-code session-filter-button" type="button" data-admin-session="${row.session}" title="Filter by session ${row.session}">${row.session}</button></td>
          <td data-label="Status"><span class="status-pill ${pillClass}">${row.status}</span></td>
          <td data-label="TTL">${row.ttl}</td>
          <td data-label="Views">${row.views}</td>
          <td data-label="Preview">${row.preview}</td>
          <td class="actions-cell" data-label="Actions">
            <div class="row-actions">
              <button class="ghost-button" type="button" data-admin-copy="${row.id}">Copy ID</button>
              <button class="ghost-button" type="button" data-admin-open="${row.id}">Open</button>
            </div>
          </td>
        </tr>
      `;
    }).join("");
  }

  elements.adminPageInfo.textContent = `Page ${state.adminPage} / ${totalPages} · ${filtered.length} matched`;
  elements.adminPrevPage.disabled = state.adminPage <= 1;
  elements.adminNextPage.disabled = state.adminPage >= totalPages;
}

function getRecipientDemoSecret() {
  if (state.recipientDemoId) {
    return state.secrets.find(item => item.id === state.recipientDemoId) || null;
  }

  return state.secrets[0] || null;
}

function renderRecipientPreview(resetReveal = true) {
  if (!elements.recipientPreviewId) {
    return;
  }

  const secret = getRecipientDemoSecret();
  const id = secret?.id || "demo";
  const hint = secret?.hint || "No public description provided";
  const status = secret ? getStatus(secret) : "active";

  elements.recipientPreviewId.textContent = id;
  elements.recipientPreviewHint.textContent = hint;
  elements.recipientPreviewHintWrap.hidden = false;
  elements.recipientPreviewStatus.textContent = status === "active" ? "ready" : status;
  elements.recipientPreviewStatus.className = `status-pill ${status === "active" ? "success" : status === "burned" ? "danger" : ""}`;

  if (resetReveal) {
    elements.recipientPreviewLocked.hidden = false;
    elements.recipientPreviewSecret.hidden = true;
  }
}

function revealRecipientSecret() {
  const secret = getRecipientDemoSecret();

  if (secret) {
    const status = getStatus(secret);
    if (status !== "active" && status !== "viewed") {
      showToast("This secret is no longer available.", "failure");
      renderRecipientPreview(false);
      return;
    }

    secret.views += 1;
    if (secret.burnAfterRead || secret.views >= secret.maxViews) {
      secret.destroyedAt = now();
    }
    saveSecrets();
    renderSecrets();
  }

  elements.recipientPreviewLocked.hidden = true;
  elements.recipientPreviewSecret.hidden = false;
  if (elements.recipientPreviewStatus) {
    elements.recipientPreviewStatus.textContent = "revealed";
    elements.recipientPreviewStatus.className = "status-pill danger";
  }
  showToast("Secret revealed once.");
}

function createSecret(event) {
  event.preventDefault();

  const text = elements.secretText.value.trim();
  if (!text) {
    showToast("Secret text is required.", "failure");
    return;
  }

  if (elements.notifyOnOpen?.checked && !elements.notifyEmail.value.trim()) {
    showToast("Email is required when open notification is enabled.", "failure");
    elements.notifyEmail.focus();
    return;
  }

  const id = createId();
  const ttlSeconds = Number(elements.ttlSelect.value);
  const maxViews = Math.max(1, Number(elements.maxViews.value || 1));
  const secret = {
    id,
    destroyToken: createToken(),
    maxViews,
    views: 0,
    burnAfterRead: elements.burnAfterRead.checked,
    hint: elements.hint.value.trim(),
    notifyOnOpen: Boolean(elements.notifyOnOpen?.checked),
    notifyEmail: elements.notifyOnOpen?.checked ? elements.notifyEmail.value.trim() : "",
    createdAt: now(),
    expiresAt: now() + ttlSeconds * 1000,
    destroyedAt: null,
  };

  state.secrets.unshift(secret);
  saveSecrets();
  renderSecrets();

  elements.createdLink.value = buildSecretUrl(id);
  elements.destroyToken.value = secret.destroyToken;
  elements.destroyToken.type = "password";
  elements.toggleDestroyTokenButton.textContent = "Show";
  elements.resultPanel.hidden = false;

  if (elements.hideFromHistory.checked) {
    elements.secretText.value = "";
  }

  elements.hint.value = "";
  if (elements.notifyEmail && !elements.notifyOnOpen?.checked) {
    elements.notifyEmail.value = "";
  }
  showToast("Disposable link created.");
}

function copyText(value) {
  navigator.clipboard.writeText(value).then(
    () => showToast("Copied to clipboard."),
    () => showToast("Clipboard access was blocked.", "failure")
  );
}

function simulateView(id) {
  const secret = state.secrets.find(item => item.id === id);
  if (!secret) {
    return;
  }

  const status = getStatus(secret);
  if (status !== "active" && status !== "viewed") {
    showToast("This secret is no longer available.", "failure");
    return;
  }

  secret.views += 1;

  if (secret.burnAfterRead || secret.views >= secret.maxViews) {
    secret.destroyedAt = now();
  }

  saveSecrets();
  renderSecrets();
  showToast("View simulated.");
}

function destroySecret(id) {
  const secret = state.secrets.find(item => item.id === id);
  if (!secret) {
    return;
  }

  if (!confirm(`Destroy secret ${id}?`)) {
    return;
  }

  secret.destroyedAt = now();
  saveSecrets();
  renderSecrets();
  showToast("Secret destroyed.");
}

function destroyByToken(token) {
  const secret = state.secrets.find(item => item.destroyToken === token.trim());
  if (!secret) {
    showToast("Destroy token was not found in this demo session.", "failure");
    return;
  }

  if (getStatus(secret) === "burned") {
    showToast("This secret is already destroyed.", "failure");
    return;
  }

  secret.destroyedAt = now();
  saveSecrets();
  renderSecrets();
  elements.destroyTokenInput.value = "";
  showToast(`Secret ${secret.id} destroyed.`);
}

function clearSession() {
  if (state.secrets.length === 0) {
    showToast("Session is already empty.");
    return;
  }

  if (!confirm("Remove all local demo metadata from this browser session?")) {
    return;
  }

  state.secrets = [];
  saveSecrets();
  renderSecrets();
  elements.resultPanel.hidden = true;
  showToast("Session metadata cleared.");
}

function getSavedTheme() {
  const savedTheme = localStorage.getItem(THEME_KEY);
  if (savedTheme === "light" || savedTheme === "dark") {
    return savedTheme;
  }

  return "dark";
}

function applyTheme(theme) {
  document.documentElement.dataset.theme = theme;
  localStorage.setItem(THEME_KEY, theme);
  elements.themeToggleButton.textContent = theme === "dark" ? "☼" : "☾";
  elements.themeToggleButton.setAttribute("aria-label", theme === "dark" ? "Switch to light theme" : "Switch to dark theme");
  elements.themeToggleButton.title = theme === "dark" ? "Switch to light theme" : "Switch to dark theme";
}

function toggleTheme() {
  const currentTheme = document.documentElement.dataset.theme || getSavedTheme();
  applyTheme(currentTheme === "dark" ? "light" : "dark");
}

function setView(name) {
  if (!document.querySelector(`#view-${name}`)) {
    name = "home";
  }

  document.querySelectorAll(".view").forEach(view => {
    view.classList.toggle("active", view.id === `view-${name}`);
  });

  document.querySelectorAll(".nav-link").forEach(link => {
    link.classList.toggle("active", link.dataset.view === name);
  });

  if (name === "secret-preview") {
    renderRecipientPreview();
  }

  if (location.hash !== `#${name}` && !location.hash.startsWith("#/s/")) {
    history.replaceState(null, "", `#${name}`);
  }
}

function setInitialView() {
  const hash = location.hash.replace("#", "");
  if (hash.startsWith("/s/")) {
    state.recipientDemoId = hash.replace("/s/", "") || null;
    setView("secret-preview");
    return;
  }

  setView(hash || "home");
}

document.addEventListener("click", event => {
  const navLink = event.target.closest("[data-view]");
  if (navLink) {
    event.preventDefault();
    setView(navLink.dataset.view);
  }

  const scrollCreate = event.target.closest("[data-scroll-create]");
  if (scrollCreate) {
    setView("home");
    document.querySelector("#createPanel").scrollIntoView({ behavior: "smooth" });
  }

  const copyButton = event.target.closest("[data-copy]");
  if (copyButton) {
    copyText(buildSecretUrl(copyButton.dataset.copy));
  }

  const viewButton = event.target.closest("[data-view-secret]");
  if (viewButton) {
    simulateView(viewButton.dataset.viewSecret);
  }

  const destroyButton = event.target.closest("[data-destroy]");
  if (destroyButton) {
    destroySecret(destroyButton.dataset.destroy);
  }

  const previewRecipientButton = event.target.closest("#previewRecipientButton");
  if (previewRecipientButton) {
    state.recipientDemoId = state.secrets[0]?.id || null;
    renderRecipientPreview();
    setView("secret-preview");
  }

  const sessionFilterButton = event.target.closest("[data-admin-session]");
  if (sessionFilterButton) {
    const session = sessionFilterButton.dataset.adminSession;
    state.adminQuery = session;
    state.adminPage = 1;
    if (elements.adminSearch) {
      elements.adminSearch.value = session;
    }
    if (elements.adminSearchClear) {
      elements.adminSearchClear.hidden = false;
    }
    renderAdminLinks();
  }

  const apiTryButton = event.target.closest("[data-api-try]");
  if (apiTryButton) {
    const method = apiTryButton.closest(".api-method");
    const input = method?.querySelector(".api-try-line input");
    const output = method?.querySelector("[data-api-result]");
    const code = output?.querySelector("code");
    const value = input?.value.trim();
    if (code) {
      code.textContent = JSON.stringify({
        ok: true,
        demo: true,
        method: apiTryButton.dataset.apiTry,
        requestValue: value || null,
        note: "This is where the backend response will be rendered.",
      }, null, 2);
    }
    if (output) {
      output.hidden = false;
    }
    const suffix = value ? " with provided value" : " in demo mode";
    showToast(`API ${apiTryButton.dataset.apiTry} request prepared${suffix}.`);
  }

  const adminCopyButton = event.target.closest("[data-admin-copy]");
  if (adminCopyButton) {
    copyText(adminCopyButton.dataset.adminCopy);
  }

  const adminOpenButton = event.target.closest("[data-admin-open]");
  if (adminOpenButton) {
    showToast(`Admin open request prepared for ${adminOpenButton.dataset.adminOpen}.`);
  }

  const notifyTestButton = event.target.closest("[data-notify-test]");
  if (notifyTestButton) {
    if (notifyTestButton.dataset.notifyTest === "failure") {
      showToast("Backend returned 500. Please retry later.", "failure");
    } else {
      showToast("Backend accepted the request.");
    }
  }
});

document.querySelectorAll("[data-filter]").forEach(button => {
  button.addEventListener("click", () => {
    state.filter = button.dataset.filter;
    document.querySelectorAll("[data-filter]").forEach(item => {
      item.classList.toggle("active", item === button);
    });
    renderSecrets();
  });
});

elements.form.addEventListener("submit", createSecret);
elements.copyLastButton.addEventListener("click", () => {
  if (elements.createdLink.value) {
    copyText(elements.createdLink.value);
  }
});
elements.copyDestroyTokenButton.addEventListener("click", () => {
  if (elements.destroyToken.value) {
    copyText(elements.destroyToken.value);
  }
});
elements.toggleDestroyTokenButton.addEventListener("click", () => {
  const isHidden = elements.destroyToken.type === "password";
  elements.destroyToken.type = isHidden ? "text" : "password";
  elements.toggleDestroyTokenButton.textContent = isHidden ? "Hide" : "Show";
});
if (elements.revealRecipientSecretButton) {
  elements.revealRecipientSecretButton.addEventListener("click", revealRecipientSecret);
}
elements.clearSessionButton.addEventListener("click", clearSession);
elements.themeToggleButton.addEventListener("click", toggleTheme);
if (elements.notifyOnOpen) {
  elements.notifyOnOpen.addEventListener("change", () => {
    elements.notifyEmailField.hidden = !elements.notifyOnOpen.checked;
    if (elements.notifyOnOpen.checked) {
      elements.notifyEmail.focus();
    }
  });
}
elements.destroyForm.addEventListener("submit", event => {
  event.preventDefault();
  destroyByToken(elements.destroyTokenInput.value);
});

if (elements.adminSearch) {
  elements.adminSearch.addEventListener("input", () => {
    state.adminQuery = elements.adminSearch.value;
    state.adminPage = 1;
    if (elements.adminSearchClear) {
      elements.adminSearchClear.hidden = elements.adminSearch.value.length === 0;
    }
    renderAdminLinks();
  });
}

if (elements.adminSearchClear) {
  elements.adminSearchClear.addEventListener("click", () => {
    state.adminQuery = "";
    state.adminPage = 1;
    elements.adminSearch.value = "";
    elements.adminSearchClear.hidden = true;
    elements.adminSearch.focus();
    renderAdminLinks();
  });
}

if (elements.adminPageSize) {
  elements.adminPageSize.addEventListener("change", () => {
    state.adminPageSize = Number(elements.adminPageSize.value || 10);
    state.adminPage = 1;
    renderAdminLinks();
  });
}

if (elements.adminPrevPage) {
  elements.adminPrevPage.addEventListener("click", () => {
    state.adminPage -= 1;
    renderAdminLinks();
  });
}

if (elements.adminNextPage) {
  elements.adminNextPage.addEventListener("click", () => {
    state.adminPage += 1;
    renderAdminLinks();
  });
}

applyTheme(getSavedTheme());
window.addEventListener("hashchange", setInitialView);
window.setInterval(renderSecrets, 1000);
setInitialView();
renderSecrets();
renderAdminLinks();
