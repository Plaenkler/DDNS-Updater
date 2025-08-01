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

// Handle backup export
document.getElementById('export-backup-btn').addEventListener('click', async () => {
  try {
    const response = await fetch('/api/backup/export');
    if (!response.ok) {
      throw new Error('Failed to export backup');
    }
    
    // Get filename from Content-Disposition header
    const contentDisposition = response.headers.get('Content-Disposition');
    let filename = 'ddns-backup.json';
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="(.+)"/);
      if (filenameMatch) {
        filename = filenameMatch[1];
      }
    }
    
    // Create download
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
    
    // Close modal
    bootstrap.Modal.getInstance(document.getElementById('backup-modal')).hide();
  } catch (error) {
    console.error('Export failed:', error);
    alert('Failed to export backup. Please try again.');
  }
});

// Handle backup import
document.getElementById('import-backup-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const fileInput = document.getElementById('backup-file-input');
  const file = fileInput.files[0];
  
  if (!file) {
    alert('Please select a backup file');
    return;
  }
  
  if (!confirm('This will replace all current jobs and configuration. Are you sure you want to continue?')) {
    return;
  }
  
  try {
    const formData = new FormData();
    formData.append('backup', file);
    
    const response = await fetch('/api/backup/import', {
      method: 'POST',
      body: formData
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'Failed to import backup');
    }
    
    const result = await response.json();
    alert(result.message || 'Backup imported successfully');
    
    // Close modal and reload page
    bootstrap.Modal.getInstance(document.getElementById('backup-modal')).hide();
    window.location.reload();
  } catch (error) {
    console.error('Import failed:', error);
    alert('Failed to import backup: ' + error.message);
  }
});
