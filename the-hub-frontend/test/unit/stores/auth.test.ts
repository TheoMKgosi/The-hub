import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { ref, computed, watch, nextTick } from 'vue'

// The auth store reaches for Nuxt auto-imports (useNuxtApp, useToast, useRouter,
// $fetch, useRuntimeConfig, ref/computed/watch) and imports a dozen sibling
// stores that it resets on logout. Stub the Vue globals with the real
// implementations, mock the sibling store modules, and stand in $api/$fetch
// with mocks.
// NOTE: auth.ts imports useValidation via the `@/` alias, which does not
// resolve in the unit project — re-map that specifier to the real module below.
const apiMock = vi.fn()
const fetchMock = vi.fn()
const toastAddMock = vi.fn()
const pushMock = vi.fn()
const replaceMock = vi.fn()

const storeMocks = vi.hoisted(() => {
  const resets: Record<string, ReturnType<typeof vi.fn>> = {}
  const makeComposable = (key: string) => () => {
    resets[key] ??= vi.fn()
    return { reset: resets[key] }
  }
  return { resets, makeComposable }
})

vi.mock('@/composables/useValidation', async () => await import('../../../app/composables/useValidation'))

vi.mock('../../../app/stores/topics', () => ({ useTopicStore: storeMocks.makeComposable('topic') }))
vi.mock('../../../app/stores/tasks', () => ({ useTaskStore: storeMocks.makeComposable('task') }))
vi.mock('../../../app/stores/tags', () => ({ useTagStore: storeMocks.makeComposable('tag') }))
vi.mock('../../../app/stores/goals', () => ({ useGoalStore: storeMocks.makeComposable('goal') }))
vi.mock('../../../app/stores/decks', () => ({ useDeckStore: storeMocks.makeComposable('deck') }))
vi.mock('../../../app/stores/cards', () => ({ useCardStore: storeMocks.makeComposable('card') }))
vi.mock('../../../app/stores/income', () => ({ useIncomeStore: storeMocks.makeComposable('income') }))
vi.mock('../../../app/stores/schedule', () => ({ useScheduleStore: storeMocks.makeComposable('schedule') }))
vi.mock('../../../app/stores/finance', () => ({
  useBudgetStore: storeMocks.makeComposable('budget'),
  useCategoryStore: storeMocks.makeComposable('category')
}))
vi.mock('../../../app/stores/financialGoal', () => ({ useFinancialGoalStore: storeMocks.makeComposable('financialGoal') }))
vi.mock('../../../app/stores/themes', () => ({ useThemeStore: storeMocks.makeComposable('theme') }))

// In-memory stand-in for the browser localStorage the store persists to.
const storageData = new Map<string, string>()
const localStorageMock = {
  getItem: vi.fn((key: string) => storageData.get(key) ?? null),
  setItem: vi.fn((key: string, value: string) => { storageData.set(key, String(value)) }),
  removeItem: vi.fn((key: string) => { storageData.delete(key) }),
  clear: vi.fn(() => storageData.clear())
}

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('useRouter', () => ({ push: pushMock, replace: replaceMock }))
  vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: 'https://api.example.com' } }))
  vi.stubGlobal('$fetch', fetchMock)
  vi.stubGlobal('localStorage', localStorageMock)
  vi.stubGlobal('ref', ref)
  vi.stubGlobal('computed', computed)
  vi.stubGlobal('watch', watch)
})

const { useAuthStore } = await import('../../../app/stores/auth')

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  fetchMock.mockReset()
  toastAddMock.mockClear()
  pushMock.mockClear()
  replaceMock.mockClear()
  storageData.clear()
  for (const fn of Object.values(localStorageMock)) fn.mockClear()
  vi.spyOn(console, 'warn').mockImplementation(() => {})
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  // The persistence watcher is only registered when `process.client` is set —
  // individual tests opt in and must not leak the flag into other tests.
  ;(process as unknown as { client?: boolean }).client = false
  vi.restoreAllMocks()
})

const makeUser = (role = 'user') => ({
  user_id: 'u1',
  name: 'Test User',
  email: 'test@example.com',
  role
})

const makeAuthResponse = (role = 'user', expiresIn = 3600) => ({
  access_token: 'at-1',
  refresh_token: 'rt-1',
  token_type: 'bearer',
  expires_in: expiresIn,
  user: makeUser(role)
})

