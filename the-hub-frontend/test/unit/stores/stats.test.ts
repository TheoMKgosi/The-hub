import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref, computed } from 'vue'

// The stats store relies on Vue auto-imports (ref, computed) and Nuxt
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

const { useStatsStore } = await import('../../../app/stores/stats')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const stats = (overrides: Record<string, unknown> = {}) => ({
  period: { start_date: '2026-11-01', end_date: '2026-11-02' },
  summary: {
    total_tasks: 10,
    completed_tasks: 6,
    pending_tasks: 3,
    overdue_tasks: 1,
    completion_rate: 60
  },
  priority_distribution: { '1': 2, '2': 2, '3': 2, '4': 2, '5': 2 },
  priority_completion: { '1': 1, '2': 1, '3': 1, '4': 1, '5': 1 },
  goal_distribution: { with_goals: 7, without_goals: 3 },
  time_based: { due_today: 2, due_tomorrow: 1, due_this_week: 5 },
  ...overrides
}) as any

const trend = (date: string, completionRate = 50) => ({
  date,
  period: { start_date: date, end_date: date },
  summary: {
    total_tasks: 10,
    completed_tasks: 5,
    pending_tasks: 4,
    overdue_tasks: 1,
    completion_rate: completionRate
  },
  priority_distribution: { '1': 2, '2': 2, '3': 2, '4': 2, '5': 2 },
  priority_completion: { '1': 1, '2': 1, '3': 1, '4': 1, '5': 1 },
  goal_distribution: { with_goals: 7, without_goals: 3 },
  time_based: { due_today: 2, due_tomorrow: 1, due_this_week: 5 }
}) as any

describe('useStatsStore', () => {
  describe('initial state', () => {
    it('starts empty with no error and no loading', () => {
      const store = useStatsStore()

      expect(store.currentStats).toBeNull()
      expect(store.trends).toEqual([])
      expect(store.loading).toBe(false)
      expect(store.trendsLoading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('computed metrics fall back to zero when there are no stats', () => {
      const store = useStatsStore()

      expect(store.completionRate).toBe(0)
      expect(store.totalTasks).toBe(0)
      expect(store.completedTasks).toBe(0)
      expect(store.pendingTasks).toBe(0)
      expect(store.overdueTasks).toBe(0)
      expect(store.tasksDueToday).toBe(0)
      expect(store.tasksDueTomorrow).toBe(0)
      expect(store.tasksDueThisWeek).toBe(0)
      expect(store.priorityData).toEqual([])
      expect(store.goalData).toEqual([])
    })
  })

  describe('fetchStats', () => {
    it('requests today by default and stores the response', async () => {
      const store = useStatsStore()
      const payload = stats()
      apiMock.mockResolvedValue({ stats: payload })

      await store.fetchStats()

      expect(apiMock).toHaveBeenCalledWith('/stats/tasks?period=today')
      expect(store.currentStats).toEqual(payload)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(toastAddMock).not.toHaveBeenCalled()
    })

    it('builds query params from the given period and date range', async () => {
      const store = useStatsStore()
      apiMock.mockResolvedValue({ stats: stats() })

      await store.fetchStats('week', '2026-10-01', '2026-10-31')

      expect(apiMock).toHaveBeenCalledWith('/stats/tasks?period=week&start_date=2026-10-01&end_date=2026-10-31')
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useStatsStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchStats()
      expect(store.loading).toBe(true)

      resolveApi({ stats: stats() })
      await pending
      expect(store.loading).toBe(false)
    })

    it('keeps existing stats when the response has none', async () => {
      const store = useStatsStore()
      store.currentStats = stats()
      apiMock.mockResolvedValue({})

      await store.fetchStats()

      expect(store.currentStats).toEqual(stats())
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('sets the error, shows a toast, and clears loading on failure', async () => {
      const store = useStatsStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchStats()

      expect(store.error).toBe('Failed to fetch task statistics')
      expect(store.currentStats).toBeNull()
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load statistics', color: 'error' })
    })
  })

  describe('fetchTrends', () => {
    it('requests 30 days by default and stores the trends', async () => {
      const store = useStatsStore()
      const trends = [trend('2026-11-01'), trend('2026-11-02')]
      apiMock.mockResolvedValue({ trends })

      await store.fetchTrends()

      expect(apiMock).toHaveBeenCalledWith('/stats/tasks/trends?days=30')
      expect(store.trends).toEqual(trends)
      expect(store.trendsLoading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('requests the given number of days', async () => {
      const store = useStatsStore()
      apiMock.mockResolvedValue({ trends: [] })

      await store.fetchTrends(7)

      expect(apiMock).toHaveBeenCalledWith('/stats/tasks/trends?days=7')
    })

    it('toggles trendsLoading while the request is in flight', async () => {
      const store = useStatsStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchTrends()
      expect(store.trendsLoading).toBe(true)

      resolveApi({ trends: [] })
      await pending
      expect(store.trendsLoading).toBe(false)
    })

    it('sets the error, shows a toast, and clears trendsLoading on failure', async () => {
      const store = useStatsStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchTrends()

      expect(store.error).toBe('Failed to fetch statistics trends')
      expect(store.trendsLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to load trends', color: 'error' })
    })
  })

  describe('computed metrics', () => {
    it('exposes summary metrics from the current stats', () => {
      const store = useStatsStore()
      store.currentStats = stats()

      expect(store.completionRate).toBe(60)
      expect(store.totalTasks).toBe(10)
      expect(store.completedTasks).toBe(6)
      expect(store.pendingTasks).toBe(3)
      expect(store.overdueTasks).toBe(1)
      expect(store.tasksDueToday).toBe(2)
      expect(store.tasksDueTomorrow).toBe(1)
      expect(store.tasksDueThisWeek).toBe(5)
    })

    it('maps the priority distribution into chart-ready data', () => {
      const store = useStatsStore()
      store.currentStats = stats()

      expect(store.priorityData).toEqual([
        { priority: '1 - Low', count: 2, completed: 1 },
        { priority: '2 - Low', count: 2, completed: 1 },
        { priority: '3 - Medium', count: 2, completed: 1 },
        { priority: '4 - High', count: 2, completed: 1 },
        { priority: '5 - High', count: 2, completed: 1 }
      ])
    })

    it('maps the goal distribution into chart-ready data', () => {
      const store = useStatsStore()
      store.currentStats = stats()

      expect(store.goalData).toEqual([
        { name: 'With Goals', value: 7 },
        { name: 'Without Goals', value: 3 }
      ])
    })

    it('maps trends into time-series data', () => {
      const store = useStatsStore()
      store.trends = [trend('2026-11-01', 40), trend('2026-11-02', 80)]

      expect(store.trendsCompletionData).toEqual([
        { date: '2026-11-01', completionRate: 40, completedTasks: 5, totalTasks: 10 },
        { date: '2026-11-02', completionRate: 80, completedTasks: 5, totalTasks: 10 }
      ])
    })
  })

  describe('reset', () => {
    it('clears stats, trends, and loading/error state', async () => {
      const store = useStatsStore()
      apiMock.mockResolvedValue({ stats: stats(), trends: [] })
      await store.fetchStats()
      await store.fetchTrends()

      store.reset()

      expect(store.currentStats).toBeNull()
      expect(store.trends).toEqual([])
      expect(store.loading).toBe(false)
      expect(store.trendsLoading).toBe(false)
      expect(store.error).toBeNull()
    })
  })
})
