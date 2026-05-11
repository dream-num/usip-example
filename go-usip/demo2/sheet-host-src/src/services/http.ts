import type { APIError } from '../types/api'

export async function apiFetch<T>(input: RequestInfo | URL, init?: RequestInit): Promise<T> {
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
