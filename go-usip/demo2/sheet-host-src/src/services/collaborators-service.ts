import type { CollaboratorsResp } from '../types/collaborators'
import { apiFetch } from './http'

export async function fetchCollaborators(unitId: string) {
  const payload = await apiFetch<CollaboratorsResp>('/usip/collaborators', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ unitIDs: [unitId] }),
  })

  const target = payload.collaborators.find(item => item.unitID === unitId)
  return target?.subjects ?? []
}
