import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useScheduleStore } from '../../../app/stores/schedule'

// The schedule store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them both; the $api mock stands in for the API plugin. The store
// persists suggestions to the bare `localStorage` global (which Node itself
// provides in this environment), so a Map-backed implementation is stubbed
// over it.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

const storageMap = new Map<string, string>()
const localStorageMock = {
  getItem: (key: string) => (storageMap.has(key) ? storageMap.get(key)! : null),
  setItem: (key: string, value: string) => { storageMap.set(key, String(value)) },
  removeItem: (key: string) => { storageMap.delete(key) }
}

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('window', { localStorage: localStorageMock })
  vi.stubGlobal('localStorage', localStorageMock)
})

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  storageMap.clear()
})

const event = (id: string, overrides: Record<string, unknown> = {}) => ({
  id,
  title: `Event ${id}`,
  start: new Date('2030-01-01T10:00:00.000Z'),
  end: new Date('2030-01-01T11:00:00.000Z'),
  ...overrides
}) as any

const serverEvent = (id: string, overrides: Record<string, unknown> = {}) => ({
  id,
  title: `Event ${id}`,
  start: '2030-01-01T10:00:00.000Z',
  end: '2030-01-01T11:00:00.000Z',
  ...overrides
})

const suggestion = (id: string, overrides: Record<string, unknown> = {}) => ({
  id,
  title: `Suggestion ${id}`,
  start: '2030-01-01T10:00:00.000Z',
  end: '2030-01-01T11:00:00.000Z',
  ...overrides
})

