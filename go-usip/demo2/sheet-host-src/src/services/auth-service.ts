import type { AuthResp } from '../types/auth'
import { apiFetch } from './http'

export async function login(username: string, password: string) {
  return apiFetch<AuthResp>('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
}

export async function register(nickname: string, username: string, password: string) {
  return apiFetch<AuthResp>('/api/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname, username, password }),
  })
}

export async function logout() {
  return apiFetch<{ ok: boolean }>('/api/auth/logout', { method: 'POST' })
}
