import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useThemeStore } from '../../../app/stores/themes'

// The theme store reaches for Nuxt auto-imports (useNuxtApp, useToast,
// useAuthStore). Stub them all; the $api mock stands in for the API plugin and
// the auth stub is backed by a mutable object so login state can be toggled.
const apiMock = vi.fn()
const toastAddMock = vi.fn()
const authState = { isLoggedIn: true, user: { user_id: 'u1' } }

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('useAuthStore', () => authState)
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  authState.isLoggedIn = true
  authState.user = { user_id: 'u1' }
})

afterEach(() => {
  vi.restoreAllMocks()
})

const theme = (id: string, days: string[], overrides: Record<string, unknown> = {}) => ({
  theme_id: id,
  user_id: 'u1',
  name: `theme-${id}`,
  description: '',
  color: '#00ff00',
  days,
  ...overrides
})

const todayWeekday = () => new Date().toLocaleDateString('en-US', { weekday: 'long' }).toLowerCase()

// A weekday guaranteed to be different from today (never matches activeThemeToday)
const otherWeekday = () => (todayWeekday() === 'monday' ? 'tuesday' : 'monday')

const appTask = (id: string, themeId?: string) => ({
  task_id: id,
  title: `task-${id}`,
  description: '',
  status: 'not_started',
  start_time: new Date('2026-11-02T10:00:00.000Z'),
  time_spent_minutes: 0,
  is_recurring: false,
  theme_id: themeId
})

