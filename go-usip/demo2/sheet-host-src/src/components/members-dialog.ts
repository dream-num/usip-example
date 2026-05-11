import { fetchCollaborators } from '../services/collaborators-service'
import type { CollaboratorItem } from '../types/collaborators'
import { escapeHtml } from '../utils/html'

const DIALOG_ID = 'members-dialog'

function ensureDialog() {
  let dialog = document.getElementById(DIALOG_ID) as HTMLDialogElement | null
  if (dialog)
    return dialog

  dialog = document.createElement('dialog')
  dialog.id = DIALOG_ID
  dialog.className = 'dialog-panel members-dialog'
  dialog.innerHTML = `
    <div class="dialog-header">
      <h3>Members</h3>
      <p>People who can access this file.</p>
    </div>
    <div class="members-list" id="members-list"></div>
    <div class="dialog-actions">
      <button id="members-close-btn" class="demo-btn-secondary" type="button">Close</button>
    </div>
  `

  document.body.appendChild(dialog)

  const closeBtn = dialog.querySelector<HTMLButtonElement>('#members-close-btn')
  closeBtn?.addEventListener('click', () => dialog?.close())

  return dialog
}

function roleClass(role: string) {
  if (role === 'owner' || role === 'editor' || role === 'reader')
    return role
  return 'reader'
}

function renderMembers(listEl: HTMLDivElement, members: CollaboratorItem[]) {
  if (!members.length) {
    listEl.innerHTML = '<p class="members-empty">No members found.</p>'
    return
  }

  listEl.innerHTML = members.map((member) => {
    const role = roleClass(member.role)
    return `
      <div class="members-row">
        <div class="members-profile">
          <img class="members-avatar" src="${escapeHtml(member.subject.avatar)}" alt="avatar" loading="lazy" decoding="async" />
          <span class="members-name">${escapeHtml(member.subject.name)}</span>
        </div>
        <span class="file-role-badge role-${role}">${role}</span>
      </div>
    `
  }).join('')
}

export async function openMembersDialog(unitId: string) {
  const dialog = ensureDialog()
  const listEl = dialog.querySelector<HTMLDivElement>('#members-list')
  if (!listEl)
    return

  listEl.innerHTML = '<p class="members-empty">Loading...</p>'
  dialog.showModal()

  try {
    const members = await fetchCollaborators(unitId)
    renderMembers(listEl, members)
  }
  catch {
    listEl.innerHTML = '<p class="members-empty">Failed to load members.</p>'
  }
}
