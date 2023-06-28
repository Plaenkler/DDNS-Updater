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
  document.getElementById('edit-id').value = row.querySelector('td:nth-child(1)').textContent;
  document.getElementById('delete-button').href = `/api/job/delete?ID=${row.querySelector('td:nth-child(1)').textContent}`;
  document.getElementById('edit-provider-select').value = row.querySelector('td:nth-child(2)').textContent;
  const params = JSON.parse(row.querySelector('td:nth-child(4)').textContent);
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
