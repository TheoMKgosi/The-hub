import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useTopicStore } from '../../../app/stores/topics'

// The topic store reaches for Nuxt auto-imports (useNuxtApp, useToast).
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

afterEach(() => {
  vi.restoreAllMocks()
})

const topic = (id: number, overrides: Record<string, unknown> = {}) => ({
  topic_id: id,
  title: `topic-${id}`,
  description: '',
  status: 'active',
  deadline: null,
  tags: [] as number[],
  ...overrides
})

describe('useTopicStore', () => {
  it('starts with no topics and no selection', () => {
    const store = useTopicStore()

    expect(store.topics).toEqual([])
    expect(store.topic).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.fetchError).toBeNull()
  })

  describe('fetchTopics', () => {
    it('populates topics and clears loading', async () => {
      const store = useTopicStore()
      apiMock.mockResolvedValue({ topics: [topic(1), topic(2)] })

      await store.fetchTopics()

      expect(apiMock).toHaveBeenCalledWith('topics')
      expect(store.topics).toEqual([topic(1), topic(2)])
      expect(store.loading).toBe(false)
    })

    it('keeps existing topics when the API returns nothing', async () => {
      const store = useTopicStore()
      store.topics.push(topic(1))
      apiMock.mockResolvedValue(null)

      await store.fetchTopics()

      expect(store.topics).toEqual([topic(1)])
      expect(store.loading).toBe(false)
    })

    it('NOTE: documents current behavior — fetchTopics has no error handling, so the rejection propagates and loading stays true', async () => {
      const store = useTopicStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await expect(store.fetchTopics()).rejects.toThrow('network down')

      expect(store.loading).toBe(true)
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('fetchTopic', () => {
    it('stores the fetched topic and clears loading', async () => {
      const store = useTopicStore()
      apiMock.mockResolvedValue(topic(5, { title: 'Single' }))

      await store.fetchTopic(5)

      expect(apiMock).toHaveBeenCalledWith('topics/5')
      expect(store.topic).toEqual(topic(5, { title: 'Single' }))
      expect(store.loading).toBe(false)
    })

    it('NOTE: documents current behavior — fetchTopic has no error handling, so the rejection propagates and loading stays true', async () => {
      const store = useTopicStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await expect(store.fetchTopic(5)).rejects.toThrow('network down')

      expect(store.loading).toBe(true)
      expect(toastAddMock).not.toHaveBeenCalled()
    })
  })

  describe('submitForm', () => {
    it('adds the topic optimistically, then replaces it with the server response', async () => {
      const store = useTopicStore()
      const payload = { title: 'New topic', description: 'desc', status: 'active', deadline: null, tags: [1, 2] }
      apiMock.mockResolvedValue(topic(7, { title: 'New topic' }))

      await store.submitForm(payload)

      expect(apiMock).toHaveBeenCalledWith('topics', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.topics).toHaveLength(1)
      expect(store.topics[0]).toEqual(topic(7, { title: 'New topic' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Topic added succesfully', color: 'success' })
    })

    it('rolls back the optimistic topic on failure', async () => {
      const store = useTopicStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ title: 'New topic', description: 'desc', status: 'active', deadline: null, tags: [] })

      expect(store.topics).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Topic not added', color: 'error' })
    })
  })

  describe('editTopic', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useTopicStore()
      store.topics.push(topic(1, { title: 'old-name' }))
      apiMock.mockResolvedValue(topic(1, { title: 'new-name' }))

      await store.editTopic(topic(1, { title: 'new-name' }))

      expect(apiMock).toHaveBeenCalledWith('topics/1', {
        method: 'PATCH',
        body: JSON.stringify(topic(1, { title: 'new-name' }))
      })
      expect(store.topics[0]).toEqual(topic(1, { title: 'new-name' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Edited topic succesfully', color: 'success' })
    })

    it('reverts to the original topic on failure', async () => {
      const store = useTopicStore()
      store.topics.push(topic(1, { title: 'original' }))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editTopic(topic(1, { title: 'changed' }))

      expect(store.topics[0]).toEqual(topic(1, { title: 'original' }))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing topic failed', color: 'error' })
    })

    it('is a no-op for a topic that is not in the list', async () => {
      const store = useTopicStore()
      apiMock.mockResolvedValue(topic(99))

      await store.editTopic(topic(99))

      expect(store.topics).toHaveLength(0)
    })
  })

  describe('deleteTopic', () => {
    it('removes the topic optimistically and succeeds', async () => {
      const store = useTopicStore()
      store.topics.push(topic(1), topic(2))
      apiMock.mockResolvedValue(undefined)

      await store.deleteTopic(1)

      expect(apiMock).toHaveBeenCalledWith('topics/1', { method: 'DELETE' })
      expect(store.topics.map(t => t.topic_id)).toEqual([2])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Topic deleted succesfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useTopicStore()

      await store.deleteTopic(42)

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Topic not found', color: 'error' })
    })

    it('restores the topic when deletion fails', async () => {
      const store = useTopicStore()
      store.topics.push(topic(1))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteTopic(1)

      expect(store.topics.map(t => t.topic_id)).toEqual([1])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Topic did not delete', color: 'error' })
    })
  })

  describe('reset', () => {
    it('clears all topics', () => {
      const store = useTopicStore()
      store.topics.push(topic(1))

      store.reset()

      expect(store.topics).toHaveLength(0)
    })
  })
})
