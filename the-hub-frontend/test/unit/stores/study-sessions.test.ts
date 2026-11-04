import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import { useStudySessionStore } from '../../../app/stores/study-sessions'

// The study-sessions store reaches for Vue auto-imports (ref) and Nuxt
// auto-imports (useNuxtApp, useToast). Stub them all; the $api mock stands in
// for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('ref', ref)
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const session = (id: string, durationMin = 25, overrides: Record<string, unknown> = {}) => ({
  id,
  user_id: 'u1',
  topic_id: 'topic-1',
  task_id: undefined,
  duration_min: durationMin,
  started_at: '2026-11-01T10:00:00.000Z',
  ended_at: '2026-11-01T10:25:00.000Z',
  ...overrides
})

const stats = (overrides: Record<string, unknown> = {}) => ({
  total_minutes: 120,
  total_hours: 2,
  days: 30,
  daily_stats: [{ date: '2026-11-01', minutes: 40 }],
  topic_stats: [{ topic_id: 'topic-1', topic_title: 'Math', minutes: 90, sessions: 3 }],
  average_daily: 4,
  ...overrides
})

describe('useStudySessionStore', () => {
  it('starts with empty sessions, no stats and no loading', () => {
    const store = useStudySessionStore()

    expect(store.sessions).toEqual([])
    expect(store.stats).toBeNull()
    expect(store.loading).toBe(false)
  })

  describe('createSession', () => {
    it('adds the session optimistically, then replaces it with the server response', async () => {
      const store = useStudySessionStore()
      const serverSession = session('s1')
      apiMock.mockResolvedValue(serverSession)

      const result = await store.createSession({ topic_id: 'topic-1', duration_min: 25 })

      expect(apiMock).toHaveBeenCalledWith('/study-sessions', {
        method: 'POST',
        body: JSON.stringify({ topic_id: 'topic-1', duration_min: 25 })
      })
      // Optimistic temp entry is swapped for the confirmed session
      expect(store.sessions).toHaveLength(1)
      expect(store.sessions[0]).toEqual(serverSession)
      expect(result).toEqual(serverSession)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Learning',
        description: 'Study session logged successfully',
        color: 'success'
      })
    })

    it('rolls back the optimistic session on failure', async () => {
      const store = useStudySessionStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.createSession({ topic_id: 'topic-1', duration_min: 25 })

      expect(store.sessions).toHaveLength(0)
      expect(result).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to log study session',
        color: 'error'
      })
    })
  })

  describe('fetchSessions', () => {
    it('builds query params from provided filters and sets sessions', async () => {
      const store = useStudySessionStore()
      const serverSessions = [session('s1'), session('s2')]
      apiMock.mockResolvedValue({ study_sessions: serverSessions })

      await store.fetchSessions({ topic_id: 'topic-1', date_from: '2026-11-01', limit: 10 })

      expect(apiMock).toHaveBeenCalledWith(
        `/study-sessions?${new URLSearchParams({ topic_id: 'topic-1', date_from: '2026-11-01', limit: '10' })}`
      )
      expect(store.sessions).toEqual(serverSessions)
      expect(store.loading).toBe(false)
    })

    it('requests a bare path when no filters are given', async () => {
      const store = useStudySessionStore()
      apiMock.mockResolvedValue({ study_sessions: [] })

      await store.fetchSessions()

      expect(apiMock).toHaveBeenCalledWith('/study-sessions?')
      expect(store.sessions).toEqual([])
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useStudySessionStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchSessions()
      expect(store.loading).toBe(true)

      resolveApi({ study_sessions: [] })
      await pending
      expect(store.loading).toBe(false)
    })

    it('shows an error toast on failure and clears loading', async () => {
      const store = useStudySessionStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchSessions()

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch study sessions',
        color: 'error'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchStats', () => {
    it('defaults to a 30 day range and stores the stats payload', async () => {
      const store = useStudySessionStore()
      const statsPayload = stats()
      apiMock.mockResolvedValue(statsPayload)

      await store.fetchStats()

      expect(apiMock).toHaveBeenCalledWith('/study-sessions/stats?days=30')
      expect(store.stats).toEqual(statsPayload)
      expect(store.loading).toBe(false)
    })

    it('honours a custom range', async () => {
      const store = useStudySessionStore()
      apiMock.mockResolvedValue(stats())

      await store.fetchStats(7)

      expect(apiMock).toHaveBeenCalledWith('/study-sessions/stats?days=7')
    })

    it('shows an error toast on failure', async () => {
      const store = useStudySessionStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchStats(7)

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch study statistics',
        color: 'error'
      })
      expect(store.stats).toBeNull()
      expect(store.loading).toBe(false)
    })
  })

  describe('stat getters', () => {
    it('return zero/empty defaults when no stats are loaded', () => {
      const store = useStudySessionStore()

      expect(store.getTotalStudyTime()).toBe(0)
      expect(store.getAverageDailyTime()).toBe(0)
      expect(store.getTopicBreakdown()).toEqual([])
      expect(store.getDailyStats()).toEqual([])
    })

    it('read the loaded stats payload', async () => {
      const store = useStudySessionStore()
      apiMock.mockResolvedValue(stats())
      await store.fetchStats()

      expect(store.getTotalStudyTime()).toBe(120)
      expect(store.getAverageDailyTime()).toBe(4)
      expect(store.getTopicBreakdown()).toEqual([
        { topic_id: 'topic-1', topic_title: 'Math', minutes: 90, sessions: 3 }
      ])
      expect(store.getDailyStats()).toEqual([{ date: '2026-11-01', minutes: 40 }])
    })
  })
})
