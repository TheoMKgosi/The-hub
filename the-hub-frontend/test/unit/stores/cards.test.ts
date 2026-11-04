import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useCardStore } from '../../../app/stores/cards'

// cards.ts explicitly imports useAuthStore from './auth' (a real runtime
// import). Mock that module so loading cards.ts doesn't pull in the whole
// auth store; exportCards only reads getAccessToken().
vi.mock('../../../app/stores/auth', () => ({
  useAuthStore: () => ({ getAccessToken: vi.fn().mockResolvedValue('test-token') })
}))

// The card store reaches for Nuxt auto-imports (useNuxtApp, useToast).
// Stub them all; the $api mock stands in for the API plugin.
const apiMock = vi.fn()
const toastAddMock = vi.fn()
const fetchMock = vi.fn()

// Plain (non-mock) DOM stubs for exportCards: it appends an <a> and clicks it.
const createdAnchors: Array<{ href: string; download: string; clicks: number }> = []

beforeAll(() => {
  vi.stubGlobal('useNuxtApp', () => ({ $api: apiMock }))
  vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
  vi.stubGlobal('fetch', fetchMock)
  // exportCards builds a download link through the DOM
  vi.stubGlobal('window', {
    URL: {
      createObjectURL: () => 'blob:export',
      revokeObjectURL: () => {}
    }
  })
  vi.stubGlobal('document', {
    createElement: () => {
      const anchor = { href: '', download: '', clicks: 0, click() { this.clicks++ } }
      createdAnchors.push(anchor)
      return anchor
    },
    body: {
      appendChild: () => {},
      removeChild: () => {}
    }
  })
})

