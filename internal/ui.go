package internal

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>LogTail</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }

  body {
    font-family: "SF Mono", "Cascadia Code", "Fira Code", Consolas, monospace;
    background: #1a1b26;
    color: #c0caf5;
    height: 100vh;
    display: flex;
    flex-direction: column;
  }

  header {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: #16161e;
    border-bottom: 1px solid #292e42;
    flex-shrink: 0;
  }

  header h1 {
    font-size: 16px;
    font-weight: 600;
    color: #7aa2f7;
    margin-right: 8px;
  }

  #search {
    flex: 1;
    max-width: 400px;
    padding: 6px 12px;
    background: #1a1b26;
    border: 1px solid #292e42;
    border-radius: 6px;
    color: #c0caf5;
    font-family: inherit;
    font-size: 13px;
    outline: none;
  }

  #search:focus { border-color: #7aa2f7; }
  #search::placeholder { color: #565f89; }

  .controls {
    display: flex;
    gap: 8px;
    margin-left: auto;
  }

  button {
    padding: 6px 14px;
    background: #292e42;
    border: 1px solid #3b4261;
    border-radius: 6px;
    color: #c0caf5;
    font-family: inherit;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.15s;
  }

  button:hover { background: #3b4261; }
  button.active { background: #7aa2f7; color: #1a1b26; border-color: #7aa2f7; }

  #status {
    font-size: 12px;
    color: #565f89;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  #status .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #9ece6a;
  }

  #status .dot.disconnected { background: #f7768e; }

  #log-container {
    flex: 1;
    overflow-y: auto;
    padding: 8px 0;
  }

  .log-line {
    padding: 1px 16px;
    font-size: 13px;
    line-height: 20px;
    white-space: pre-wrap;
    word-break: break-all;
    border-left: 3px solid transparent;
  }

  .log-line:hover { background: rgba(122, 162, 247, 0.06); }
  .log-line.hidden { display: none; }
  .log-line .line-num {
    display: inline-block;
    width: 50px;
    color: #3b4261;
    text-align: right;
    margin-right: 12px;
    user-select: none;
  }

  .log-line.level-debug { border-left-color: #9ece6a; }
  .log-line.level-info { border-left-color: #7aa2f7; }
  .log-line.level-warn { border-left-color: #e0af68; }
  .log-line.level-error { border-left-color: #f7768e; }
  .log-line.level-fatal { border-left-color: #ff5370; }

  .log-line mark {
    background: #e0af68;
    color: #1a1b26;
    border-radius: 2px;
    padding: 0 2px;
  }

  footer {
    padding: 6px 16px;
    background: #16161e;
    border-top: 1px solid #292e42;
    font-size: 12px;
    color: #565f89;
    display: flex;
    justify-content: space-between;
    flex-shrink: 0;
  }
</style>
</head>
<body>
  <header>
    <h1>LogTail</h1>
    <input type="text" id="search" placeholder="Search logs... (Ctrl+F)" autocomplete="off" />
    <div class="controls">
      <button id="pause-btn" title="Pause/Resume streaming">Pause</button>
      <button id="clear-btn" title="Clear log view">Clear</button>
      <button id="scroll-btn" class="active" title="Toggle auto-scroll">Auto-scroll</button>
    </div>
    <div id="status">
      <span class="dot" id="status-dot"></span>
      <span id="status-text">Connecting...</span>
    </div>
  </header>

  <div id="log-container"></div>

  <footer>
    <span id="line-count">0 lines</span>
    <span id="match-count"></span>
  </footer>

  <script>
    const container = document.getElementById('log-container');
    const searchInput = document.getElementById('search');
    const pauseBtn = document.getElementById('pause-btn');
    const clearBtn = document.getElementById('clear-btn');
    const scrollBtn = document.getElementById('scroll-btn');
    const statusDot = document.getElementById('status-dot');
    const statusText = document.getElementById('status-text');
    const lineCountEl = document.getElementById('line-count');
    const matchCountEl = document.getElementById('match-count');

    let lineNumber = 0;
    let autoScroll = true;
    let paused = false;
    let eventSource = null;
    let pauseBuffer = [];
    let searchTerm = '';

    function detectLevel(text) {
      const t = text.toUpperCase();
      if (t.includes('"FATAL"') || t.includes('[FATAL]') || t.includes(' FATAL ')) return 'fatal';
      if (t.includes('"ERROR"') || t.includes('[ERROR]') || t.includes(' ERROR ')) return 'error';
      if (t.includes('"WARN"') || t.includes('[WARN]') || t.includes(' WARN ')) return 'warn';
      if (t.includes('"INFO"') || t.includes('[INFO]') || t.includes(' INFO ')) return 'info';
      if (t.includes('"DEBUG"') || t.includes('[DEBUG]') || t.includes(' DEBUG ')) return 'debug';
      return '';
    }

    function escapeHTML(str) {
      return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    }

    function highlightText(html, term) {
      if (!term) return html;
      const escaped = term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      return html.replace(new RegExp('(' + escaped + ')', 'gi'), '<mark>$1</mark>');
    }

    function addLine(text) {
      lineNumber++;
      const level = detectLevel(text);
      const div = document.createElement('div');
      div.className = 'log-line' + (level ? ' level-' + level : '');
      div.dataset.raw = text;

      let html = '<span class="line-num">' + lineNumber + '</span>';
      let content = escapeHTML(text);
      if (searchTerm) {
        content = highlightText(content, searchTerm);
        if (!text.toLowerCase().includes(searchTerm.toLowerCase())) {
          div.classList.add('hidden');
        }
      }
      html += content;
      div.innerHTML = html;

      container.appendChild(div);

      if (autoScroll) {
        container.scrollTop = container.scrollHeight;
      }

      updateCounts();
    }

    function updateCounts() {
      const total = container.children.length;
      const visible = container.querySelectorAll('.log-line:not(.hidden)').length;
      lineCountEl.textContent = total + ' lines';
      if (searchTerm) {
        matchCountEl.textContent = visible + ' matches';
      } else {
        matchCountEl.textContent = '';
      }
    }

    function applySearch() {
      const lines = container.querySelectorAll('.log-line');
      lines.forEach(line => {
        const raw = line.dataset.raw || '';
        let content = escapeHTML(raw);

        if (searchTerm) {
          content = highlightText(content, searchTerm);
          if (!raw.toLowerCase().includes(searchTerm.toLowerCase())) {
            line.classList.add('hidden');
          } else {
            line.classList.remove('hidden');
          }
        } else {
          line.classList.remove('hidden');
        }

        const numSpan = line.querySelector('.line-num');
        const num = numSpan ? numSpan.outerHTML : '';
        line.innerHTML = num + content;
      });
      updateCounts();
    }

    function connect() {
      if (eventSource) eventSource.close();

      eventSource = new EventSource('/logs/stream');

      eventSource.onopen = function() {
        statusDot.className = 'dot';
        statusText.textContent = 'Connected';
      };

      eventSource.onmessage = function(e) {
        if (paused) {
          pauseBuffer.push(e.data);
          statusText.textContent = 'Paused (' + pauseBuffer.length + ' buffered)';
          return;
        }
        addLine(e.data);
      };

      eventSource.onerror = function() {
        statusDot.className = 'dot disconnected';
        statusText.textContent = 'Disconnected - retrying...';
      };
    }

    searchInput.addEventListener('input', function() {
      searchTerm = this.value;
      applySearch();
    });

    pauseBtn.addEventListener('click', function() {
      paused = !paused;
      if (paused) {
        this.classList.add('active');
        this.textContent = 'Resume';
        statusText.textContent = 'Paused';
      } else {
        this.classList.remove('active');
        this.textContent = 'Pause';
        pauseBuffer.forEach(addLine);
        pauseBuffer = [];
        statusText.textContent = 'Connected';
      }
    });

    clearBtn.addEventListener('click', function() {
      container.innerHTML = '';
      lineNumber = 0;
      updateCounts();
    });

    scrollBtn.addEventListener('click', function() {
      autoScroll = !autoScroll;
      this.classList.toggle('active', autoScroll);
      if (autoScroll) {
        container.scrollTop = container.scrollHeight;
      }
    });

    container.addEventListener('scroll', function() {
      const atBottom = container.scrollTop + container.clientHeight >= container.scrollHeight - 30;
      if (!atBottom && autoScroll) {
        autoScroll = false;
        scrollBtn.classList.remove('active');
      }
    });

    document.addEventListener('keydown', function(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
        e.preventDefault();
        searchInput.focus();
        searchInput.select();
      }
      if (e.key === 'Escape') {
        searchInput.value = '';
        searchTerm = '';
        applySearch();
        searchInput.blur();
      }
    });

    connect();
  </script>
</body>
</html>`
