import './style.css'
import { setupUniver } from './setup-univer'

function main() {
  const univerAPI = setupUniver()

  // test on dev
  window.univerAPI = univerAPI

  // Toolbar debug actions are intentionally disabled in demo2 runtime UI.
  // Keep only core sheet rendering behavior.
  void univerAPI
}

main()
