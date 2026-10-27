// Utility function to encode text content
function encodeTextContent(text) {
  const div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}

// Append inputs to form create job
document.getElementById('add-provider-tiles').addEventListener('change', async (e) => {
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
    document.getElementById('add-inputs-container').scrollIntoView({ behavior: 'smooth', block: 'nearest' });
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
  const checked = document.querySelector('#add-provider-tiles input:checked');
  if (!checked) return;
  const provider = checked.value;
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

// Toggle visibility of TOTP QR code
document.getElementById('totp-qr-switch').addEventListener('change', (e) => {
  document.getElementById('totp-qr').classList.toggle('d-none', !e.target.checked);
});

// Hide TOTP QR code when settings close
document.getElementById('config-modal').addEventListener('hidden.bs.modal', () => {
  document.getElementById('totp-qr-switch').checked = false;
  document.getElementById('totp-qr').classList.add('d-none');
});

// Create log line with colored level like the CLI
function createLogLine(entry) {
  const line = document.createElement('div');
  const match = entry.match(/^(\S+ \S+ )(INF|ERR|FAT)( .*)$/);
  if (!match) {
    line.textContent = entry;
    return line;
  }
  const level = document.createElement('span');
  level.textContent = match[2];
  level.classList.add(match[2] === 'INF' ? 'text-success' : 'text-danger');
  line.append(match[1], level, match[3]);
  return line;
}

// Load log entries into modal logs
async function loadLogs(scrollToEnd) {
  try {
    const entries = await (await fetch('/api/logs')).json();
    const container = document.getElementById('logs-container');
    const atBottom = container.scrollHeight - container.scrollTop - container.clientHeight < 10;
    container.replaceChildren();
    for (const entry of entries) {
      container.appendChild(createLogLine(entry));
    }
    if (!scrollToEnd && !atBottom) return;
    container.scrollTop = container.scrollHeight;
  } catch (error) {
    console.error(error);
  }
}

// Refresh logs periodically while modal is open
let logsInterval;
document.getElementById('logs-modal').addEventListener('shown.bs.modal', () => {
  loadLogs(true);
  logsInterval = setInterval(loadLogs, 5000);
});

// Stop refreshing logs when modal closes
document.getElementById('logs-modal').addEventListener('hidden.bs.modal', () => {
  clearInterval(logsInterval);
});

// Release focus before a modal hides
document.addEventListener('hide.bs.modal', () => {
  document.activeElement.blur();
});