describe('useAuthStore', () => {
  describe('initial state', () => {
    it('starts with empty credentials and no user', () => {
      const store = useAuthStore()

      expect(store.user).toBeNull()
      expect(store.accessToken).toBe('')
      expect(store.refreshToken).toBe('')
      expect(store.isLoggedIn).toBe(false)
    })

    it('reports the token as expired when no expiry is recorded', () => {
      const store = useAuthStore()

      expect(store.isTokenExpired()).toBe(true)
    })
  })

  describe('login', () => {
    it('stores tokens and user data on success', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse())

      await store.login({ email: 'test@example.com', password: 'secret' })

      expect(apiMock).toHaveBeenCalledWith('/login', {
        method: 'POST',
        body: JSON.stringify({ email: 'test@example.com', password: 'secret' })
      })
      expect(store.accessToken).toBe('at-1')
      expect(store.refreshToken).toBe('rt-1')
      expect(store.user).toEqual(makeUser())
      expect(store.isLoggedIn).toBe(true)
      expect(store.isTokenExpired()).toBe(false)
    })

    it('redirects a normal user to /dashboard with a welcome toast', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('user'))

      await store.login({ email: 'test@example.com', password: 'secret' })

      expect(pushMock).toHaveBeenCalledWith('/dashboard')
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Success', description: 'Welcome back!', color: 'success' })
      )
    })

    it('redirects an admin to /admin with an admin toast', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('admin'))

      await store.login({ email: 'admin@example.com', password: 'secret' })

      expect(pushMock).toHaveBeenCalledWith('/admin')
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Success', description: 'Welcome to admin panel', color: 'success' })
      )
    })

    it('rejects with the validation message and never calls the API for an invalid email', async () => {
      const store = useAuthStore()

      await expect(store.login({ email: 'not-an-email', password: 'secret' })).rejects.toThrow(
        'email format is invalid'
      )

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).not.toHaveBeenCalled()
      // NOTE: login rethrows errors without toasting — the calling page is
      // expected to catch and display them (register/logout behave differently).
    })

    it('rethrows API errors and leaves auth state untouched', async () => {
      const store = useAuthStore()
      apiMock.mockRejectedValue(new Error('invalid credentials'))

      await expect(store.login({ email: 'test@example.com', password: 'wrong' })).rejects.toThrow(
        'invalid credentials'
      )

      expect(apiMock).toHaveBeenCalledWith('/login', expect.anything())
      expect(store.accessToken).toBe('')
      expect(store.user).toBeNull()
      expect(store.isLoggedIn).toBe(false)
      expect(toastAddMock).not.toHaveBeenCalled()
      expect(pushMock).not.toHaveBeenCalled()
    })
  })

  describe('register', () => {
    it('stores tokens and user data, redirects to /dashboard and toasts success', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse())

      await store.register({ name: 'Test User', email: 'test@example.com', password: 'secret' })

      expect(apiMock).toHaveBeenCalledWith('/register', {
        method: 'POST',
        body: JSON.stringify({ name: 'Test User', email: 'test@example.com', password: 'secret' })
      })
      expect(store.accessToken).toBe('at-1')
      expect(store.user).toEqual(makeUser())
      expect(pushMock).toHaveBeenCalledWith('/dashboard')
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Success', description: 'Account created successfully!', color: 'success' })
      )
    })

    it('toasts the first validation error and skips the API call', async () => {
      const store = useAuthStore()

      await store.register({ name: 'A', email: 'not-an-email', password: 'secret' })

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', description: 'name must be at least 2 characters', color: 'error' })
      )
      expect(pushMock).not.toHaveBeenCalled()
    })

    it('toasts the API error message and swallows the failure', async () => {
      const store = useAuthStore()
      apiMock.mockRejectedValue(new Error('email already taken'))

      await store.register({ name: 'Test User', email: 'test@example.com', password: 'secret' })

      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', description: 'email already taken', color: 'error' })
      )
      expect(store.accessToken).toBe('')
    })
  })

  describe('logout', () => {
    it('revokes the refresh token, clears all state and resets the other stores', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse())
      await store.login({ email: 'test@example.com', password: 'secret' })
      apiMock.mockReset()
      apiMock.mockResolvedValue(undefined)

      await store.logout()

      expect(apiMock).toHaveBeenCalledWith('/auth/logout', {
        method: 'POST',
        body: JSON.stringify({ refresh_token: 'rt-1' })
      })
      expect(store.accessToken).toBeNull()
      expect(store.refreshToken).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isLoggedIn).toBe(false)
      expect(pushMock).toHaveBeenCalledWith('/login')

      const expectedResets = ['task', 'goal', 'schedule', 'tag', 'topic', 'deck', 'card', 'income', 'budget', 'category', 'financialGoal', 'theme']
      for (const key of expectedResets) {
        expect(storeMocks.resets[key], key).toHaveBeenCalledOnce()
      }
    })

    it('clears the localStorage auth keys', async () => {
      const store = useAuthStore()
      storageData.set('accessToken', 'at-1')
      storageData.set('refreshToken', 'rt-1')
      storageData.set('tokenExpiry', '123')
      storageData.set('user', '{"user_id":"u1"}')
      apiMock.mockResolvedValue(makeAuthResponse())
      await store.login({ email: 'test@example.com', password: 'secret' })

      await store.logout()

      expect(localStorageMock.removeItem).toHaveBeenCalledWith('accessToken')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('refreshToken')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('tokenExpiry')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('user')
      expect(storageData.size).toBe(0)
    })

    it('still clears local state when the revoke API call fails', async () => {
      const store = useAuthStore()
      store.refreshToken = 'rt-1'
      apiMock.mockRejectedValue(new Error('network down'))

      await store.logout()

      expect(store.accessToken).toBeNull()
      expect(store.refreshToken).toBeNull()
      expect(store.user).toBeNull()
      expect(pushMock).toHaveBeenCalledWith('/login')
      expect(console.warn).toHaveBeenCalled()
    })

    it('does not call the API when no refresh token exists', async () => {
      const store = useAuthStore()

      await store.logout()

      expect(apiMock).not.toHaveBeenCalled()
      expect(store.accessToken).toBeNull()
      expect(pushMock).toHaveBeenCalledWith('/login')
    })
  })

  describe('forgotPassword', () => {
    it('posts the email and returns success', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(undefined)

      const result = await store.forgotPassword({ email: 'test@example.com' })

      expect(apiMock).toHaveBeenCalledWith('/auth/forgot-password', {
        method: 'POST',
        body: JSON.stringify({ email: 'test@example.com' })
      })
      expect(result).toEqual({ success: true })
      expect(toastAddMock).toHaveBeenCalledWith(expect.objectContaining({ title: 'Success', color: 'success' }))
    })

    it('rejects and toasts for an invalid email', async () => {
      const store = useAuthStore()

      await expect(store.forgotPassword({ email: 'nope' })).rejects.toThrow('email format is invalid')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', description: 'email format is invalid', color: 'error' })
      )
    })

    it('toasts and rethrows API errors', async () => {
      const store = useAuthStore()
      apiMock.mockRejectedValue(new Error('rate limited'))

      await expect(store.forgotPassword({ email: 'test@example.com' })).rejects.toThrow('rate limited')

      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', description: 'rate limited', color: 'error' })
      )
    })
  })

  describe('resetPassword', () => {
    it('posts the token and password and returns success', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(undefined)
      const payload = { token: 'a'.repeat(32), password: 'new-secret' }

      const result = await store.resetPassword(payload)

      expect(apiMock).toHaveBeenCalledWith('/auth/reset-password', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(result).toEqual({ success: true })
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Success', description: 'Password reset successful', color: 'success' })
      )
    })

    it('rejects and toasts when the token is too short', async () => {
      const store = useAuthStore()

      await expect(store.resetPassword({ token: 'short', password: 'new-secret' })).rejects.toThrow(
        'token must be at least 32 characters'
      )

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', color: 'error' })
      )
    })

    it('toasts and rethrows API errors', async () => {
      const store = useAuthStore()
      apiMock.mockRejectedValue(new Error('expired token'))

      await expect(store.resetPassword({ token: 'a'.repeat(32), password: 'new-secret' })).rejects.toThrow(
        'expired token'
      )

      expect(toastAddMock).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'Error', description: 'expired token', color: 'error' })
      )
    })
  })

  describe('isTokenExpired', () => {
    it('is false for a token expiring in one hour', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('user', 3600))

      await store.login({ email: 'test@example.com', password: 'secret' })

      expect(store.isTokenExpired()).toBe(false)
    })

    it('is true for a token that expired immediately', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('user', 0))

      await store.login({ email: 'test@example.com', password: 'secret' })

      expect(store.isTokenExpired()).toBe(true)
    })
  })

  describe('refreshAccessToken', () => {
    it('throws when no refresh token is available', async () => {
      const store = useAuthStore()

      await expect(store.refreshAccessToken()).rejects.toThrow('No refresh token available')
      expect(fetchMock).not.toHaveBeenCalled()
    })

    it('refreshes via $fetch with the runtime apiBase and updates the access token', async () => {
      const store = useAuthStore()
      store.refreshToken = 'rt-1'
      fetchMock.mockResolvedValue({ access_token: 'at-2', expires_in: 3600 })

      await store.refreshAccessToken()

      expect(fetchMock).toHaveBeenCalledWith('/auth/refresh', {
        baseURL: 'https://api.example.com',
        method: 'POST',
        body: JSON.stringify({ refresh_token: 'rt-1' })
      })
      expect(store.accessToken).toBe('at-2')
      expect(store.isTokenExpired()).toBe(false)
    })

    it('shares a single in-flight request between concurrent callers', async () => {
      const store = useAuthStore()
      store.refreshToken = 'rt-1'
      let resolveFetch!: (value: unknown) => void
      fetchMock.mockReturnValue(new Promise((resolve) => { resolveFetch = resolve }))

      const first = store.refreshAccessToken()
      const second = store.refreshAccessToken()

      expect(fetchMock).toHaveBeenCalledTimes(1)
      resolveFetch({ access_token: 'at-2', expires_in: 3600 })
      await first
      await second
      expect(store.accessToken).toBe('at-2')
    })

    it('makes a fresh request after a failed refresh (in-flight state reset in finally)', async () => {
      const store = useAuthStore()
      store.refreshToken = 'rt-1'
      fetchMock.mockRejectedValueOnce(new Error('refresh failed'))
      fetchMock.mockResolvedValueOnce({ access_token: 'at-2', expires_in: 3600 })

      await expect(store.refreshAccessToken()).rejects.toThrow('refresh failed')

      const token = await store.refreshAccessToken()
      expect(token).toEqual({ access_token: 'at-2', expires_in: 3600 })
      expect(store.accessToken).toBe('at-2')
    })
  })

  describe('getAccessToken', () => {
    it('returns null when there is no access token', async () => {
      const store = useAuthStore()

      expect(await store.getAccessToken()).toBeNull()
      expect(fetchMock).not.toHaveBeenCalled()
    })

    it('returns the current token without refreshing when it is still valid', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('user', 3600))
      await store.login({ email: 'test@example.com', password: 'secret' })

      const token = await store.getAccessToken()

      expect(token).toBe('at-1')
      expect(fetchMock).not.toHaveBeenCalled()
    })

    it('refreshes first and returns the new token when the current one is expired', async () => {
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse('user', 0))
      await store.login({ email: 'test@example.com', password: 'secret' })
      fetchMock.mockResolvedValue({ access_token: 'at-2', expires_in: 3600 })

      const token = await store.getAccessToken()

      expect(token).toBe('at-2')
      expect(store.accessToken).toBe('at-2')
      expect(fetchMock).toHaveBeenCalledTimes(1)
    })
  })

  describe('localStorage persistence', () => {
    it('writes tokens and user to localStorage on the client', async () => {
      // Opt into the client-only persistence watcher for this test.
      ;(process as unknown as { client?: boolean }).client = true
      const store = useAuthStore()
      apiMock.mockResolvedValue(makeAuthResponse())

      await store.login({ email: 'test@example.com', password: 'secret' })
      await nextTick()

      expect(localStorageMock.setItem).toHaveBeenCalledWith('accessToken', 'at-1')
      expect(localStorageMock.setItem).toHaveBeenCalledWith('refreshToken', 'rt-1')
      expect(localStorageMock.setItem).toHaveBeenCalledWith('user', JSON.stringify(makeUser()))
      expect(storageData.get('accessToken')).toBe('at-1')
      expect(storageData.get('user')).toBe(JSON.stringify(makeUser()))
    })
  })
})
