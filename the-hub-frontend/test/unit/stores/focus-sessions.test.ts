import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref, computed } from 'vue'

// The focus-sessions store relies on Vue auto-imports (ref/computed) and Nuxt
// auto-imports (defineStore, useNuxtApp, useToast). Stub them all before the
// store module loads (defineStore is referenced at module top level, so the
// store must be imported dynamically below); the $api mock stands in for the
// API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('computed', computed)
vi.stubGlobal('defineStore', defineStore)

const { useFocusSessionStore } = await import('../../../app/stores/focus-sessions')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const session = (id: string, durationMin = 25, overrides: Record<string, unknown> = {}) => ({
  focus_session_id: id,
  user_id: 'u1',
  duration_min: durationMin,
  started_at: '2026-11-01T10:00:00.000Z',
  status: 'completed',
  session_type: 'focus',
  cycle_order: 0,
  notes: '',
  ...overrides
})

describe('useFocusSessionStore', () => {
  describe('startSession', () => {
    it('optimistically activates a session, then swaps in the server session', async () => {
      const store = useFocusSessionStore()
      const serverSession = session('s1')
      apiMock.mockResolvedValue(serverSession)

      const result = await store.startSession({ task_id: 't1', session_type: 'focus', planned_duration: 25 })

      expect(apiMock).toHaveBeenCalledWith('/focus-sessions', {
        method: 'POST',
        body: JSON.stringify({ task_id: 't1', session_type: 'focus', planned_duration: 25 })
      })
      // Sessions list keeps the confirmed session
      expect(store.sessions).toEqual([serverSession])
      expect(store.activeSession).toEqual(serverSession)
      expect(result).toEqual(serverSession)
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ color: 'success', description: 'Focus session started' })
      )
    })

    it('applies defaults for optional fields', async () => {
      const store = useFocusSessionStore()
      apiMock.mockResolvedValue(session('s1'))

      await store.startSession({})

      expect(apiMock).toHaveBeenCalledWith('/focus-sessions', {
        method: 'POST',
        body: JSON.stringify({})
      })
    })

    it('clears the active session on failure', async () => {
      const store = useFocusSessionStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.startSession({})

      expect(store.activeSession).toBeNull()
      expect(store.sessions).toHaveLength(0)
      expect(result).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ color: 'error', description: 'Failed to start focus session' })
      )
    })
  })

  describe('stopSession', () => {
    it('updates the session in the list and clears the active session', async () => {
      const store = useFocusSessionStore()
      const stopped = session('s1', 25, { status: 'completed' })
      apiMock.mockResolvedValue(stopped)
      store.sessions.push(session('s1', 0, { status: 'active' }))
      store.activeSession = session('s1', 0, { status: 'active' })

      const result = await store.stopSession('s1', 'completed')

      expect(apiMock).toHaveBeenCalledWith('/focus-sessions/s1', {
        method: 'PATCH',
        body: JSON.stringify({ status: 'completed' })
      })
      expect(store.activeSession).toBeNull()
      expect(store.sessions[0]).toEqual(stopped)
      expect(result).toEqual(stopped)
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ description: 'Focus session completed' })
      )
    })

    it('keeps a different session active when stopping another one', async () => {
      const store = useFocusSessionStore()
      apiMock.mockResolvedValue(session('s1'))
      const other = session('s2', 0, { status: 'active' })
      store.activeSession = other

      await store.stopSession('s1')

      expect(store.activeSession).toEqual(other)
    })

    it('shows an error toast on failure and leaves state unchanged', async () => {
      const store = useFocusSessionStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.stopSession('s1')

      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ color: 'error', description: 'Failed to stop session' })
      )
    })
  })

  describe('fetchSessions', () => {
    it('builds query params from provided filters and sets sessions', async () => {
      const store = useFocusSessionStore()
      const serverSessions = [session('s1'), session('s2')]
      apiMock.mockResolvedValue({ focus_sessions: serverSessions })

      await store.fetchSessions({ task_id: 't1', status: 'completed', limit: 10 })

      expect(apiMock).toHaveBeenCalledWith(
        `/focus-sessions?${new URLSearchParams({ task_id: 't1', status: 'completed', limit: '10' })}`
      )
      expect(store.sessions).toEqual(serverSessions)
      expect(store.loading).toBe(false)
    })

    it('requests a bare path when no filters are given', async () => {
      const store = useFocusSessionStore()
      apiMock.mockResolvedValue({ focus_sessions: [] })

      await store.fetchSessions()

      expect(apiMock).toHaveBeenCalledWith('/focus-sessions?')
      expect(store.sessions).toEqual([])
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useFocusSessionStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchSessions()
      expect(store.loading).toBe(true)

      resolveApi({ focus_sessions: [] })
      await pending
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchStats', () => {
    it('stores the stats payload with the requested range', async () => {
      const store = useFocusSessionStore()
      const statsPayload = { total_minutes: 120, daily_stats: [], weekly_stats: [] }
      apiMock.mockResolvedValue(statsPayload)

      await store.fetchStats(7, 'break')

      expect(apiMock).toHaveBeenCalledWith('/focus-sessions/stats?days=7&session_type=break')
      expect(store.stats).toEqual(statsPayload)
      expect(store.loading).toBe(false)
    })
  })

  describe('computed stats', () => {
    it('getTodayMinutes returns 0 when there are no stats', () => {
      const store = useFocusSessionStore()

      expect(store.getTodayMinutes).toBe(0)
    })

    it('getTodayMinutes returns today\'s minutes from daily stats', () => {
      const store = useFocusSessionStore()
      const today = new Date().toISOString().split('T')[0]
      store.stats = {
        total_minutes: 0,
        total_hours: 0,
        total_sessions: 0,
        days: 0,
        daily_stats: [
          { date: '2020-01-01', minutes: 15, sessions_count: 1 },
          { date: today, minutes: 45, sessions_count: 2 }
        ],
        task_breakdown: [],
        weekly_stats: [],
        average_daily: 0
      } as any

      expect(store.getTodayMinutes).toBe(45)
    })

    it('getWeeklyMinutes sums minutes across weekly stats', () => {
      const store = useFocusSessionStore()
      store.stats = {
        daily_stats: [],
        weekly_stats: [
          { week_start: '2026-10-27', minutes: 100, sessions_count: 4 },
          { week_start: '2026-11-03', minutes: 80, sessions_count: 3 }
        ]
      } as any

      expect(store.getWeeklyMinutes).toBe(180)
    })
  })
})
