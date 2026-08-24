const state = { sources: [], metrics: [], rules: [], events: [] };

async function request(path, options) {
  const response = await fetch(path, options);
  const body = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(body.error || `Request failed (${response.status})`);
  return body;
}

function option(value, label) {
  const node = document.createElement("option");
  node.value = value;
  node.textContent = label;
  return node;
}

function updateSummary() {
  document.querySelector("#source-count").textContent = state.sources.length;
  document.querySelector("#metric-count").textContent = state.metrics.length;
  document.querySelector("#rule-count").textContent = state.rules.filter((rule) => rule.enabled).length;
  document.querySelector("#event-count").textContent = state.events.filter((event) => event.status === "open").length;
}

function renderSelectors() {
  const sourceSelect = document.querySelector("#source-select");
  const metricSelect = document.querySelector("#metric-select");
  sourceSelect.replaceChildren(option("", "Select source"), ...state.sources.map((item) => option(item.id, item.name || item.id)));
  metricSelect.replaceChildren(option("", "Select metric"), ...state.metrics.map((item) => option(item.id, `${item.name || item.id}${item.unit ? ` (${item.unit})` : ""}`)));
}

function renderRules() {
  const grid = document.querySelector("#rules-grid");
  if (!state.rules.length) {
    grid.innerHTML = '<p class="empty">No evaluation rules configured.</p>';
    return;
  }
  grid.replaceChildren(...state.rules.map((rule) => {
    const node = document.createElement("article");
    node.className = "rule";
    const status = document.createElement("span");
    status.className = "rule-state";
    status.textContent = rule.enabled ? "ACTIVE" : "PAUSED";
    const title = document.createElement("h3");
    title.textContent = rule.name || rule.id;
    title.prepend(status);
    const detail = document.createElement("p");
    detail.textContent = `${rule.type || "rule"} · ${rule.metric_id || "unbound metric"} · ${rule.operator || ""} ${rule.threshold ?? ""}`;
    node.append(title, detail);
    return node;
  }));
}

function renderEvents() {
  const status = document.querySelector("#event-filter").value;
  const events = state.events.filter((event) => !status || event.status === status);
  const body = document.querySelector("#events-body");
  const empty = document.querySelector("#events-empty");
  body.replaceChildren(...events.map((event) => {
    const row = document.createElement("tr");
    const ruleCell = document.createElement("td");
    ruleCell.textContent = event.rule_id || "Unknown rule";
    const targetCell = document.createElement("td");
    targetCell.textContent = event.source_id || "Unknown source";
    const metric = document.createElement("div");
    metric.className = "subtle";
    metric.textContent = event.metric_id || "Unknown metric";
    targetCell.append(metric);
    const statusCell = document.createElement("td");
    const badge = document.createElement("span");
    badge.className = "badge";
    badge.textContent = event.status || "open";
    statusCell.append(badge);
    const observedCell = document.createElement("td");
    observedCell.textContent = event.observed_at ? new Date(event.observed_at).toLocaleString() : "Pending timestamp";
    row.append(ruleCell, targetCell, statusCell, observedCell);
    return row;
  }));
  empty.hidden = events.length > 0;
}

async function refresh() {
  const dot = document.querySelector("#health-dot");
  const label = document.querySelector("#health-label");
  dot.className = "status-dot";
  label.textContent = "Refreshing";
  try {
    await request("/healthz");
    const [sources, metrics, rules, events] = await Promise.all([
      request("/api/v1/sources"),
      request("/api/v1/metrics"),
      request("/api/v1/rules"),
      request("/api/v1/events?limit=100")
    ]);
    Object.assign(state, { sources, metrics, rules, events });
    updateSummary();
    renderSelectors();
    renderRules();
    renderEvents();
    dot.classList.add("healthy");
    label.textContent = "Service healthy";
  } catch (error) {
    dot.classList.add("failed");
    label.textContent = error.message;
  }
}

document.querySelector("#refresh").addEventListener("click", refresh);
document.querySelector("#event-filter").addEventListener("change", renderEvents);
document.querySelector("#sample-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const result = document.querySelector("#sample-result");
  result.className = "result";
  result.textContent = "Submitting sample...";
  try {
    const payload = {
      source_id: document.querySelector("#source-select").value,
      metric_id: document.querySelector("#metric-select").value,
      value: Number(document.querySelector("#sample-value").value)
    };
    const response = await request("/api/v1/telemetry", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    const count = Array.isArray(response.events) ? response.events.length : 0;
    result.textContent = `Sample accepted. ${count} event${count === 1 ? "" : "s"} emitted.`;
    await refresh();
  } catch (error) {
    result.classList.add("error");
    result.textContent = error.message;
  }
});

refresh();