beforeEach(() => {
  setActivePinia(createPinia())
  apiMock.mockReset()
  toastAddMock.mockClear()
  fetchMock.mockReset()
  createdAnchors.length = 0
  apiMock.defaults = { baseURL: 'https://api.test' }
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const card = (id: string, deckId = 'deck1', overrides: Record<string, unknown> = {}) => ({
  card_id: id,
  deck_id: deckId,
  question: `question-${id}`,
  answer: `answer-${id}`,
  ...overrides
})

describe('useCardStore', () => {
  it('starts with empty cards and review cards', () => {
    const store = useCardStore()

    expect(store.cards).toEqual([])
    expect(store.reviewCards).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.fetchError).toBeNull()
  })

  describe('fetchCards', () => {
    it('populates cards for a deck and clears loading', async () => {
      const store = useCardStore()
      apiMock.mockResolvedValue({ cards: [card('c1'), card('c2')] })

      await store.fetchCards('deck1')

      expect(apiMock).toHaveBeenCalledWith('/decks/cards/deck1')
      expect(store.cards).toEqual([card('c1'), card('c2')])
      expect(store.loading).toBe(false)
      expect(store.fetchError).toBeNull()
    })

    it('sets fetchError, shows a toast and clears loading on failure', async () => {
      const store = useCardStore()
      const failure = new Error('network down')
      apiMock.mockRejectedValue(failure)

      await store.fetchCards('deck1')

      expect(store.fetchError).toBe(failure)
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch cards', color: 'error' })
    })
  })

  describe('fetchDueCards', () => {
    it('populates reviewCards for a deck', async () => {
      const store = useCardStore()
      apiMock.mockResolvedValue({ cards: [card('c1')] })

      await store.fetchDueCards('deck1')

      expect(apiMock).toHaveBeenCalledWith('/cards/due/deck1')
      expect(store.reviewCards).toEqual([card('c1')])
      expect(store.cards).toEqual([])
      expect(store.loading).toBe(false)
    })

    it('shows an error toast on failure', async () => {
      const store = useCardStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.fetchDueCards('deck1')

      expect(store.fetchError).toBeInstanceOf(Error)
      expect(store.reviewCards).toEqual([])
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch due cards', color: 'error' })
    })
  })

  describe('submitForm', () => {
    it('adds the card optimistically, then replaces it with the server response', async () => {
      const store = useCardStore()
      const payload = card('temp', 'deck1', { question: 'What?', answer: 'This.' })
      apiMock.mockResolvedValue(card('c9'))

      await store.submitForm('deck1', payload)

      expect(apiMock).toHaveBeenCalledWith('cards', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
      expect(store.cards).toHaveLength(1)
      expect(store.cards[0]).toEqual(card('c9'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card added successfully', color: 'success' })
    })

    it('rolls back the optimistic card on failure', async () => {
      const store = useCardStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm('deck1', card('temp'))

      expect(store.cards).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card not added', color: 'error' })
    })
  })

  describe('editCard', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useCardStore()
      store.cards.push(card('c1'))
      const edited = card('c1', 'deck1', { question: 'Edited?' })
      apiMock.mockResolvedValue(edited)

      await store.editCard(edited)

      expect(apiMock).toHaveBeenCalledWith('cards/c1', {
        method: 'PATCH',
        body: JSON.stringify(edited)
      })
      expect(store.cards[0]).toEqual(edited)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card edited successfully', color: 'success' })
    })

    it('reverts to the original card on failure', async () => {
      const store = useCardStore()
      store.cards.push(card('c1'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editCard(card('c1', 'deck1', { question: 'Edited?' }))

      expect(store.cards[0]).toEqual(card('c1'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing card failed', color: 'error' })
    })

    it('is a no-op for a card that is not in the list', async () => {
      const store = useCardStore()
      apiMock.mockResolvedValue(card('c99', 'deck2'))

      await store.editCard(card('c99', 'deck2'))

      expect(store.cards).toHaveLength(0)
    })
  })

  describe('deleteCard', () => {
    it('removes the card optimistically and succeeds', async () => {
      const store = useCardStore()
      store.cards.push(card('c1'), card('c2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteCard('deck1', 'c1')

      expect(apiMock).toHaveBeenCalledWith('cards/c1', { method: 'DELETE' })
      expect(store.cards.map(c => c.card_id)).toEqual(['c2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card deleted successfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useCardStore()

      await store.deleteCard('deck1', 'c42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card not found', color: 'error' })
    })

    it('restores the card when deletion fails', async () => {
      const store = useCardStore()
      store.cards.push(card('c1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteCard('deck1', 'c1')

      expect(store.cards.map(c => c.card_id)).toEqual(['c1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Card not deleted', color: 'error' })
    })
  })

  describe('reviewCard', () => {
    it('posts the review rating', async () => {
      const store = useCardStore()
      apiMock.mockResolvedValue(undefined)

      await store.reviewCard('c1', 5)

      expect(apiMock).toHaveBeenCalledWith('cards/review/c1', {
        method: 'POST',
        body: { quality: 5 }
      })
      expect(toastAddMock).not.toHaveBeenCalled()
    })

    it('shows an error toast on failure', async () => {
      const store = useCardStore()
      apiMock.mockRejectedValue(new Error('server error'))

      await store.reviewCard('c1', 5)

      expect(toastAddMock).toHaveBeenCalledWith({
        title: 'An error occurred while reviewing card',
        color: 'error'
      })
    })
  })

  describe('exportCards', () => {
    it('downloads the export as a link click and shows a success toast', async () => {
      const store = useCardStore()
      fetchMock.mockResolvedValue({
        ok: true,
        statusText: 'OK',
        headers: {
          get: (name: string) => (name === 'content-disposition' ? 'attachment; filename="cards.json"' : null)
        },
        blob: async () => new Blob(['card data'])
      })

      await store.exportCards('deck1', 'json')

      expect(fetchMock).toHaveBeenCalledWith('https://api.test/decks/export/deck1/cards?format=json', {
        method: 'GET',
        headers: { Authorization: 'Bearer test-token' }
      })
      expect(createdAnchors).toHaveLength(1)
      expect(createdAnchors[0].download).toBe('cards.json')
      expect(createdAnchors[0].clicks).toBe(1)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Cards exported successfully', color: 'success' })
    })

    it('falls back to a default filename without a content-disposition header', async () => {
      const store = useCardStore()
      fetchMock.mockResolvedValue({
        ok: true,
        statusText: 'OK',
        headers: { get: () => null },
        blob: async () => new Blob(['card data'])
      })

      await store.exportCards('deck1', 'csv')

      expect(createdAnchors[0].download).toBe('cards_export.csv')
    })

    it('shows an error toast and rethrows when the download fails', async () => {
      const store = useCardStore()
      fetchMock.mockResolvedValue({
        ok: false,
        statusText: 'Forbidden',
        headers: { get: () => null },
        blob: async () => new Blob()
      })

      await expect(store.exportCards('deck1')).rejects.toThrow('Export failed: Forbidden')

      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Export failed', color: 'error' })
    })
  })

  describe('importCards', () => {
    it('posts the file as FormData and refreshes the cards list', async () => {
      const store = useCardStore()
      apiMock
        .mockResolvedValueOnce(undefined) // POST import
        .mockResolvedValueOnce({ cards: [card('c1')] }) // fetchCards refresh
      const file = new File(['cards'], 'cards.json', { type: 'application/json' })

      await store.importCards('deck1', file)

      expect(apiMock).toHaveBeenNthCalledWith(1, '/decks/import/deck1/cards?format=json', {
        method: 'POST',
        body: expect.any(FormData)
      })
      expect(apiMock).toHaveBeenNthCalledWith(2, '/decks/cards/deck1')
      expect(store.cards).toEqual([card('c1')])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Successfully imported cards' })
    })

    it('shows an error toast and rethrows on failure', async () => {
      const store = useCardStore()
      apiMock.mockRejectedValue(new Error('bad file'))
      const file = new File(['cards'], 'cards.csv')

      await expect(store.importCards('deck1', file, 'csv')).rejects.toThrow('bad file')

      expect(apiMock).toHaveBeenCalledWith('/decks/import/deck1/cards?format=csv', {
        method: 'POST',
        body: expect.any(FormData)
      })
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Import failed', color: 'error' })
    })
  })

  describe('reset', () => {
    it('clears all cards', () => {
      const store = useCardStore()
      store.cards.push(card('c1'))

      store.reset()

      expect(store.cards).toHaveLength(0)
    })
  })
})
