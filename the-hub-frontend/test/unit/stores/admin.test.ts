import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
import { useAdminStore } from '../../../app/stores/admin'

// The admin store reaches for Nuxt auto-imports (useNuxtApp, useToast, ref).
// Stub them; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('ref', ref)
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

const adminUser = (id: string, role = 'user') => ({
  user_id: id,
  name: `User ${id}`,
  email: `${id}@example.com`,
  role,
  created_at: '2026-01-01T00:00:00.000Z',
  updated_at: '2026-01-01T00:00:00.000Z'
})

describe('useAdminStore', () => {
  describe('getAllUsers', () => {
    it('populates users and clears loading', async () => {
      const store = useAdminStore()
      apiMock.mockResolvedValue({ users: [adminUser('u1'), adminUser('u2')] })

      await store.getAllUsers()

      expect(apiMock).toHaveBeenCalledWith('/admin/users')
      expect(store.users).toEqual([adminUser('u1'), adminUser('u2')])
      expect(store.loading).toBe(false)
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useAdminStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.getAllUsers()
      expect(store.loading).toBe(true)

      resolveApi({ users: [] })
      await pending
      expect(store.loading).toBe(false)
    })

    it('shows an error toast and keeps existing users on failure', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      apiMock.mockRejectedValue(new Error('forbidden'))

      await store.getAllUsers()

      expect(store.users).toEqual([adminUser('u1')])
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch users', color: 'error' })
    })
  })

  describe('getSystemStats', () => {
    it('stores the stats payload', async () => {
      const store = useAdminStore()
      const stats = {
        total_users: 10,
        active_users: 8,
        total_tasks: 42,
        total_goals: 5,
        total_decks: 3,
        total_study_sessions: 120
      }
      apiMock.mockResolvedValue(stats)

      await store.getSystemStats()

      expect(apiMock).toHaveBeenCalledWith('/admin/stats')
      expect(store.stats).toEqual(stats)
      expect(store.loading).toBe(false)
    })

    it('shows an error toast and keeps stats null on failure', async () => {
      const store = useAdminStore()
      apiMock.mockRejectedValue(new Error('forbidden'))

      await store.getSystemStats()

      expect(store.stats).toBeNull()
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch system statistics', color: 'error' })
    })
  })

  describe('updateUserRole', () => {
    it('calls the API, replaces the user in the list and toasts success', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      const promoted = adminUser('u1', 'admin')
      apiMock.mockResolvedValue(promoted)

      const result = await store.updateUserRole('u1', 'admin')

      expect(apiMock).toHaveBeenCalledWith('/admin/users/u1/role', {
        method: 'PUT',
        body: JSON.stringify({ role: 'admin' })
      })
      expect(store.users[0]).toEqual(promoted)
      expect(result).toEqual(promoted)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'User role updated to admin', color: 'success' })
    })

    it('still toasts success for a user id that is not in the list', async () => {
      // NOTE: the store cannot find the user locally but reports success
      // anyway — the list stays unchanged.
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      apiMock.mockResolvedValue(adminUser('u9', 'admin'))

      const result = await store.updateUserRole('u9', 'admin')

      expect(store.users.map(u => u.user_id)).toEqual(['u1'])
      expect(result).toEqual(adminUser('u9', 'admin'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'User role updated to admin', color: 'success' })
    })

    it('toasts an error and rethrows on failure', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      apiMock.mockRejectedValue(new Error('forbidden'))

      await expect(store.updateUserRole('u1', 'admin')).rejects.toThrow('forbidden')

      expect(store.users[0]).toEqual(adminUser('u1'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to update user role', color: 'error' })
    })
  })

  describe('promoteToAdmin', () => {
    it('calls the API, replaces the user in the list and toasts success', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      const promoted = adminUser('u1', 'admin')
      apiMock.mockResolvedValue(promoted)

      const result = await store.promoteToAdmin('u1')

      expect(apiMock).toHaveBeenCalledWith('/admin/users/u1/promote', { method: 'POST' })
      expect(store.users[0]).toEqual(promoted)
      expect(result).toEqual(promoted)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'User promoted to admin successfully', color: 'success' })
    })

    it('toasts an error and rethrows on failure', async () => {
      const store = useAdminStore()
      apiMock.mockRejectedValue(new Error('forbidden'))

      await expect(store.promoteToAdmin('u1')).rejects.toThrow('forbidden')

      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to promote user to admin', color: 'error' })
    })
  })

  describe('deleteUser', () => {
    it('calls the API, removes the user and toasts success', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'), adminUser('u2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteUser('u1')

      expect(apiMock).toHaveBeenCalledWith('/admin/users/u1', { method: 'DELETE' })
      expect(store.users.map(u => u.user_id)).toEqual(['u2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'User deleted successfully', color: 'success' })
    })

    it('keeps the user and rethrows on failure', async () => {
      const store = useAdminStore()
      store.users.push(adminUser('u1'))
      apiMock.mockRejectedValue(new Error('forbidden'))

      await expect(store.deleteUser('u1')).rejects.toThrow('forbidden')

      expect(store.users.map(u => u.user_id)).toEqual(['u1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to delete user', color: 'error' })
    })
  })
})
