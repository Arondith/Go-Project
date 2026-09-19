const form = document.querySelector("#analyze-form");
const input = document.querySelector("#repository");
const statusText = document.querySelector("#status");
const results = document.querySelector("#results");

const formatNumber = new Intl.NumberFormat();

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const raw = input.value.trim().replace(/^https?:\/\/github\.com\//i, "").replace(/\/$/, "");
  const [owner, repo, ...extra] = raw.split("/");

  if (!owner || !repo || extra.length) {
    showStatus("Use the format owner/repository.", true);
    return;
  }

  showStatus("Analyzing repository…");
  results.classList.add("hidden");

  try {
    const response = await fetch(
      `/api/analyze?owner=${encodeURIComponent(owner)}&repo=${encodeURIComponent(repo)}`
    );
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Unable to analyze this repository.");
    }

    render(data);
    showStatus("");
  } catch (error) {
    showStatus(error.message, true);
  }
});

function render(data) {
  document.querySelector("#repo-name").textContent = data.fullName;
  document.querySelector("#description").textContent =
    data.description || "No repository description provided.";
  document.querySelector("#repo-link").href = data.url;
  document.querySelector("#stars").textContent = formatNumber.format(data.stars);
  document.querySelector("#forks").textContent = formatNumber.format(data.forks);
  document.querySelector("#issues").textContent = formatNumber.format(data.openIssues);
  document.querySelector("#branch").textContent = data.defaultBranch || "—";
  document.querySelector("#license").textContent = data.license || "Not specified";
  document.querySelector("#updated").textContent = new Date(data.updatedAt).toLocaleDateString();

  const languageEntries = Object.entries(data.languages || {})
    .sort((a, b) => b[1] - a[1]);

  document.querySelector("#language-count").textContent =
    `${languageEntries.length} language${languageEntries.length === 1 ? "" : "s"}`;

  const languageContainer = document.querySelector("#languages");
  languageContainer.innerHTML = "";

  if (!languageEntries.length) {
    languageContainer.textContent = "GitHub did not report language data for this repository.";
  } else {
    for (const [name, bytes] of languageEntries) {
      const percent = data.totalBytes ? (bytes / data.totalBytes) * 100 : 0;
      const row = document.createElement("div");
      row.className = "language-row";
      row.innerHTML = `
        <div class="language-meta">
          <strong>${escapeHTML(name)}</strong>
          <span>${percent.toFixed(1)}%</span>
        </div>
        <div class="bar"><span style="width: ${percent}%"></span></div>
      `;
      languageContainer.appendChild(row);
    }
  }

  results.classList.remove("hidden");
}

function showStatus(message, isError = false) {
  statusText.textContent = message;
  statusText.classList.toggle("error", isError);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}
