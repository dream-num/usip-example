import './style.css'
import { setupUniver } from './setup-univer'

type APIError = { error?: string }

type User = {
  userId: string
  nickname: string
  username: string
}

type AuthResp = {
  user: User
}

type FileItem = {
  id: number
  name: string
  unitId: string
  unitType: number
  updatedAt: string
  openUrl: string
  exportUrl: string
}

type FilesResp = {
  userId: string
  files: FileItem[]
}

async function apiFetch<T>(input: RequestInfo | URL, init?: RequestInit): Promise<T> {
  const resp = await fetch(input, init)

  if (resp.status === 401) {
    if (location.pathname !== '/login')
      location.href = '/login'
    throw new Error('unauthorized')
  }

  if (!resp.ok) {
    let message = `request failed: ${resp.status}`
    try {
      const payload = (await resp.json()) as APIError
      if (payload.error)
        message = payload.error
    }
    catch {
      // ignore
    }
    throw new Error(message)
  }

  if (resp.status === 204)
    return {} as T

  return (await resp.json()) as T
}

function renderAuthShell(title: string, subtitle: string, content: string) {
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

function attachMessage(message: string, isError = true) {
  const el = document.querySelector<HTMLDivElement>('#form-message')
  if (!el)
    return
  el.textContent = message
  el.className = isError ? 'form-message is-error' : 'form-message is-success'
}

function renderLogin() {
  renderAuthShell(
    'Welcome Back',
    'Sign in to continue to your workspace.',
    `<form id="login-form" class="auth-form">
      <label for="login-username"><b>Username</b></label>
      <input id="login-username" type="text" name="username" required>
      <label for="login-password"><b>Password</b></label>
      <input id="login-password" type="password" name="password" required>
      <div id="form-message" class="form-message"></div>
      <button type="submit" class="auth-submit">Login</button>
    </form>
    <footer class="auth-footer">
      <span>No account yet?</span>
      <a href="/register">Go to register</a>
    </footer>`,
  )

  const form = document.querySelector<HTMLFormElement>('#login-form')
  if (!form)
    return

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    try {
      await apiFetch<AuthResp>('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: String(formData.get('username') ?? ''),
          password: String(formData.get('password') ?? ''),
        }),
      })
      location.href = '/files'
    }
    catch (error) {
      attachMessage((error as Error).message)
    }
  })
}

function renderRegister() {
  renderAuthShell(
    'Create Account',
    'Join and start using your workspace in seconds.',
    `<form id="register-form" class="auth-form">
      <label for="register-nickname"><b>Nickname</b></label>
      <input id="register-nickname" type="text" name="nickname" required>
      <label for="register-username"><b>Username</b></label>
      <input id="register-username" type="text" name="username" required>
      <label for="register-password"><b>Password</b></label>
      <input id="register-password" type="password" name="password" required>
      <div id="form-message" class="form-message"></div>
      <button type="submit" class="auth-submit">Register</button>
    </form>
    <footer class="auth-footer">
      <span>Already have an account?</span>
      <a href="/login">Go to login</a>
    </footer>`,
  )

  const form = document.querySelector<HTMLFormElement>('#register-form')
  if (!form)
    return

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    try {
      await apiFetch<AuthResp>('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nickname: String(formData.get('nickname') ?? ''),
          username: String(formData.get('username') ?? ''),
          password: String(formData.get('password') ?? ''),
        }),
      })
      location.href = '/files'
    }
    catch (error) {
      attachMessage((error as Error).message)
    }
  })
}