describe('useThemeStore', () => {
  it('starts with no themes, suggestions, feature disabled and no error', () => {
    const store = useThemeStore()

    expect(store.themes).toEqual([])
    expect(store.suggestions).toEqual([])
    expect(store.enabled).toBe(false)
    expect(store.loading).toBe(false)
    expect(store.suggestLoading).toBe(false)
    expect(store.fetchError).toBeNull()
  })

  describe('fetchThemes', () => {
    it('populates themes and clears loading', async () => {
      const store = useThemeStore()
      const serverThemes = [theme('th1', ['monday']), theme('th2', ['friday'])]
      apiMock.mockResolvedValue({ themes: serverThemes })

      await store.fetchThemes()

      expect(apiMock).toHaveBeenCalledWith('/themes')
      expect(store.themes).toEqual(serverThemes)
      expect(store.loading).toBe(false)
    })

    it('records the error, shows a toast and clears loading on failure', async () => {
      const store = useThemeStore()
      const failure = new Error('server error')
      apiMock.mockRejectedValue(failure)

      await store.fetchThemes()

      expect(store.fetchError).toBe(failure)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Failed to fetch themes',
        color: 'error'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchEnabled', () => {
    it('enables theme days when the setting is on', async () => {
      const store = useThemeStore()
      apiMock.mockResolvedValue({ settings: { theme_days_enabled: true } })

      await store.fetchEnabled()

      expect(apiMock).toHaveBeenCalledWith('/users/u1/settings')
      expect(store.enabled).toBe(true)
    })

    it('keeps theme days disabled when the setting is off', async () => {
      const store = useThemeStore()
      apiMock.mockResolvedValue({ settings: { theme_days_enabled: false } })

      await store.fetchEnabled()

      expect(store.enabled).toBe(false)
    })

    it('does not call the API when the user is logged out', async () => {
      const store = useThemeStore()
      authState.isLoggedIn = false

      await store.fetchEnabled()

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.enabled).toBe(false)
    })

    it('leaves the setting untouched when the response has no settings', async () => {
      const store = useThemeStore()
      apiMock.mockResolvedValue(null)

      await store.fetchEnabled()

      expect(store.enabled).toBe(false)
    })
  })

  describe('setEnabled', () => {
    it('persists the setting and updates local state', async () => {
      const store = useThemeStore()
      apiMock.mockResolvedValue(undefined)

      await store.setEnabled(true)

      expect(apiMock).toHaveBeenCalledWith('/users/u1/settings', {
        method: 'PATCH',
        body: JSON.stringify({ theme_days_enabled: true })
      })
      expect(store.enabled).toBe(true)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme Days',
        description: 'Theme Days enabled',
        color: 'success'
      })
    })

    it('shows an error toast and keeps state when the user is logged out', async () => {
      const store = useThemeStore()
      authState.isLoggedIn = false

      await store.setEnabled(true)

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.enabled).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'You must be logged in',
        color: 'error'
      })
    })

    it('shows an error toast and keeps state on failure', async () => {
      const store = useThemeStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.setEnabled(true)

      expect(store.enabled).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Failed to update setting',
        color: 'error'
      })
    })
  })

  describe('createTheme', () => {
    it('appends the created theme and returns it', async () => {
      const store = useThemeStore()
      const created = theme('th1', ['monday'])
      apiMock.mockResolvedValue(created)

      const result = await store.createTheme({ name: 'theme-th1', description: '', color: '#00ff00', days: ['monday'] })

      expect(apiMock).toHaveBeenCalledWith('/themes', {
        method: 'POST',
        body: JSON.stringify({ name: 'theme-th1', description: '', color: '#00ff00', days: ['monday'] })
      })
      expect(store.themes).toEqual([created])
      expect(result).toEqual(created)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Theme created successfully',
        color: 'success'
      })
    })

    it('rethrows and surfaces the API error message', async () => {
      const store = useThemeStore()
      const failure = Object.assign(new Error('bad request'), { data: { error: 'Name already used' } })
      apiMock.mockRejectedValue(failure)

      await expect(
        store.createTheme({ name: 'theme-th1', description: '', color: '#00ff00', days: ['monday'] })
      ).rejects.toThrow('bad request')

      expect(store.themes).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Name already used',
        color: 'error'
      })
    })
  })

  describe('updateTheme', () => {
    it('replaces the theme with the server response', async () => {
      const store = useThemeStore()
      store.themes.push(theme('th1', ['monday']))
      const updated = theme('th1', ['tuesday'], { name: 'renamed' })
      apiMock.mockResolvedValue(updated)

      const result = await store.updateTheme('th1', { name: 'renamed', description: '', color: '#00ff00', days: ['tuesday'] })

      expect(apiMock).toHaveBeenCalledWith('/themes/th1', {
        method: 'PATCH',
        body: JSON.stringify({ name: 'renamed', description: '', color: '#00ff00', days: ['tuesday'] })
      })
      expect(store.themes[0]).toEqual(updated)
      expect(result).toEqual(updated)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Theme updated successfully',
        color: 'success'
      })
    })

    it('reverts to the original theme on failure and rethrows', async () => {
      const store = useThemeStore()
      store.themes.push(theme('th1', ['monday']))
      apiMock.mockRejectedValue(new Error('conflict'))

      await expect(
        store.updateTheme('th1', { name: 'renamed', description: '', color: '#00ff00', days: ['tuesday'] })
      ).rejects.toThrow('conflict')

      expect(store.themes[0]).toEqual(theme('th1', ['monday']))
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Failed to update theme',
        color: 'error'
      })
    })
  })

  describe('deleteTheme', () => {
    it('removes the theme and succeeds', async () => {
      const store = useThemeStore()
      store.themes.push(theme('th1', ['monday']), theme('th2', ['friday']))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTheme('th1')

      expect(apiMock).toHaveBeenCalledWith('/themes/th1', { method: 'DELETE' })
      expect(store.themes.map(t => t.theme_id)).toEqual(['th2'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Theme deleted successfully',
        color: 'success'
      })
    })

    it('keeps the theme on failure and rethrows', async () => {
      const store = useThemeStore()
      store.themes.push(theme('th1', ['monday']))
      apiMock.mockRejectedValue(new Error('server error'))

      await expect(store.deleteTheme('th1')).rejects.toThrow('server error')

      expect(store.themes.map(t => t.theme_id)).toEqual(['th1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Failed to delete theme',
        color: 'error'
      })
    })
  })

  describe('getAISuggestions', () => {
    it('stores and returns the suggestions payload', async () => {
      const store = useThemeStore()
      const serverSuggestions = [
        { name: 'Deep Work', description: '', color: '#0000ff', days: ['monday'], tasks: [], reasoning: '' }
      ]
      apiMock.mockResolvedValue({ suggestions: serverSuggestions })

      const result = await store.getAISuggestions()

      expect(apiMock).toHaveBeenCalledWith('/themes/ai/suggest', { method: 'POST' })
      expect(store.suggestions).toEqual(serverSuggestions)
      expect(result).toEqual(serverSuggestions)
      expect(store.suggestLoading).toBe(false)
    })

    it('returns an empty list and shows a toast on failure', async () => {
      const store = useThemeStore()
      apiMock.mockRejectedValue(new Error('server error'))

      const result = await store.getAISuggestions()

      expect(result).toEqual([])
      expect(store.suggestions).toEqual([])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Failed to get AI theme suggestions',
        color: 'error'
      })
      expect(store.suggestLoading).toBe(false)
    })
  })

  describe('applyAISuggestions', () => {
    it('applies themes, refetches themes, clears suggestions and returns true', async () => {
      const store = useThemeStore()
      store.suggestions.push({ name: 'Deep Work', description: '', color: '#0000ff', days: ['monday'], tasks: [], reasoning: '' })
      const payload = {
        themes: [{ name: 'Deep Work', description: '', color: '#0000ff', days: ['monday' as const] }],
        assignments: [{ task_id: 't1', theme_name: 'Deep Work' }]
      }
      apiMock
        .mockResolvedValueOnce(undefined)
        .mockResolvedValueOnce({ themes: [theme('th1', ['monday'])] })

      const result = await store.applyAISuggestions(payload)

      expect(apiMock).toHaveBeenNthCalledWith(1, '/themes/ai/apply', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      // fetchThemes follows the apply call
      expect(apiMock).toHaveBeenNthCalledWith(2, '/themes')
      expect(result).toBe(true)
      expect(store.suggestions).toEqual([])
      expect(store.themes).toHaveLength(1)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Themes applied successfully',
        color: 'success'
      })
      expect(store.suggestLoading).toBe(false)
    })

    it('returns false and surfaces the API error message on failure', async () => {
      const store = useThemeStore()
      const failure = Object.assign(new Error('bad request'), { data: { error: 'Invalid themes' } })
      apiMock.mockRejectedValue(failure)

      const result = await store.applyAISuggestions({ themes: [], assignments: [] })

      expect(result).toBe(false)
      expect(apiMock).toHaveBeenCalledTimes(1)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Theme',
        description: 'Invalid themes',
        color: 'error'
      })
      expect(store.suggestLoading).toBe(false)
    })
  })

  describe('theme lookup and task filtering', () => {
    // NOTE: themeForWeekday is defined in the store but NOT returned/exposed,
    // so the weekday lookup is only observable through activeThemeToday.
    it('activeThemeToday resolves the theme for today\'s weekday', async () => {
      const store = useThemeStore()
      const todayTheme = theme('th-today', [todayWeekday()])
      store.themes.push(theme('th-other', [otherWeekday()]), todayTheme)

      expect(store.activeThemeToday).toEqual(todayTheme)
    })

    it('filterTasksByTheme passes tasks through when theme days are disabled', async () => {
      const store = useThemeStore()
      store.themes.push(theme('th1', [todayWeekday()]))
      const tasks = [appTask('t1', 'th1'), appTask('t2', 'th2')]

      expect(store.filterTasksByTheme(tasks)).toEqual(tasks)
    })

    it('filterTasksByTheme hides everything when enabled with no theme today', async () => {
      const store = useThemeStore()
      store.enabled = true
      const tasks = [appTask('t1', 'th1')]

      expect(store.filterTasksByTheme(tasks)).toEqual([])
    })

    it('filterTasksByTheme keeps only tasks assigned to today\'s theme when enabled', async () => {
      const store = useThemeStore()
      const todayTheme = theme('th-today', [todayWeekday()])
      store.themes.push(todayTheme)
      store.enabled = true

      const tasks = [appTask('t1', 'th-today'), appTask('t2', 'th-other')]
      expect(store.filterTasksByTheme(tasks)).toEqual([tasks[0]])
    })
  })

  it('reset clears themes, suggestions and the enabled flag', () => {
    const store = useThemeStore()
    store.themes.push(theme('th1', ['monday']))
    store.suggestions.push({ name: 'Deep Work', description: '', color: '#0000ff', days: ['monday'], tasks: [], reasoning: '' })
    store.enabled = true

    store.reset()

    expect(store.themes).toHaveLength(0)
    expect(store.suggestions).toHaveLength(0)
    expect(store.enabled).toBe(false)
  })
})
