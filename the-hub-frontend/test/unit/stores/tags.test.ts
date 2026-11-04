import { describe, it, expect, vi, beforeAll, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useTagStore } from '../../../app/stores/tags'

// The tag store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them all; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
})

beforeEach(() => {
  // Fresh pinia per test gives fresh store state
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
})

const tag = (id: number, name = `tag-${id}`) => ({ tag_id: id, name, color: '#ff0000' })

describe('useTagStore', () => {
  describe('fetchTags', () => {
    it('populates tags and clears loading', async () => {
      const store = useTagStore()
      apiMock.mockResolvedValue({ tags: [tag(1), tag(2)] })

      await store.fetchTags()

      expect(apiMock).toHaveBeenCalledWith('tags')
      expect(store.tags).toEqual([tag(1), tag(2)])
      expect(store.loading).toBe(false)
    })

    it('keeps existing tags when the API returns nothing', async () => {
      const store = useTagStore()
      store.tags.push(tag(1))
      apiMock.mockResolvedValue(null)

      await store.fetchTags()

      expect(store.tags).toEqual([tag(1)])
      expect(store.loading).toBe(false)
    })
  })

  describe('submitForm', () => {
    it('adds the tag optimistically, then replaces it with the server response', async () => {
      const store = useTagStore()
      apiMock.mockResolvedValue(tag(7, 'work'))

      await store.submitForm({ name: 'work', color: '#ff0000' })

      expect(apiMock).toHaveBeenCalledWith('tags', {
        method: 'POST',
        body: JSON.stringify({ name: 'work', color: '#ff0000' })
      })
      expect(store.tags).toHaveLength(1)
      expect(store.tags[0]).toEqual(tag(7, 'work'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Tag added succesfully', color: 'success' })
    })

    it('rolls back the optimistic tag on failure', async () => {
      const store = useTagStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ name: 'work', color: '#ff0000' })

      expect(store.tags).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Tag not added', color: 'error' })
    })
  })

  describe('editTag', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useTagStore()
      store.tags.push(tag(1, 'old-name'))
      apiMock.mockResolvedValue(tag(1, 'new-name'))

      await store.editTag(tag(1, 'new-name'))

      expect(apiMock).toHaveBeenCalledWith('tags/1', {
        method: 'PATCH',
        body: JSON.stringify(tag(1, 'new-name'))
      })
      expect(store.tags[0]).toEqual(tag(1, 'new-name'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Edited tag succesfully', color: 'success' })
    })

    it('reverts to the original tag on failure', async () => {
      const store = useTagStore()
      store.tags.push(tag(1, 'original'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editTag(tag(1, 'changed'))

      expect(store.tags[0]).toEqual(tag(1, 'original'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing tag failed', color: 'error' })
    })

    it('is a no-op for a tag that is not in the list', async () => {
      const store = useTagStore()
      apiMock.mockResolvedValue(tag(99, 'remote'))

      await store.editTag(tag(99, 'remote'))

      expect(store.tags).toHaveLength(0)
    })
  })

  describe('deleteTag', () => {
    it('removes the tag optimistically and succeeds', async () => {
      const store = useTagStore()
      store.tags.push(tag(1), tag(2))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTag(1)

      expect(apiMock).toHaveBeenCalledWith('tags/1', { method: 'DELETE' })
      expect(store.tags.map(t => t.tag_id)).toEqual([2])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Tag deleted succesfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useTagStore()

      await store.deleteTag(42)

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Tag not found', color: 'error' })
    })

    it('restores the tag when deletion fails', async () => {
      const store = useTagStore()
      store.tags.push(tag(1))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteTag(1)

      expect(store.tags.map(t => t.tag_id)).toEqual([1])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Tag deletion failed', color: 'error' })
    })
  })

  it('reset clears all tags', () => {
    const store = useTagStore()
    store.tags.push(tag(1))

    store.reset()

    expect(store.tags).toHaveLength(0)
  })
})
