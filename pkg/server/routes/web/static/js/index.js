// Utility function to encode text content
function encodeTextContent(text) {
  const div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}

// Append inputs to form create job
document.getElementById('add-provider-select').addEventListener('change', async (e) => {
  try {
    const inputs = await (await fetch(`/api/inputs?provider=${e.target.value}`)).json();
    let html = '';
    for (const key in inputs) {
      html += `
        <div class="mb-3">
          <label for="add-input-${key}" class="form-label">${key}</label>
          <input type="text" class="form-control" id="add-input-${key}" name="${key}" value="${inputs[key]}">
        </div>
      `;
    }
    document.getElementById('add-inputs-container').innerHTML = html;
  } catch (error) {
    console.error(error);
  }
});

// Add inputs & values to form edit
document.querySelector('table tbody').addEventListener('click', function(event) {
  const row = event.target.closest('tr');
  if (!row) return;
  const idText = row.querySelector('td:nth-child(1)').textContent;
  document.getElementById('edit-id').value = idText;
  document.getElementById('delete-button').href = `/api/job/delete?ID=${encodeTextContent(idText)}`;
  document.getElementById('edit-provider-select').value = row.querySelector('td:nth-child(2)').textContent;
  const params = JSON.parse(row.querySelector('td:nth-child(4)').getAttribute('json'));
  let html = '';
  for (const key in params) {
    html += `
      <div class="mb-3">
        <label for="edit-input-${key}" class="form-label">${key}</label>
        <input type="text" class="form-control" id="edit-input-${key}" name="${key}" value="${params[key]}">
      </div>
    `;
  }
  document.getElementById('edit-inputs-container').innerHTML = html;
  new bootstrap.Modal(document.getElementById('edit-modal')).show();
});

// Send request to job api endpoint
async function sendRequest(url, method) {
  try {
    await fetch(url, {
      method: method
    });
    window.location.reload();
  } catch (error) {
    console.error(error);
  }
}

// Handle submit for form create job
document.getElementById('add-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const data = Object.fromEntries(formData.entries());
  const provider = document.getElementById('add-provider-select').value;
  const url = `/api/job/create?provider=${provider}&params=${JSON.stringify(data)}`;
  await sendRequest(url, 'POST');
});

// Handle submit for form edit job
document.getElementById('edit-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const formData = new FormData(e.target);
  const data = Object.fromEntries(formData.entries());
  const provider = document.getElementById('edit-provider-select').value;
  const id = document.getElementById('edit-id').value;
  const url = `/api/job/update?ID=${id}&provider=${provider}&params=${JSON.stringify(data)}`;
  await sendRequest(url, 'POST');
});

// Logs functionality
let logsAutoRefreshInterval = null;

// Function to fetch and display logs
async function fetchAndDisplayLogs() {
  try {
    const response = await fetch('/api/logs');
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const logs = await response.json();
    
    const logsContainer = document.getElementById('logs-container');
    if (logs.length === 0 || (logs.length === 1 && logs[0] === "")) {
      logsContainer.innerHTML = '<div class="text-center text-muted">No logs available</div>';
      return;
    }
    
    // Filter out empty entries and format logs
    const formattedLogs = logs
      .filter(log => log.trim() !== '')
      .map(log => formatLogEntry(log))
      .join('\n');
    
    logsContainer.innerHTML = formattedLogs;
    
    // Auto-scroll to bottom to show latest logs
    logsContainer.scrollTop = logsContainer.scrollHeight;
  } catch (error) {
    console.error('Error fetching logs:', error);
    document.getElementById('logs-container').innerHTML = 
      '<div class="text-danger text-center">Error loading logs: ' + error.message + '</div>';
  }
}

// Function to format individual log entries with color coding
function formatLogEntry(logEntry) {
  if (!logEntry || logEntry.trim() === '') return '';
  
  // Parse log format: "2025/08/01 06:31:03 INF origin:main-main line:16 message:started database connection"
  const logRegex = /^(\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}) (INF|ERR|FAT) (.*)/;
  const match = logEntry.match(logRegex);
  
  if (match) {
    const [, timestamp, level, rest] = match;
    let levelClass = '';
    let levelIcon = '';
    
    switch (level) {
      case 'INF':
        levelClass = 'text-success';
        levelIcon = '●';
        break;
      case 'ERR':
        levelClass = 'text-danger';
        levelIcon = '●';
        break;
      case 'FAT':
        levelClass = 'text-danger';
        levelIcon = '●';
        break;
      default:
        levelClass = 'text-info';
        levelIcon = '●';
    }
    
    return `<span class="text-muted">${timestamp}</span> <span class="${levelClass}">${levelIcon} ${level}</span> ${encodeTextContent(rest)}`;
  }
  
  // If log doesn't match expected format, just return it as-is
  return encodeTextContent(logEntry);
}

// Refresh logs button handler
document.getElementById('refresh-logs-btn').addEventListener('click', fetchAndDisplayLogs);

// Clear logs display button handler
document.getElementById('clear-logs-display-btn').addEventListener('click', () => {
  document.getElementById('logs-container').innerHTML = '<div class="text-center text-muted">Display cleared</div>';
});

// Auto-refresh checkbox handler
document.getElementById('auto-refresh-logs').addEventListener('change', (e) => {
  if (e.target.checked) {
    // Start auto-refresh every 5 seconds
    logsAutoRefreshInterval = setInterval(fetchAndDisplayLogs, 5000);
  } else {
    // Stop auto-refresh
    if (logsAutoRefreshInterval) {
      clearInterval(logsAutoRefreshInterval);
      logsAutoRefreshInterval = null;
    }
  }
});

// Load logs when modal is shown
document.getElementById('logs-modal').addEventListener('shown.bs.modal', fetchAndDisplayLogs);

// Clean up auto-refresh when modal is hidden
document.getElementById('logs-modal').addEventListener('hidden.bs.modal', () => {
  if (logsAutoRefreshInterval) {
    clearInterval(logsAutoRefreshInterval);
    logsAutoRefreshInterval = null;
  }
  // Uncheck auto-refresh
  document.getElementById('auto-refresh-logs').checked = false;
});
