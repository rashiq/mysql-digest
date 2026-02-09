const elements = {
  sql: document.getElementById("sql"),
  version: document.getElementById("version"),
  output: document.getElementById("output"),
  text: document.getElementById("text"),
  hash: document.getElementById("hash"),
  error: document.getElementById("error"),
  errorMessage: document.getElementById("error-message"),
  loading: document.getElementById("loading"),
};

let ready = false;

function loadFromURL() {
  const params = new URLSearchParams(window.location.search);
  const query = params.get("query");
  const version = params.get("version");
  if (query !== null) {
    elements.sql.value = query;
  }
  if (version !== null && ["0", "1", "2", "3"].includes(version)) {
    elements.version.value = version;
  }
}

function updateURL() {
  const sql = elements.sql.value.trim();
  const version = elements.version.value;
  const params = new URLSearchParams();
  if (sql) {
    params.set("query", sql);
  }
  if (version !== "2") {
    params.set("version", version);
  }
  const qs = params.toString();
  const url = window.location.pathname + (qs ? "?" + qs : "");
  history.replaceState(null, "", url);
}

async function init() {
  loadFromURL();
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch("digest.wasm"),
    go.importObject
  );
  go.run(result.instance);
  ready = true;
  elements.loading.classList.add("hidden");
  compute();
}

function compute() {
  if (!ready) return;

  updateURL();

  const sql = elements.sql.value.trim();
  if (!sql) {
    elements.output.classList.add("hidden");
    elements.error.classList.add("hidden");
    return;
  }

  const version = parseInt(elements.version.value, 10);
  const resultJson = computeDigest(sql, version);
  const result = JSON.parse(resultJson);

  if (result.error) {
    elements.output.classList.add("hidden");
    elements.error.classList.remove("hidden");
    elements.errorMessage.textContent = result.error;
    return;
  }

  elements.error.classList.add("hidden");
  elements.output.classList.remove("hidden");
  elements.text.textContent = result.text;
  elements.hash.textContent = result.hash;
}

function copyLink() {
  navigator.clipboard.writeText(window.location.href);
}

function copyText(id) {
  const text = document.getElementById(id).textContent;
  navigator.clipboard.writeText(text);
}

let debounceTimer;
function debounce(fn, delay) {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fn, delay);
}

elements.sql.addEventListener("input", () => debounce(compute, 150));
elements.version.addEventListener("change", compute);

init();
