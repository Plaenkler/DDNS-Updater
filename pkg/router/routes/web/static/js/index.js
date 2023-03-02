// Load input`s values for edit job modal
function loadEdit(event) {
  const [id, provider, domain, user] = event.target.parentNode.children
  document.querySelector('#edit-id').value = id.textContent
  document.querySelector('#edit-provider').value = provider.textContent
  document.querySelector('#edit-domain').value = domain.textContent
  document.querySelector('#edit-user').value = user.textContent
  document.querySelector(
    '#delete-button'
  ).href = `/api/job/delete?ID=${id.textContent}`
  new bootstrap.Modal(document.querySelector('#edit-modal')).show()
}

// Handle submit for different forms
function handleForm(event) { 
  event.preventDefault();
  const provider = event.target.querySelector("#provider").value;
  const inputs = event.target.querySelectorAll("input");
  const data = {};
  inputs.forEach(input => {
    data[input.name] = input.value;
  });
  fetch(`/api/job/create?provider=${provider}&params=${JSON.stringify(data)}`, {
    method: "POST",
  })
  .catch(error => {
    console.error(error);
  });
}

// Add inputs to form create job
function onLoad() {
  const providerDropdown = document.getElementById("provider");
  const contentContainer = document.getElementById("inputs-add-container");
  providerDropdown.addEventListener("change", function() {
    const selectedProvider = providerDropdown.value;
    if (selectedProvider === "Please select") {
      contentContainer.innerHTML = "";
      return
    }
    fetch(`/api/inputs?provider=${selectedProvider}`)
    .catch(error => {
      console.error(error);
    });
  });
}

document.getElementById("add-form").addEventListener("submit", handleForm);
document.getElementById("edit-form").addEventListener("submit", handleForm);
document.addEventListener("DOMContentLoaded", onLoad);
document.querySelectorAll('table tbody tr').forEach((row) => {
  row.addEventListener('click', loadEdit)
})