import { fetchPeople, inviteUsers } from '../services/files-service'

export function wireInviteDialog(fileContainer: HTMLDivElement) {
  const dialog = document.querySelector<HTMLDialogElement>('#dialog')
  const dialogMsg = document.querySelector<HTMLDivElement>('#dialog-msg')
  const roleSelect = document.querySelector<HTMLSelectElement>('#select-role')
  const okBtn = document.querySelector<HTMLButtonElement>('#dialog-ok-btn')
  const cancelInviteBtn = document.querySelector<HTMLButtonElement>('#dialog-cancel-btn')

  if (!dialog || !dialogMsg || !roleSelect || !okBtn || !cancelInviteBtn)
    return

  let inviteFileId = 0
  let invite: string[] = []
  const pageStack: number[] = [0]

  async function renderPeople(next: number) {
    const data = await fetchPeople(next)
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

    if (pageStack.length > 1 || data.next !== 0) {
      const pagination = document.createElement('div')
      pagination.className = 'dialog-pagination'

      if (pageStack.length > 1) {
        const prevBtn = document.createElement('button')
        prevBtn.type = 'button'
        prevBtn.className = 'dialog-page-btn'
        prevBtn.innerText = 'Prev'
        prevBtn.onclick = async () => {
          pageStack.pop()
          const prevNext = pageStack[pageStack.length - 1]
          await renderPeople(prevNext)
        }
        pagination.appendChild(prevBtn)
      }

      if (data.next !== 0) {
        const nextBtn = document.createElement('button')
        nextBtn.type = 'button'
        nextBtn.className = 'dialog-page-btn'
        nextBtn.innerText = 'Next'
        nextBtn.onclick = async () => {
          pageStack.push(data.next)
          await renderPeople(data.next)
        }
        pagination.appendChild(nextBtn)
      }

      dialogMsg.appendChild(pagination)
    }
  }

  fileContainer.querySelectorAll<HTMLButtonElement>('.invite-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
      invite = []
      inviteFileId = Number(btn.dataset.fileId ?? 0)
      pageStack.length = 0
      pageStack.push(0)
      await renderPeople(0)
      dialog.showModal()
    })
  })

  okBtn.addEventListener('click', async () => {
    if (!inviteFileId)
      return

    await inviteUsers(inviteFileId, invite, roleSelect.value)
    dialog.close()
  })

  cancelInviteBtn.addEventListener('click', () => dialog.close())
}