describe('useScheduleStore', () => {
  describe('fetchSchedule', () => {
    it('populates the schedule with Date objects and clears loading', async () => {
      const store = useScheduleStore()
      apiMock.mockResolvedValue({ schedule: [serverEvent('e1')] })

      await store.fetchSchedule()

      expect(apiMock).toHaveBeenCalledWith('/schedule')
      expect(store.schedule).toEqual([event('e1')])
      expect(store.schedule[0].start).toBeInstanceOf(Date)
      expect(store.loading).toBe(false)
      expect(store.fetchError).toBeNull()
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useScheduleStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchSchedule()
      expect(store.loading).toBe(true)

      resolveApi({ schedule: [] })
      await pending
      expect(store.loading).toBe(false)
    })

    it('records the error without toasting when the fetch fails', async () => {
      // NOTE (actual behavior): fetch failures only set fetchError; no toast.
      const store = useScheduleStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchSchedule()

      expect(store.fetchError).toBeInstanceOf(Error)
      expect(store.schedule).toEqual([])
      expect(store.loading).toBe(false)
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('submitForm', () => {
    it('adds the item optimistically, then replaces it with the server response', async () => {
      const store = useScheduleStore()
      const formData = {
        title: 'Standup',
        start: '2030-01-01T10:00:00.000Z',
        end: '2030-01-01T11:00:00.000Z'
      }
      apiMock.mockResolvedValue(serverEvent('e1', { title: 'Standup' }))

      await store.submitForm(formData)

      expect(apiMock).toHaveBeenCalledWith('schedule', {
        method: 'POST',
        body: JSON.stringify(formData)
      })
      expect(store.schedule).toHaveLength(1)
      expect(store.schedule[0]).toEqual(event('e1', { title: 'Standup' }))
      expect(store.schedule[0].start).toBeInstanceOf(Date)
      expect(store.fetchError).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Calendar',
        description: 'Schedule Added successfully',
        color: 'success'
      })
    })

    it('rolls back the optimistic item and records the error on failure', async () => {
      const store = useScheduleStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({
        title: 'Standup',
        start: '2030-01-01T10:00:00.000Z',
        end: '2030-01-01T11:00:00.000Z'
      })

      expect(store.schedule).toHaveLength(0)
      expect(store.fetchError).toBeInstanceOf(Error)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Schedule not added',
        color: 'error'
      })
    })
  })

  describe('updateSchedule', () => {
    it('returns false and sets fetchError for an unknown id', async () => {
      const store = useScheduleStore()

      const result = await store.updateSchedule('e-missing', { title: 'x' })

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('Schedule item not found')
    })

    it('rejects a start time in the past without calling the API', async () => {
      const store = useScheduleStore()
      const start = new Date(Date.now() + 60 * 60 * 1000)
      const end = new Date(Date.now() + 2 * 60 * 60 * 1000)
      store.schedule.push(event('e1', { start, end }))

      const result = await store.updateSchedule('e1', {
        start: new Date(Date.now() - 2 * 60 * 60 * 1000)
      })

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('Cannot schedule events in the past')
    })

    it('rejects an event longer than 24 hours without calling the API', async () => {
      const store = useScheduleStore()
      store.schedule.push(event('e1'))

      const result = await store.updateSchedule('e1', {
        end: new Date('2030-01-03T12:00:00.000Z')
      })

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('Event duration cannot exceed 24 hours')
    })

    it('rejects a start that is not before the end without calling the API', async () => {
      const store = useScheduleStore()
      store.schedule.push(event('e1'))

      const result = await store.updateSchedule('e1', {
        start: new Date('2030-01-01T11:00:00.000Z'),
        end: new Date('2030-01-01T11:00:00.000Z')
      })

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('Start time must be before end time')
    })

    it('rejects a time slot that conflicts with another event', async () => {
      const store = useScheduleStore()
      store.schedule.push(event('e1'), event('e2', {
        start: new Date('2030-01-01T12:00:00.000Z'),
        end: new Date('2030-01-01T13:00:00.000Z')
      }))

      const result = await store.updateSchedule('e1', {
        start: new Date('2030-01-01T12:30:00.000Z'),
        end: new Date('2030-01-01T13:30:00.000Z')
      })

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('This time slot conflicts with an existing scheduled event')
    })

    it('applies the update optimistically, then reconciles with the server response', async () => {
      const store = useScheduleStore()
      store.schedule.push(event('e1'))
      apiMock.mockResolvedValue(serverEvent('e1', { title: 'Renamed' }))

      const result = await store.updateSchedule('e1', { title: 'Renamed' })

      expect(apiMock).toHaveBeenCalledWith('schedule/e1', {
        method: 'PUT',
        body: JSON.stringify({ title: 'Renamed' })
      })
      expect(result).toBe(true)
      expect(store.schedule[0]).toEqual(event('e1', { title: 'Renamed' }))
      expect(store.schedule[0].start).toBeInstanceOf(Date)
      expect(store.fetchError).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Calendar',
        description: 'Schedule successfully updated',
        color: 'success'
      })
    })

    it('restores the original item and returns false on failure', async () => {
      const store = useScheduleStore()
      const original = event('e1')
      store.schedule.push(original)
      apiMock.mockRejectedValue(new Error('conflict'))

      const result = await store.updateSchedule('e1', { title: 'changed' })

      expect(result).toBe(false)
      expect(store.schedule[0]).toEqual(original)
      expect(store.fetchError).toBeInstanceOf(Error)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Schedule did not update',
        color: 'error'
      })
    })
  })

  describe('deleteSchedule', () => {
    it('removes the item optimistically and succeeds', async () => {
      const store = useScheduleStore()
      store.schedule.push(event('e1'), event('e2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteSchedule('e1')

      expect(apiMock).toHaveBeenCalledWith('schedule/e1', { method: 'DELETE' })
      expect(store.schedule.map(e => e.id)).toEqual(['e2'])
      expect(store.fetchError).toBeNull()
    })

    it('sets fetchError and does not call the API for an unknown id', async () => {
      const store = useScheduleStore()

      await store.deleteSchedule('e-missing')

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.fetchError?.message).toBe('Schedule item not found')
    })

    it('restores the item on failure without toasting', async () => {
      // NOTE (actual behavior): delete failures restore the item and set
      // fetchError, but no toast is shown.
      const store = useScheduleStore()
      const original = event('e1')
      store.schedule.push(original)
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteSchedule('e1')

      expect(store.schedule).toEqual([original])
      expect(store.fetchError).toBeInstanceOf(Error)
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('getSuggestions', () => {
    it('stores the suggestions and persists them to localStorage', async () => {
      const store = useScheduleStore()
      const fetched = [suggestion('s1'), suggestion('s2')]
      apiMock.mockResolvedValue({ suggestions: fetched })

      const result = await store.getSuggestions()

      expect(apiMock).toHaveBeenCalledWith('/schedule/suggestions')
      expect(store.suggestions).toEqual(fetched)
      expect(result).toEqual(fetched)
      expect(JSON.parse(storageMap.get('schedule-suggestions')!)).toEqual(fetched)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Calendar',
        description: 'Got AI suggestions',
        color: 'success'
      })
    })

    it('keeps existing suggestions when the API returns none', async () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'))
      apiMock.mockResolvedValue({ suggestions: [] })

      const result = await store.getSuggestions()

      expect(store.suggestions).toEqual([suggestion('s1')])
      expect(result).toEqual([])
      expect(storageMap.has('schedule-suggestions')).toBe(false)
    })

    it('shows an error toast, caches nothing, and resolves undefined on failure', async () => {
      // NOTE (actual behavior): on failure the action returns undefined (no
      // explicit return) rather than falling back to cached suggestions.
      const store = useScheduleStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.getSuggestions()

      expect(result).toBeUndefined()
      expect(store.fetchError).toBeInstanceOf(Error)
      expect(storageMap.has('schedule-suggestions')).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Could not get AI Suggestions',
        color: 'error'
      })
    })
  })

  describe('localStorage persistence', () => {
    it('loadSuggestionsFromStorage rehydrates persisted suggestions', () => {
      const store = useScheduleStore()
      const persisted = [suggestion('s1'), suggestion('s2')]
      storageMap.set('schedule-suggestions', JSON.stringify(persisted))

      const result = store.loadSuggestionsFromStorage()

      expect(result).toEqual(persisted)
      expect(store.suggestions).toEqual(persisted)
    })

    it('loadSuggestionsFromStorage returns an empty list when nothing is stored', () => {
      const store = useScheduleStore()

      const result = store.loadSuggestionsFromStorage()

      expect(result).toEqual([])
      expect(store.suggestions).toEqual([])
    })

    it('loadSuggestionsFromStorage tolerates corrupted storage', () => {
      const store = useScheduleStore()
      storageMap.set('schedule-suggestions', '{not json')

      const result = store.loadSuggestionsFromStorage()

      expect(result).toEqual([])
      expect(store.suggestions).toEqual([])
    })

    it('removeSuggestion filters the list and persists the result', () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'), suggestion('s2'))

      store.removeSuggestion('s1')

      expect(store.suggestions.map(s => s.id)).toEqual(['s2'])
      expect(JSON.parse(storageMap.get('schedule-suggestions')!)).toEqual([suggestion('s2')])
    })

    it('clearSuggestions empties the list and removes the storage key', () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'))
      storageMap.set('schedule-suggestions', JSON.stringify([suggestion('s1')]))

      store.clearSuggestions()

      expect(store.suggestions).toEqual([])
      expect(storageMap.has('schedule-suggestions')).toBe(false)
    })
  })

  describe('applySuggestions', () => {
    it('warns and returns false when no suggestions match the ids', async () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'))

      const result = await store.applySuggestions(['s-other'])

      expect(result).toBe(false)
      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Warning',
        description: 'No suggestions to apply',
        color: 'warning'
      })
    })

    it('posts the selected suggestions in bulk, prunes them, and refetches the schedule', async () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'), suggestion('s2'))
      apiMock.mockResolvedValueOnce(undefined).mockResolvedValueOnce({ schedule: [] })

      const result = await store.applySuggestions(['s1'])

      expect(apiMock).toHaveBeenNthCalledWith(1, 'schedule/bulk', {
        method: 'POST',
        body: JSON.stringify([{ title: 'Suggestion s1', start: '2030-01-01T10:00:00.000Z', end: '2030-01-01T11:00:00.000Z' }])
      })
      expect(apiMock).toHaveBeenNthCalledWith(2, '/schedule')
      expect(result).toBe(true)
      expect(store.suggestions.map(s => s.id)).toEqual(['s2'])
      expect(JSON.parse(storageMap.get('schedule-suggestions')!)).toEqual([suggestion('s2')])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Calendar',
        description: 'Added 1 suggestions to schedule',
        color: 'success'
      })
    })

    it('returns false, keeps the suggestions, and shows an error toast on failure', async () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'))
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.applySuggestions(['s1'])

      expect(result).toBe(false)
      expect(store.suggestions.map(s => s.id)).toEqual(['s1'])
      expect(store.fetchError).toBeInstanceOf(Error)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Could not apply suggestions',
        color: 'error'
      })
    })
  })

  describe('updateSuggestion', () => {
    it('updates the timing and persists the change', () => {
      const store = useScheduleStore()
      store.suggestions.push(suggestion('s1'))

      const result = store.updateSuggestion('s1', '2030-01-02T09:00:00.000Z', '2030-01-02T10:00:00.000Z')

      expect(result).toBe(true)
      expect(store.suggestions[0].start).toBe('2030-01-02T09:00:00.000Z')
      expect(JSON.parse(storageMap.get('schedule-suggestions')!)[0].start).toBe('2030-01-02T09:00:00.000Z')
    })

    it('returns false for an unknown suggestion', () => {
      const store = useScheduleStore()

      const result = store.updateSuggestion('s-missing', '2030-01-02T09:00:00.000Z', '2030-01-02T10:00:00.000Z')

      expect(result).toBe(false)
      expect(storageMap.has('schedule-suggestions')).toBe(false)
    })
  })

  it('reset clears the schedule but keeps suggestions', () => {
    // NOTE (actual behavior): reset() only clears the schedule; suggestions
    // (and localStorage) are left untouched.
    const store = useScheduleStore()
    store.schedule.push(event('e1'))
    store.suggestions.push(suggestion('s1'))

    store.reset()

    expect(store.schedule).toEqual([])
    expect(store.suggestions).toEqual([suggestion('s1')])
  })
})
