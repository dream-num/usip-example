import { setupUniver } from '../setup-univer'
import { openMembersDialog } from '../components/members-dialog'

export function renderSheetPage() {
  const app = document.querySelector<HTMLDivElement>('#app')
  if (!app)
    return

  app.innerHTML = `
    <div class="sheet-shell">
      <div class="sheet-toolbar">
        <div class="sheet-toolbar-left">Sheet</div>
        <button id="sheet-members-btn" class="demo-btn-secondary" type="button">Members</button>
      </div>
      <div id="univer"></div>
    </div>
  `
  const univerAPI = setupUniver()
  window.univerAPI = univerAPI

  const btn = document.querySelector<HTMLButtonElement>('#sheet-members-btn')
  btn?.addEventListener('click', async () => {
    const unitId = new URLSearchParams(window.location.search).get('unit') ?? ''
    if (!unitId) {
      alert('Current sheet has no unit id')
      return
    }
    await openMembersDialog(unitId)
  })
}
