import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref } from 'vue'

// The resource store relies on the Vue auto-import `ref` (it is not imported
// explicitly in the module). Stub the real implementation before the store
// module loads, then import it dynamically; the $api mock stands in for the
// API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('ref', ref)
vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))

const { useResourceStore } = await import('../../../app/stores/resources')

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

const resource = (id: string, overrides: Record<string, unknown> = {}) => ({
  id,
  topic_id: 't1',
  task_id: undefined,
  title: `resource-${id}`,
  link: `https://example.com/${id}`,
  type: 'article' as const,
  notes: '',
  ...overrides
})

describe('useResourceStore', () => {
  it('starts with an empty list and no loading', () => {
    const store = useResourceStore()

    expect(store.resources).toEqual([])
    expect(store.loading).toBe(false)
  })

  describe('createResource', () => {
    it('prepends the resource optimistically, then replaces it with the server response', async () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'))
      const data = { title: 'New resource', link: 'https://example.com/new', type: 'video' as const, topic_id: 't1' }
      apiMock.mockResolvedValue(resource('r9', { title: 'New resource', type: 'video' }))

      const result = await store.createResource(data)

      expect(apiMock).toHaveBeenCalledWith('/resources', {
        method: 'POST',
        body: JSON.stringify(data)
      })
      // The new resource is confirmed at the front, the pre-existing one stays behind
      expect(store.resources).toHaveLength(2)
      expect(store.resources[0]).toEqual(resource('r9', { title: 'New resource', type: 'video' }))
      expect(store.resources[1]).toEqual(resource('r1'))
      expect(result).toEqual(resource('r9', { title: 'New resource', type: 'video' }))
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Resource added successfully',
        color: 'success'
      })
    })

    it('rolls back the optimistic resource on failure and returns undefined', async () => {
      const store = useResourceStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.createResource({
        title: 'New resource',
        link: 'https://example.com/new',
        type: 'video'
      })

      expect(store.resources).toHaveLength(0)
      expect(result).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to add resource',
        color: 'error'
      })
    })
  })

  describe('fetchResources', () => {
    it('builds query params from provided filters and sets resources', async () => {
      const store = useResourceStore()
      const serverResources = [resource('r1'), resource('r2')]
      apiMock.mockResolvedValue({ resources: serverResources })

      await store.fetchResources({ topic_id: 't1', type: 'video', search: 'css' })

      expect(apiMock).toHaveBeenCalledWith('/resources?topic_id=t1&type=video&search=css')
      expect(store.resources).toEqual(serverResources)
      expect(store.loading).toBe(false)
    })

    it('requests a bare query when no filters are given', async () => {
      const store = useResourceStore()
      apiMock.mockResolvedValue({ resources: [] })

      await store.fetchResources()

      expect(apiMock).toHaveBeenCalledWith('/resources?')
      expect(store.resources).toEqual([])
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useResourceStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchResources()
      expect(store.loading).toBe(true)

      resolveApi({ resources: [] })
      await pending
      expect(store.loading).toBe(false)
    })

    it('NOTE: documents current behavior — the fetch error toast uses the misspelled color "errro"', async () => {
      const store = useResourceStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchResources()

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to fetch resources',
        color: 'errro'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('updateResource', () => {
    it('applies the update optimistically, then reconciles with the server response', async () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'))
      apiMock.mockResolvedValue(resource('r1', { title: 'Updated', type: 'book' }))

      const result = await store.updateResource('r1', { title: 'Updated', type: 'book' })

      expect(apiMock).toHaveBeenCalledWith('/resources/r1', {
        method: 'PUT',
        body: JSON.stringify({ title: 'Updated', type: 'book' })
      })
      expect(store.resources[0]).toEqual(resource('r1', { title: 'Updated', type: 'book' }))
      expect(result).toEqual(resource('r1', { title: 'Updated', type: 'book' }))
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Resource updated successfully',
        color: 'success'
      })
    })

    it('reverts to the original resource on failure', async () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'))
      apiMock.mockRejectedValue(new Error('conflict'))

      const result = await store.updateResource('r1', { title: 'Updated' })

      expect(store.resources[0]).toEqual(resource('r1'))
      expect(result).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to update resource',
        color: 'error'
      })
    })

    it('reports success but leaves state untouched for an unknown id', async () => {
      const store = useResourceStore()
      apiMock.mockResolvedValue(resource('r99'))

      await store.updateResource('r99', { title: 'Updated' })

      expect(store.resources).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Resource updated successfully',
        color: 'success'
      })
    })
  })

  describe('deleteResource', () => {
    it('removes the resource optimistically and succeeds', async () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'), resource('r2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteResource('r1')

      expect(apiMock).toHaveBeenCalledWith('/resources/r1', { method: 'DELETE' })
      expect(store.resources.map(r => r.id)).toEqual(['r2'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Task',
        description: 'Resource deleted successfully',
        color: 'success'
      })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useResourceStore()

      await store.deleteResource('r42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Resource not found',
        color: 'error'
      })
    })

    it('restores the resource when deletion fails', async () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteResource('r1')

      expect(store.resources.map(r => r.id)).toEqual(['r1'])
      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'Error',
        description: 'Failed to delete resource',
        color: 'error'
      })
    })
  })

  describe('filter helpers', () => {
    it('getResourcesByTopic returns only resources for the given topic', () => {
      const store = useResourceStore()
      store.resources.push(resource('r1'), resource('r2', { topic_id: 't2' }))

      expect(store.getResourcesByTopic('t1').map(r => r.id)).toEqual(['r1'])
      expect(store.getResourcesByTopic('t2').map(r => r.id)).toEqual(['r2'])
    })

    it('getResourcesByTask returns only resources for the given task', () => {
      const store = useResourceStore()
      store.resources.push(resource('r1', { task_id: 'task1' }), resource('r2'))

      expect(store.getResourcesByTask('task1').map(r => r.id)).toEqual(['r1'])
    })
  })

  describe('type helpers', () => {
    it('getResourceTypeIcon maps every type to an icon', () => {
      const store = useResourceStore()

      expect(store.getResourceTypeIcon('video')).toBe('🎥')
      expect(store.getResourceTypeIcon('article')).toBe('📄')
      expect(store.getResourceTypeIcon('document')).toBe('📋')
      expect(store.getResourceTypeIcon('book')).toBe('📚')
      expect(store.getResourceTypeIcon('course')).toBe('🎓')
    })

    it('getResourceTypeColor maps every type to a color class', () => {
      const store = useResourceStore()

      expect(store.getResourceTypeColor('video')).toBe('text-red-600 dark:text-red-400')
      expect(store.getResourceTypeColor('article')).toBe('text-blue-600 dark:text-blue-400')
      expect(store.getResourceTypeColor('document')).toBe('text-green-600 dark:text-green-400')
      expect(store.getResourceTypeColor('book')).toBe('text-purple-600 dark:text-purple-400')
      expect(store.getResourceTypeColor('course')).toBe('text-orange-600 dark:text-orange-400')
    })
  })
})
