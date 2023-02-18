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