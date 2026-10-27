// Create input for job params
function createInput(prefix, key, value) {
  const inputId = `${prefix}-input-${key}`;
  const label = document.createElement('label');
  label.htmlFor = inputId;
  label.classList.add('form-label');
  label.textContent = key;
  const input = document.createElement('input');
  input.type = 'text';
  input.classList.add('form-control');
  input.id = inputId;
  input.name = key;
  input.value = value;
  const wrapper = document.createElement('div');
  wrapper.classList.add('mb-3');
  wrapper.appendChild(label);
  wrapper.appendChild(input);
  return wrapper;
}

// Render inputs for job params into form
function renderInputs(prefix, params) {
  const container = document.getElementById(`${prefix}-inputs-container`);
  container.replaceChildren();
  for (const key of Object.keys(params)) {
    const input = createInput(prefix, key, params[key]);
    container.appendChild(input);
  }
  return container;
}

// Fetch JSON from api endpoint
async function fetchJSON(url) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`${url} responded ${response.status}`);
  }
  return await response.json();
}

// Send form data to api endpoint
async function postForm(url, data) {
  const body = new URLSearchParams(data);
  const response = await fetch(url, { method: 'POST', body: body });
  if (!response.ok) {
    throw new Error(`${url} responded ${response.status}`);
  }
  window.location.reload();
}

// Read job params from form
function readParams(form) {
  const formData = new FormData(form);
  const params = Object.fromEntries(formData);
  return JSON.stringify(params);
}

// Load inputs of selected provider into form create job
async function loadInputs(event) {
  const tile = event.target;
  const query = new URLSearchParams({ provider: tile.value });
  const inputs = await fetchJSON(`/api/inputs?${query}`);
  if (!tile.checked) {
    return;
  }
  const container = renderInputs('add', inputs);
  container.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

// Fill form edit job with clicked row
function openEdit(event) {
  const row = event.target.closest('tr');
  if (!row) {
    return;
  }
  const id = row.cells[0].textContent;
  const provider = row.cells[1].textContent;
  const params = JSON.parse(row.cells[3].getAttribute('json'));
  const query = new URLSearchParams({ ID: id });
  document.getElementById('edit-id').value = id;
  document.getElementById('delete-button').href = `/api/job/delete?${query}`;
  document.getElementById('edit-provider-select').value = provider;
  renderInputs('edit', params);
  const modal = bootstrap.Modal.getOrCreateInstance(document.getElementById('edit-modal'));
  modal.show();
}

// Submit form create job
async function createJob(event) {
  event.preventDefault();
  const checked = document.querySelector('#add-provider-tiles input:checked');
  if (!checked) {
    return;
  }
  const data = {
    provider: checked.value,
    params: readParams(event.target)
  };
  await postForm('/api/job/create', data);
}

// Submit form edit job
async function updateJob(event) {
  event.preventDefault();
  const data = {
    ID: document.getElementById('edit-id').value,
    provider: document.getElementById('edit-provider-select').value,
    params: readParams(event.target)
  };
  await postForm('/api/job/update', data);
}

// Toggle visibility of TOTP QR code
function toggleTotpQR(event) {
  const hidden = !event.target.checked;
  document.getElementById('totp-qr').classList.toggle('d-none', hidden);
}

// Hide TOTP QR code when settings close
function hideTotpQR() {
  document.getElementById('totp-qr-switch').checked = false;
  document.getElementById('totp-qr').classList.add('d-none');
}

// Get color class of log level like the CLI
function getLevelClass(levelName) {
  if (levelName === 'INF') {
    return 'text-success';
  }
  return 'text-danger';
}

// Create log line with colored level like the CLI
function createLogLine(entry) {
  const line = document.createElement('div');
  const match = entry.match(/^(\S+ \S+ )(INF|ERR|FAT)( .*)$/);
  if (!match) {
    line.textContent = entry;
    return line;
  }
  const timestamp = match[1];
  const levelName = match[2];
  const message = match[3];
  const level = document.createElement('span');
  level.textContent = levelName;
  level.classList.add(getLevelClass(levelName));
  line.appendChild(document.createTextNode(timestamp));
  line.appendChild(level);
  line.appendChild(document.createTextNode(message));
  return line;
}

// Load log entries into modal logs
async function loadLogs(scrollToEnd) {
  const entries = await fetchJSON('/api/logs');
  const container = document.getElementById('logs-container');
  const distanceToEnd = container.scrollHeight - container.scrollTop - container.clientHeight;
  const atEnd = distanceToEnd < 10;
  container.replaceChildren();
  for (const entry of entries) {
    const line = createLogLine(entry);
    container.appendChild(line);
  }
  if (!scrollToEnd && !atEnd) {
    return;
  }
  container.scrollTop = container.scrollHeight;
}

// Interval of log refresh while modal is open
let logsInterval;

// Refresh logs without jumping to the end
function refreshLogs() {
  loadLogs(false);
}

// Load logs and refresh them periodically while modal is open
function startLogs() {
  loadLogs(true);
  logsInterval = setInterval(refreshLogs, 5000);
}

// Stop refreshing logs when modal closes
function stopLogs() {
  clearInterval(logsInterval);
}

// Release focus before a modal hides
function releaseFocus() {
  document.activeElement.blur();
}

// Register event listeners
document.getElementById('add-provider-tiles').addEventListener('change', loadInputs);
document.querySelector('table tbody').addEventListener('click', openEdit);
document.getElementById('add-form').addEventListener('submit', createJob);
document.getElementById('edit-form').addEventListener('submit', updateJob);
document.getElementById('totp-qr-switch').addEventListener('change', toggleTotpQR);
document.getElementById('config-modal').addEventListener('hidden.bs.modal', hideTotpQR);
document.getElementById('logs-modal').addEventListener('shown.bs.modal', startLogs);
document.getElementById('logs-modal').addEventListener('hidden.bs.modal', stopLogs);
document.addEventListener('hide.bs.modal', releaseFocus);
