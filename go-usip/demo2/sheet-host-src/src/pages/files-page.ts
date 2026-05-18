import { wireInviteDialog } from '../components/invite-dialog'
import { openMembersDialog } from '../components/members-dialog'
import { logout } from '../services/auth-service'
import {
  createSheet,
  deleteFiles,
  fetchFiles,
  importSheet,
} from '../services/files-service'
import type { FileItem } from '../types/files'
import { escapeHtml } from '../utils/html'

function renderFilesShell() {
  const app = document.querySelector<HTMLDivElement>('#app')
  if (!app)
    return

  app.innerHTML = `
    <div id="list-page">
      <div class="demo-topbar">
        <div class="demo-title">
          <h1>Workspace</h1>
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

function wireAvatar(userId: string) {
  const avatar = document.querySelector<HTMLImageElement>('#user-avatar')
  if (!avatar)
    return

  avatar.src = `/user/avatar/${userId}`
  avatar.onerror = () => {
    avatar.onerror = null
    avatar.src = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 72"><defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stop-color="%2360a5fa"/><stop offset="100%" stop-color="%233b82f6"/></linearGradient></defs><circle cx="36" cy="36" r="34" fill="url(%23g)"/><circle cx="36" cy="36" r="31" fill="rgba(255,255,255,0.12)"/><text x="36" y="44" font-size="30" font-family="Arial,sans-serif" text-anchor="middle" fill="white" font-weight="700">U</text></svg>'
  }
}

function wireFileRowToggle(fileContainer: HTMLDivElement) {
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
}

function wireCommonActions(fileContainer: HTMLDivElement, currentUserId: string) {
  const logoutBtn = document.querySelector<HTMLButtonElement>('#logout-btn')
  const formWrap = document.querySelector<HTMLDivElement>('#div-form')
  const newBtn = document.querySelector<HTMLButtonElement>('#new-btn')
  const cancelBtn = document.querySelector<HTMLButtonElement>('#form-cancel-btn')
  const newForm = document.querySelector<HTMLFormElement>('#new-form')
  const importBtn = document.querySelector<HTMLButtonElement>('#import-btn')
  const fileInput = document.querySelector<HTMLInputElement>('#file-input')
  const deleteBtn = document.querySelector<HTMLButtonElement>('#delete-btn')

  logoutBtn?.addEventListener('click', async () => {
    await logout()
    location.href = '/login'
  })

  newBtn?.addEventListener('click', () => formWrap?.classList.add('is-visible'))
  cancelBtn?.addEventListener('click', () => formWrap?.classList.remove('is-visible'))

  newForm?.addEventListener('submit', async (event) => {
    event.preventDefault()
    const name = String(new FormData(newForm).get('name') ?? '')
    const resp = await createSheet(name)
    if (resp.redirected) {
      location.href = resp.url
      return
    }
    location.reload()
  })

  importBtn?.addEventListener('click', () => fileInput?.click())
  fileInput?.addEventListener('change', async () => {
    const file = fileInput.files?.[0]
    if (!file)
      return
    if (!file.name.toLowerCase().endsWith('.xlsx')) {
      alert('Only .xlsx import is supported')
      return
    }

    const resp = await importSheet(file)
    if (resp.redirected) {
      location.href = resp.url
      return
    }
    location.reload()
  })

  deleteBtn?.addEventListener('click', async () => {
    const ids = Array.from(document.querySelectorAll<HTMLInputElement>('.fileCheckbox:checked')).map(checkbox => checkbox.value)
    if (!ids.length) {
      alert('Please select files to delete')
      return
    }

    await deleteFiles(ids)
    location.reload()
  })

  wireInviteDialog(fileContainer, currentUserId)
}

export async function renderFilesPage() {
  renderFilesShell()
  const filesResp = await fetchFiles()

  const fileContainer = document.querySelector<HTMLDivElement>('#files-container')
  if (!fileContainer)
    return

  fileContainer.innerHTML = filesResp.files.map((file) => {
    const role = normalizeRole(file)
    return `
    <div class="file-row hover-effect" data-file-id="${file.id}">
      <input type="checkbox" name="fileIds" class="fileCheckbox file-check" value="${file.id}"/>
      <div class="file-name">
        ${file.openUrl ? `<a href="${file.openUrl}">${escapeHtml(file.name)}.xlsx</a>` : escapeHtml(file.name)}
      </div>
      <span class="file-role-badge role-${role}">${role}</span>
      <label class="file-updated">${escapeHtml(file.updatedAt)}</label>
      <div class="file-actions">
        <button class="demo-btn-secondary members-btn" type="button" data-unit-id="${file.unitId}">Members</button>
        <a class="file-action-link" href="${file.exportUrl}">Export</a>
        ${role === 'owner' ? `<button class="demo-btn-secondary invite-btn" type="button" data-file-id="${file.id}">Invite</button>` : ''}
      </div>
    </div>
  `
  }).join('')

  wireAvatar(filesResp.userId)
  wireFileRowToggle(fileContainer)
  wireCommonActions(fileContainer, filesResp.userId)

  fileContainer.querySelectorAll<HTMLButtonElement>('.members-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const unitId = btn.dataset.unitId ?? ''
      if (!unitId)
        return
      await openMembersDialog(unitId)
    })
  })
}

function normalizeRole(file: FileItem) {
  if (file.role === 'owner' || file.role === 'editor' || file.role === 'reader')
    return file.role
  return 'reader'
}
