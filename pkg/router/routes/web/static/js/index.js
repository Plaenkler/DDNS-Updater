// Append inputs to form create job
const addProviderSelect = document.getElementById('add-provider-select');
const addInputsContainer = document.getElementById('add-inputs-container');
addProviderSelect.addEventListener('change', async (e) => {
  try {
    addInputsContainer.innerHTML = await (await fetch(`/api/inputs?provider=${e.target.value}`)).text();
  } catch (error) {
    console.error(error);
  }
});

// Handle submit for form create job
const addJobForm = document.getElementById('add-form');
addJobForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  let data = {};
  e.target.querySelectorAll('input').forEach(input => {
    data[input.name] = input.value;
  });
  const provider = e.target.querySelector('select').value;
  try {
    await fetch(`/api/job/create?provider=${provider}&params=${JSON.stringify(data)}`, {
      method: 'POST',
    });
    window.location.reload();
  } catch (error) {
    console.error(error);
  }
});

// Add inputs & values to form edit
const editProviderSelect = document.getElementById('edit-provider-select');
const editModal = document.getElementById('edit-modal');
const deleteButton = document.getElementById('delete-button');
const editInputsContainer = document.getElementById('edit-inputs-container');
const table = document.querySelector('table tbody');
table.addEventListener('click', function(event) {
  const row = event.target.closest('tr');
  if (!row) return;
  deleteButton.href = `/api/job/delete?ID=${row.querySelector('td:nth-child(1)').textContent}`
  editProviderSelect.value = row.querySelector('td:nth-child(2)').textContent
  const params = JSON.parse(row.querySelector('td:nth-child(4)').textContent)
  let html = ''
  for (const key in params) {
    html += 
    `
    <div class="mb-3">
      <label for="edit-${key}" class="form-label">${key}</label>
      <input type="text" class="form-control" id="edit-${key}" name="${key}" value="${params[key]}">
    </div>
    `
  }
  editInputsContainer.innerHTML = html
  new bootstrap.Modal(editModal).show()
})

// Handle submit for form edit
// const editJobForm = document.getElementById('edit-form');
// editJobForm.addEventListener('submit', async (e) => {
//   e.preventDefault();
//   let data = {};
//   e.target.querySelectorAll('input').forEach(input => {
//     data[input.name] = input.value;
//   });
//   const provider = e.target.querySelector('select').value;
//   try {
//     await fetch(`/api/job/update?provider=${provider}&params=${JSON.stringify(data)}`, {
//       method: 'POST',
//     });
//     window.location.reload();
//   } catch (error) {
//     console.error(error);
//   }
// });
