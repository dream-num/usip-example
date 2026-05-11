import type { FilesResp, UserListResp } from '../types/files'
import { apiFetch } from './http'

export async function fetchFiles() {
  return apiFetch<FilesResp>('/api/files')
}

export async function fetchPeople(next = 0) {
  return apiFetch<UserListResp>(`/user/people?next=${next}&size=10`)
}

export async function createSheet(name: string) {
  const formData = new FormData()
  formData.set('name', name)
  formData.set('type', 'sheet')
  return fetch('/file/new', { method: 'POST', body: formData })
}

export async function importSheet(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('type', 'sheet')
  formData.append('name', file.name.replace(/\.xlsx$/i, ''))
  return fetch('/file/import', { method: 'POST', body: formData })
}

export async function deleteFiles(fileIds: string[]) {
  const params = new URLSearchParams()
  fileIds.forEach(id => params.append('fileIds', id))
  return fetch(`/file?${params.toString()}`, { method: 'DELETE' })
}

export async function inviteUsers(fileId: number, userIds: string[], role: string) {
  return fetch('/file/join', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ fileId, userIds, role }),
  })
}
