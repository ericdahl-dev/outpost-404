(function () {
  const REPO = "ericdahl-dev/outpost-404";
  const API = "https://api.github.com/repos/" + REPO + "/releases/latest";

  const versionEl = document.getElementById("latest-version");
  const statusEl = document.getElementById("release-status");
  const assetsEl = document.getElementById("release-assets");
  const notesEl = document.getElementById("release-notes");

  function platformLabel(name) {
    if (name.includes("darwin") && name.includes("arm64")) return "macOS (Apple Silicon)";
    if (name.includes("darwin")) return "macOS (Intel)";
    if (name.includes("windows")) return "Windows";
    if (name.includes("linux") && name.includes("arm64")) return "Linux (arm64)";
    if (name.includes("linux")) return "Linux (amd64)";
    return name;
  }

  async function loadLatestRelease() {
    try {
      const res = await fetch(API);
      if (!res.ok) throw new Error("HTTP " + res.status);
      const data = await res.json();
      const tag = data.tag_name || "";
      const name = data.name || tag;

      if (versionEl) versionEl.textContent = tag;
      if (statusEl) {
        statusEl.innerHTML =
          '<p class="release-title"><strong>' +
          escapeHtml(name) +
          '</strong> <span class="dim">' +
          escapeHtml(tag) +
          "</span></p>";
        if (data.published_at) {
          statusEl.innerHTML +=
            '<p class="dim">Published ' +
            new Date(data.published_at).toLocaleDateString(undefined, {
              year: "numeric",
              month: "short",
              day: "numeric",
            }) +
            "</p>";
        }
      }

      if (assetsEl && Array.isArray(data.assets)) {
        const binaries = data.assets.filter(function (a) {
          return (
            a.name &&
            (a.name.endsWith(".tar.gz") ||
              a.name.endsWith(".zip") ||
              a.name === "checksums.txt")
          );
        });
        if (binaries.length === 0) {
          assetsEl.innerHTML =
            '<p class="dim">No binaries attached. See <a href="' +
            escapeHtml(data.html_url) +
            '">release notes</a>.</p>';
        } else {
          const list = document.createElement("ul");
          list.className = "asset-list";
          binaries.forEach(function (asset) {
            const li = document.createElement("li");
            const a = document.createElement("a");
            a.href = asset.browser_download_url;
            a.textContent =
              asset.name === "checksums.txt"
                ? "checksums.txt"
                : platformLabel(asset.name);
            a.setAttribute("rel", "noopener noreferrer");
            li.appendChild(a);
            list.appendChild(li);
          });
          assetsEl.innerHTML = "";
          assetsEl.appendChild(list);
        }
      }

      if (notesEl && data.body) {
        notesEl.innerHTML = "<h3>Release notes</h3>" + simpleMarkdown(data.body);
      } else if (notesEl && data.html_url) {
        notesEl.innerHTML =
          '<p><a href="' +
          escapeHtml(data.html_url) +
          '" rel="noopener noreferrer">View release on GitHub</a></p>';
      }
    } catch (err) {
      if (versionEl) versionEl.textContent = "dev";
      if (statusEl) {
        statusEl.innerHTML =
          '<p class="dim">Could not load latest release. Install via Homebrew or <a href="https://github.com/' +
          REPO +
          '/releases">GitHub Releases</a>.</p>';
      }
    }
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function simpleMarkdown(body) {
    const lines = body.split("\n");
    let html = "";
    let inList = false;
    lines.forEach(function (line) {
      const t = line.trim();
      if (t.startsWith("- ")) {
        if (!inList) {
          html += "<ul>";
          inList = true;
        }
        html += "<li>" + escapeHtml(t.slice(2)) + "</li>";
      } else {
        if (inList) {
          html += "</ul>";
          inList = false;
        }
        if (t) html += "<p>" + escapeHtml(t) + "</p>";
      }
    });
    if (inList) html += "</ul>";
    return html;
  }

  loadLatestRelease();

  var entries = document.querySelectorAll('.log-entry');
  if ('IntersectionObserver' in window) {
    var io = new IntersectionObserver(function(records) {
      records.forEach(function(r) {
        if (r.isIntersecting) { r.target.classList.add('visible'); io.unobserve(r.target); }
      });
    }, { threshold: 0.08 });
    entries.forEach(function(el) { io.observe(el); });
  } else {
    entries.forEach(function(el) { el.classList.add('visible'); });
  }
})();
