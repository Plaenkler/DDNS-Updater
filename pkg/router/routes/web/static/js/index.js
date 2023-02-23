// Load input`s values for edit job modal
document.querySelectorAll('table tbody tr').forEach((row) => {
  row.addEventListener('click', (event) => {
    const [id, provider, domain, user] = event.target.parentNode.children
    document.querySelector('#edit-id').value = id.textContent
    document.querySelector('#edit-provider').value = provider.textContent
    document.querySelector('#edit-domain').value = domain.textContent
    document.querySelector('#edit-user').value = user.textContent
    document.querySelector(
      '#delete-button'
    ).href = `/api/job/delete?ID=${id.textContent}`
    new bootstrap.Modal(document.querySelector('#edit-modal')).show()
  })
})

// Add inputs to form create job
document.addEventListener("DOMContentLoaded", function() {
  const providerDropdown = document.getElementById("provider");
  const contentContainer = document.getElementById("inputs-add-container");
  providerDropdown.addEventListener("change", function() {
    const selectedProvider = providerDropdown.value;
    if (selectedProvider === "Please select") {
      contentContainer.innerHTML = "";
      return
    }
    fetch(`/api/inputs?provider=${selectedProvider}`)
    .then(response => response.text())
    .then(data => {
      contentContainer.innerHTML = data;
      console.log(data);
    })
    .catch(error => {
      console.error(error);
    });
  });
});