function renderFilesShell() {
  const app = document.querySelector<HTMLDivElement>('#app')
  if (!app)
    return

  app.innerHTML = `
    <div id="list-page">
      <div class="demo-topbar">
        <div class="demo-title">
          <h1>Demo2 Files</h1>
          <p>Fast file operations with a clean, modern workspace.</p>
        </div>
        <div class="demo-user">
          <div class="div-image" title="Current user avatar">
            <img id="user-avatar" alt="User avatar" loading="lazy" decoding="async" referrerpolicy="no-referrer" />
          </div>
          <button id="logout-btn" type="button">Logout</button>
        </div>
      </div>
      <div class="demo-actions">
        <button id="new-btn" class="demo-btn-primary" type="button">+ New File</button>
        <button id="import-btn" class="demo-btn-secondary" type="button">Import File</button>
        <button id="delete-btn" class="demo-btn-danger" type="button">Delete Selected</button>
      </div>
      <input type="file" id="file-input" style="display:none;" />
      <div id="div-form" class="demo-card">
        <form id="new-form" enctype="multipart/form-data">
          <label><b>Name</b></label>
          <input type="text" placeholder="Enter Name" name="name" required>
          <div class="form-actions">
            <button class="demo-btn-primary" type="submit">Create</button>
            <button id="form-cancel-btn" class="demo-btn-secondary" type="button">Cancel</button>
          </div>
        </form>
      </div>
      <div class="demo-card">
        <div class="demo-list-head"><h2>File list</h2></div>
        <div id="files-container" class="file-list"></div>
      </div>
      <dialog id="dialog" class="dialog-panel">
        <div class="dialog-header">
          <h3>Invite collaborators</h3>
          <p>Select users and assign permission.</p>
        </div>
        <div class="dialog-role-row">
          <label for="select-role">Role</label>
          <select id="select-role">
            <option value="reader">Reader</option>
            <option value="editor">Editor</option>
          </select>
        </div>
        <div id="dialog-msg" class="dialog-list"></div>
        <div class="dialog-actions">
          <button id="dialog-ok-btn" class="demo-btn-primary" type="button">Invite</button>
          <button id="dialog-cancel-btn" class="demo-btn-secondary" type="button">Cancel</button>
        </div>
      </dialog>
    </div>
  `
}

