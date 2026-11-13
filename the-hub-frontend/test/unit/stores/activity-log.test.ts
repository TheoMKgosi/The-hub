import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref, computed } from 'vue'

// The store references Nuxt auto-imports at module top level, so import it
// dynamically after the stubs are in place.
let useActivityLogStore: typeof import('../../../app/stores/activity-log')['useActivityLogStore']

// The store reaches for Nuxt auto-imports (useNuxtApp, useToast, ref, computed).
// Stub them; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(async () => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('defineStore', defineStore)
  vi.stubGlobal('ref', ref)
  vi.stubGlobal('computed', computed)

  ;({ useActivityLogStore } = await import('../../../app/stores/activity-log'))
})

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const makeEntry = (overrides: Partial<Record<string, unknown>> = {}) => ({
  activity_id: 'a1',
  user_id: 'u1',
  title: 'Worked on thesis',
  notes: '',
  category: 'Work',
  tags: ['focus'],
  activity_date: '2026-11-04T00:00:00Z',
  duration_min: 90,
  created_at: '2026-11-04T10:00:00.000Z',
  updated_at: '2026-11-04T10:00:00.000Z',
  ...overrides
})

describe('useActivityLogStore', () => {
  describe('fetchEntries', () => {
    it('populates entries and total from the API', async () => {
      const store = useActivityLogStore()
      apiMock.mockResolvedValue({ activities: [makeEntry()], total: 1 })

      await store.fetchEntries()

      expect(apiMock).toHaveBeenCalledWith('/activities?')
      expect(store.entries).toEqual([makeEntry()])
      expect(store.total).toBe(1)
      expect(store.loading).toBe(false)
    })

    it('sends filters as query params', async () => {
      const store = useActivityLogStore()
      apiMock.mockResolvedValue({ activities: [], total: 0 })

      await store.fetchEntries({ category: 'Work', date_from: '2026-11-01', limit: 50 })

      expect(apiMock).toHaveBeenCalledWith('/activities?date_from=2026-11-01&category=Work&limit=50')
    })

    it('shows an error toast on failure', async () => {
      const store = useActivityLogStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.fetchEntries()

      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to fetch activities', color: 'error' })
      expect(store.loading).toBe(false)
    })
  })

  describe('createEntry', () => {
    it('posts the entry and prepends it to the list', async () => {
      const store = useActivityLogStore()
      const created = makeEntry({ activity_id: 'a2' })
      apiMock.mockResolvedValue(created)

      const result = await store.createEntry({ title: 'Read a book', duration_min: 30 })

      expect(apiMock).toHaveBeenCalledWith('/activities', {
        method: 'POST',
        body: JSON.stringify({ title: 'Read a book', duration_min: 30 })
      })
      expect(result).toEqual(created)
      expect(store.entries[0]).toEqual(created)
      expect(store.total).toBe(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Activity', description: 'Activity logged', color: 'success' })
    })

    it('does not add the entry when the API fails', async () => {
      const store = useActivityLogStore()
      apiMock.mockRejectedValue(new Error('bad request'))

      const result = await store.createEntry({ title: 'x' })

      expect(result).toBeUndefined()
      expect(store.entries).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Error', description: 'Failed to log activity', color: 'error' })
    })
  })

  describe('updateEntry', () => {
    it('patches the entry and replaces it in place', async () => {
      const store = useActivityLogStore()
      store.entries.push(makeEntry() as never)
      const updated = makeEntry({ title: 'Updated title' })
      apiMock.mockResolvedValue(updated)

      const result = await store.updateEntry('a1', { title: 'Updated title' })

      expect(apiMock).toHaveBeenCalledWith('/activities/a1', {
        method: 'PATCH',
        body: JSON.stringify({ title: 'Updated title' })
      })
      expect(result).toEqual(updated)
      expect(store.entries[0].title).toBe('Updated title')
    })
  })

  describe('deleteEntry', () => {
    it('deletes via the API and removes it from the list', async () => {
      const store = useActivityLogStore()
      store.entries.push(makeEntry() as never, makeEntry({ activity_id: 'a2' }) as never)
      store.total = 2
      apiMock.mockResolvedValue({ message: 'Activity deleted' })

      await store.deleteEntry('a1')

      expect(apiMock).toHaveBeenCalledWith('/activities/a1', { method: 'DELETE' })
      expect(store.entries.map(e => e.activity_id)).toEqual(['a2'])
      expect(store.total).toBe(1)
    })
  })

  describe('fetchStats', () => {
    it('stores the stats payload', async () => {
      const store = useActivityLogStore()
      const stats = {
        days: 30,
        total_entries: 5,
        total_minutes: 300,
        total_hours: 5,
        categories: [{ category: 'Work', minutes: 300, entries: 5 }],
        daily_stats: [],
        weekly_stats: [],
        top_tags: [],
        current_streak: 2
      }
      apiMock.mockResolvedValue(stats)

      await store.fetchStats(30)

      expect(apiMock).toHaveBeenCalledWith('/activities/stats?days=30')
      expect(store.stats).toEqual(stats)
    })
  })

  describe('suggestions', () => {
    it('derives unique category and tag suggestions from entries', () => {
      const store = useActivityLogStore()
      store.entries.push(
        makeEntry({ category: 'Work', tags: ['focus', 'writing'] }) as never,
        makeEntry({ activity_id: 'a2', category: 'Work', tags: ['focus'] }) as never,
        makeEntry({ activity_id: 'a3', category: '', tags: [] }) as never
      )

      expect(store.categorySuggestions).toEqual(['Work'])
      expect(store.tagSuggestions).toEqual(['focus', 'writing'])
    })
  })
})
