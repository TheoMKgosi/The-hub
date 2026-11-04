import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useLearningPathStore } from '../../../app/stores/learning-paths'

// The learning-path store reaches for Nuxt auto-imports (useNuxtApp, useToast).
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
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const lpTopic = (id: string, status = 'active') => ({
  topic_id: id,
  title: `topic-${id}`,
  description: '',
  status
})

const learningPath = (id: string, topics: Array<ReturnType<typeof lpTopic>> = []) => ({
  learning_path_id: id,
  user_id: 'u1',
  title: `path-${id}`,
  description: '',
  topics
})

describe('useLearningPathStore', () => {
  it('starts with an empty list and no loading', () => {
    const store = useLearningPathStore()

    expect(store.learningPaths).toEqual([])
    expect(store.loading).toBe(false)
  })

  describe('createLearningPath', () => {
    it('prepends the learning path optimistically, then replaces it with the server response', async () => {
      const store = useLearningPathStore()
      store.learningPaths.push(learningPath('lp1'))
      const data = { title: 'New path', description: 'desc', topic_ids: ['t1', 't2'] }
      apiMock.mockResolvedValue(learningPath('lp9', [lpTopic('t1'), lpTopic('t2')]))

      const result = await store.createLearningPath(data)

      expect(apiMock).toHaveBeenCalledWith('/learning-paths', {
        method: 'POST',
        body: JSON.stringify(data)
      })
      // The new path is confirmed at the front, the pre-existing one stays behind
      expect(store.learningPaths).toHaveLength(2)
      expect(store.learningPaths[0]).toEqual(learningPath('lp9', [lpTopic('t1'), lpTopic('t2')]))
      expect(store.learningPaths[1]).toEqual(learningPath('lp1'))
      expect(result).toEqual(learningPath('lp9', [lpTopic('t1'), lpTopic('t2')]))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Learning path created successfully!', color: 'success' })
    })

    it('defaults the optimistic description to an empty string', async () => {
      const store = useLearningPathStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.createLearningPath({ title: 'New path', topic_ids: [] })

      // While the request is in flight, the optimistic path uses the payload title
      expect(store.learningPaths).toHaveLength(1)
      expect(store.learningPaths[0].title).toBe('New path')
      expect(store.learningPaths[0].description).toBe('')
      expect(store.learningPaths[0].topics).toEqual([])

      resolveApi(learningPath('lp9'))
      await pending
    })

    it('rolls back the optimistic path on failure and returns undefined', async () => {
      const store = useLearningPathStore()
      apiMock.mockRejectedValue(new Error('network down'))

      const result = await store.createLearningPath({ title: 'New path', topic_ids: [] })

      expect(store.learningPaths).toHaveLength(0)
      expect(result).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to create learning path', color: 'error' })
    })
  })

  describe('fetchLearningPaths', () => {
    it('populates learning paths and clears loading', async () => {
      const store = useLearningPathStore()
      apiMock.mockResolvedValue({ learning_paths: [learningPath('lp1'), learningPath('lp2')] })

      await store.fetchLearningPaths()

      expect(apiMock).toHaveBeenCalledWith('/learning-paths')
      expect(store.learningPaths).toEqual([learningPath('lp1'), learningPath('lp2')])
      expect(store.loading).toBe(false)
    })

    it('shows an error toast and clears loading on failure, keeping existing paths', async () => {
      const store = useLearningPathStore()
      store.learningPaths.push(learningPath('lp1'))
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchLearningPaths()

      expect(store.learningPaths).toEqual([learningPath('lp1')])
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch learning paths', color: 'error' })
    })
  })

  describe('deleteLearningPath', () => {
    it('removes the learning path optimistically and succeeds', async () => {
      const store = useLearningPathStore()
      store.learningPaths.push(learningPath('lp1'), learningPath('lp2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteLearningPath('lp1')

      expect(apiMock).toHaveBeenCalledWith('/learning-paths/lp1', { method: 'DELETE' })
      expect(store.learningPaths.map(lp => lp.learning_path_id)).toEqual(['lp2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Learning path deleted successfully!', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useLearningPathStore()

      await store.deleteLearningPath('lp42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Learning path not found', color: 'error' })
    })

    it('restores the learning path when deletion fails', async () => {
      const store = useLearningPathStore()
      store.learningPaths.push(learningPath('lp1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteLearningPath('lp1')

      expect(store.learningPaths.map(lp => lp.learning_path_id)).toEqual(['lp1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to delete learning path', color: 'error' })
    })
  })

  describe('getLearningPathProgress', () => {
    it('counts completed topics and rounds the percentage', () => {
      const store = useLearningPathStore()
      const path = learningPath('lp1', [
        lpTopic('t1', 'completed'),
        lpTopic('t2', 'completed'),
        lpTopic('t3', 'active')
      ])

      expect(store.getLearningPathProgress(path)).toEqual({
        completed: 2,
        total: 3,
        percentage: 67
      })
    })

    it('returns 0 progress when nothing is completed', () => {
      const store = useLearningPathStore()
      const path = learningPath('lp1', [lpTopic('t1'), lpTopic('t2')])

      expect(store.getLearningPathProgress(path)).toEqual({
        completed: 0,
        total: 2,
        percentage: 0
      })
    })
  })
})
