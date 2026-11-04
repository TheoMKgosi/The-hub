import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useCalendarStore } from '../../../app/stores/calendar'

// The calendar store reaches for the Nuxt auto-imports useNuxtApp (no toast
// usage). It also redirects via window.location in initiateGoogleAuth, so a
// window stub is provided; the $api mock stands in for the API plugin.
const apiMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  // Fresh window per test so redirects do not leak between tests
  vi.stubGlobal('window', { location: { href: '' } })
})

const integration = (id: string, overrides: Record<string, unknown> = {}) => ({
  id,
  provider: 'google' as const,
  provider_user_id: 'prov-1',
  calendar_id: 'cal-1',
  is_active: true,
  sync_enabled: true,
  last_sync_at: undefined as string | undefined,
  ...overrides
})

describe('useCalendarStore', () => {
  it('starts with no integrations, no error and no loading', () => {
    const store = useCalendarStore()

    expect(store.integrations).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.fetchError).toBeNull()
  })

  describe('fetchIntegrations', () => {
    it('populates integrations, clears any previous error and clears loading', async () => {
      const store = useCalendarStore()
      const serverIntegrations = [integration('i1'), integration('i2')]
      apiMock.mockResolvedValue({ integrations: serverIntegrations })

      await store.fetchIntegrations()

      expect(apiMock).toHaveBeenCalledWith('/calendar/integrations')
      expect(store.integrations).toEqual(serverIntegrations)
      expect(store.fetchError).toBeNull()
      expect(store.loading).toBe(false)
    })

    it('falls back to an empty list when the response has no integrations', async () => {
      const store = useCalendarStore()
      apiMock.mockResolvedValue({})

      await store.fetchIntegrations()

      expect(store.integrations).toEqual([])
      expect(store.loading).toBe(false)
    })

    it('records the error and clears loading on failure', async () => {
      const store = useCalendarStore()
      const failure = new Error('server error')
      apiMock.mockRejectedValue(failure)

      await store.fetchIntegrations()

      expect(store.fetchError).toBe(failure)
      expect(store.loading).toBe(false)
    })
  })

  describe('initiateGoogleAuth', () => {
    it('redirects to the OAuth url returned by the API', async () => {
      const store = useCalendarStore()
      apiMock.mockResolvedValue({ auth_url: 'https://accounts.google.com/oauth' })

      await store.initiateGoogleAuth()

      expect(apiMock).toHaveBeenCalledWith('/calendar/google/auth', { method: 'POST' })
      expect(window.location.href).toBe('https://accounts.google.com/oauth')
      expect(store.fetchError).toBeNull()
    })

    it('records the error and does not redirect on failure', async () => {
      const store = useCalendarStore()
      const failure = new Error('oauth down')
      apiMock.mockRejectedValue(failure)

      await store.initiateGoogleAuth()

      expect(store.fetchError).toBe(failure)
      expect(window.location.href).toBe('')
    })
  })

  describe('syncCalendar', () => {
    it('optimistically sets last_sync_at, then reconciles with the server response', async () => {
      const store = useCalendarStore()
      store.integrations.push(integration('i1'))
      const serverIntegration = integration('i1', { last_sync_at: '2026-11-02T09:00:00.000Z' })
      apiMock.mockResolvedValue(serverIntegration)

      await store.syncCalendar('i1')

      expect(apiMock).toHaveBeenCalledWith('/calendar/integrations/i1/sync', { method: 'POST' })
      expect(store.integrations[0]).toEqual(serverIntegration)
      expect(store.fetchError).toBeNull()
    })

    it('reverts last_sync_at on failure', async () => {
      const store = useCalendarStore()
      store.integrations.push(integration('i1'))
      apiMock.mockRejectedValue(new Error('sync failed'))

      await store.syncCalendar('i1')

      expect(store.integrations[0].last_sync_at).toBeUndefined()
      expect(store.fetchError).toBeInstanceOf(Error)
    })

    it('leaves state untouched for an unknown integration id', async () => {
      const store = useCalendarStore()
      store.integrations.push(integration('i1'))
      apiMock.mockResolvedValue(integration('i2'))

      await store.syncCalendar('i2')

      expect(store.integrations).toHaveLength(1)
      expect(store.integrations[0].id).toBe('i1')
    })
  })

  describe('deleteIntegration', () => {
    it('removes the integration optimistically and succeeds', async () => {
      const store = useCalendarStore()
      store.integrations.push(integration('i1'), integration('i2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteIntegration('i1')

      expect(apiMock).toHaveBeenCalledWith('/calendar/integrations/i1', { method: 'DELETE' })
      expect(store.integrations.map(i => i.id)).toEqual(['i2'])
      expect(store.fetchError).toBeNull()
    })

    it('sets an error and does not call the API for an unknown id', async () => {
      const store = useCalendarStore()

      await store.deleteIntegration('missing')

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError).toBeInstanceOf(Error)
    })

    it('restores the integration when deletion fails', async () => {
      const store = useCalendarStore()
      store.integrations.push(integration('i1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteIntegration('i1')

      expect(store.integrations.map(i => i.id)).toEqual(['i1'])
      expect(store.fetchError).toBeInstanceOf(Error)
    })
  })

  it('reset clears all integrations', () => {
    const store = useCalendarStore()
    store.integrations.push(integration('i1'))

    store.reset()

    expect(store.integrations).toHaveLength(0)
  })
})
