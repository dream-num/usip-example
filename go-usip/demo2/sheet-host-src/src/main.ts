import './style.css'
import { renderFilesPage } from './pages/files-page'
import { renderLoginPage } from './pages/login-page'
import { renderRegisterPage } from './pages/register-page'
import { renderSheetPage } from './pages/sheet-page'

async function main() {
  const path = window.location.pathname

  if (path === '/login') {
    renderLoginPage()
    return
  }

  if (path === '/register') {
    renderRegisterPage()
    return
  }

  if (path === '/files') {
    await renderFilesPage()
    return
  }

  renderSheetPage()
}

main().catch((error) => {
  // eslint-disable-next-line no-console
  console.error(error)
})
