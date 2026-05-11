export function renderAuthShell(title: string, subtitle: string, content: string) {
  const app = document.querySelector<HTMLDivElement>('#app')
  if (!app)
    return

  app.innerHTML = `
    <main class="auth-shell">
      <section class="auth-panel">
        <header class="auth-header">
          <h1>${title}</h1>
          <p>${subtitle}</p>
        </header>
        ${content}
      </section>
    </main>
  `
}

export function attachFormMessage(message: string, isError = true) {
  const el = document.querySelector<HTMLDivElement>('#form-message')
  if (!el)
    return
  el.textContent = message
  el.className = isError ? 'form-message is-error' : 'form-message is-success'
}