function escapeHtml(text: string) {
  return text
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

async function renderFiles() {
  renderFilesShell()
  const filesResp = await apiFetch<FilesResp>('/api/files')

  const fileContainer = document.querySelector<HTMLDivElement>('#files-container')
  if (!fileContainer)
    return

  fileContainer.innerHTML = filesResp.files.map(file => `
    <div class="file-row hover-effect" data-file-id="${file.id}">
      <input type="checkbox" name="fileIds" class="fileCheckbox file-check" value="${file.id}"/>
      <div class="file-name">
        ${file.openUrl ? `<a href="${file.openUrl}">${escapeHtml(file.name)}.xlsx</a>` : escapeHtml(file.name)}
      </div>
      <label class="file-updated">${escapeHtml(file.updatedAt)}</label>
      <div class="file-actions">
        <a class="file-action-link" href="${file.exportUrl}">Export</a>
        <button class="demo-btn-secondary invite-btn" type="button" data-file-id="${file.id}">Invite</button>
      </div>
    </div>
  `).join('')

  const avatar = document.querySelector<HTMLImageElement>('#user-avatar')
  if (avatar) {
    avatar.src = `/user/avatar/${filesResp.userId}`
    avatar.onerror = () => {
      avatar.onerror = null
      avatar.src = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 72"><defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stop-color="%2360a5fa"/><stop offset="100%" stop-color="%233b82f6"/></linearGradient></defs><circle cx="36" cy="36" r="34" fill="url(%23g)"/><circle cx="36" cy="36" r="31" fill="rgba(255,255,255,0.12)"/><text x="36" y="44" font-size="30" font-family="Arial,sans-serif" text-anchor="middle" fill="white" font-weight="700">U</text></svg>'
    }
  }

  fileContainer.querySelectorAll<HTMLDivElement>('.file-row').forEach((row) => {
    row.addEventListener('click', (event) => {
      const target = event.target as HTMLElement
      if (target.closest('.file-actions') || target.matches('a') || target.matches('button') || target.matches('input[type="checkbox"]'))
        return
      const checkbox = row.querySelector<HTMLInputElement>('.fileCheckbox')
      if (checkbox)
        checkbox.checked = !checkbox.checked
    })
  })

  const logoutBtn = document.querySelector<HTMLButtonElement>('#logout-btn')
  logoutBtn?.addEventListener('click', async () => {
    await apiFetch<{ ok: boolean }>('/api/auth/logout', { method: 'POST' })
    location.href = '/login'
  })

  const formWrap = document.querySelector<HTMLDivElement>('#div-form')
  const newBtn = document.querySelector<HTMLButtonElement>('#new-btn')
  const cancelBtn = document.querySelector<HTMLButtonElement>('#form-cancel-btn')
  const newForm = document.querySelector<HTMLFormElement>('#new-form')

  newBtn?.addEventListener('click', () => formWrap?.classList.add('is-visible'))
  cancelBtn?.addEventListener('click', () => formWrap?.classList.remove('is-visible'))

  newForm?.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(newForm)
    formData.set('type', 'sheet')
    const resp = await fetch('/file/new', { method: 'POST', body: formData })
    if (resp.redirected) {
      location.href = resp.url
      return
    }
    location.reload()
  })

  const importBtn = document.querySelector<HTMLButtonElement>('#import-btn')
  const fileInput = document.querySelector<HTMLInputElement>('#file-input')
  importBtn?.addEventListener('click', () => fileInput?.click())
  fileInput?.addEventListener('change', async () => {
    const file = fileInput.files?.[0]
    if (!file)
      return
    if (!file.name.toLowerCase().endsWith('.xlsx')) {
      alert('Only .xlsx import is supported')
      return
    }

    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', 'sheet')
    formData.append('name', file.name.replace(/\.xlsx$/i, ''))

    const resp = await fetch('/file/import', { method: 'POST', body: formData })
    if (resp.redirected) {
      location.href = resp.url
      return
    }
    location.reload()
  })

  const deleteBtn = document.querySelector<HTMLButtonElement>('#delete-btn')
  deleteBtn?.addEventListener('click', async () => {
    const params = new URLSearchParams()
    document.querySelectorAll<HTMLInputElement>('.fileCheckbox:checked').forEach((checkbox) => {
      params.append('fileIds', checkbox.value)
    })

    if (!params.size) {
      alert('Please select files to delete')
      return
    }

    await fetch(`/file?${params.toString()}`, { method: 'DELETE' })
    location.reload()
  })

  const dialog = document.querySelector<HTMLDialogElement>('#dialog')
  const dialogMsg = document.querySelector<HTMLDivElement>('#dialog-msg')
  const roleSelect = document.querySelector<HTMLSelectElement>('#select-role')
  const okBtn = document.querySelector<HTMLButtonElement>('#dialog-ok-btn')
  const cancelInviteBtn = document.querySelector<HTMLButtonElement>('#dialog-cancel-btn')

  let inviteFileId = 0
  let invite: string[] = []

  async function getPeople(next: number) {
    const resp = await apiFetch<{ users: Array<{ user_id: string, nickname: string }>, next: number }>(`/user/people?next=${next}&size=10`)
    return resp
  }

  async function renderPeople() {
    const data = await getPeople(0)
    if (!dialogMsg)
      return

    dialogMsg.innerHTML = ''
    data.users.forEach((person) => {
      const row = document.createElement('div')
      row.className = 'dialog-people-row'
      const label = document.createElement('span')
      label.className = 'dialog-people-name'
      label.innerText = person.nickname
      const dot = document.createElement('span')
      dot.className = 'dialog-people-dot'
      dot.innerText = person.nickname.slice(0, 1).toUpperCase()
      const checkbox = document.createElement('input')
      checkbox.type = 'checkbox'
      checkbox.className = 'inviteCheckbox'
      checkbox.value = person.user_id
      checkbox.checked = invite.includes(person.user_id)
      checkbox.onchange = () => {
        if (checkbox.checked)
          invite = [...new Set([...invite, checkbox.value])]
        else
          invite = invite.filter(id => id !== checkbox.value)
      }
      row.addEventListener('click', (event) => {
        const target = event.target as HTMLElement
        if (target.closest('input[type="checkbox"]'))
          return
        checkbox.click()
      })
      row.append(dot, label, checkbox)
      dialogMsg.appendChild(row)
    })
  }

  fileContainer.querySelectorAll<HTMLButtonElement>('.invite-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
      invite = []
      inviteFileId = Number(btn.dataset.fileId ?? 0)
      await renderPeople()
      dialog?.showModal()
    })
  })

  okBtn?.addEventListener('click', async () => {
    if (!inviteFileId)
      return
    await fetch('/file/join', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        userIds: invite,
        fileId: inviteFileId,
        role: roleSelect?.value ?? 'reader',
      }),
    })
    dialog?.close()
  })

  cancelInviteBtn?.addEventListener('click', () => dialog?.close())
}

function renderSheet() {
  const app = document.querySelector<HTMLDivElement>('#app')
  if (!app)
    return

  app.innerHTML = '<div id="univer"></div>'
  const univerAPI = setupUniver()
  window.univerAPI = univerAPI
}

async function main() {
  const path = window.location.pathname

  if (path === '/login') {
    renderLogin()
    return
  }

  if (path === '/register') {
    renderRegister()
    return
  }

  if (path === '/files') {
    await renderFiles()
    return
  }

  renderSheet()
}

main().catch((error) => {
  // eslint-disable-next-line no-console
  console.error(error)
})
