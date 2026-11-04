import { describe, it, expect, vi, beforeAll, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useDeckStore } from '../../../app/stores/decks'

// The deck store reaches for Nuxt auto-imports (useNuxtApp, useToast).
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

const deck = (id: string, name = `deck-${id}`) => ({ deck_id: id, name })

describe('useDeckStore', () => {
  it('starts with empty decks and no error', () => {
    const store = useDeckStore()

    expect(store.decks).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.fetchError).toBeNull()
    expect(store.pdfToFlashcardLoading).toBe(false)
  })

  describe('fetchDecks', () => {
    it('populates decks and clears loading', async () => {
      const store = useDeckStore()
      apiMock.mockResolvedValue({ decks: [deck('d1'), deck('d2')] })

      await store.fetchDecks()

      expect(apiMock).toHaveBeenCalledWith('decks')
      expect(store.decks).toEqual([deck('d1'), deck('d2')])
      expect(store.loading).toBe(false)
      expect(store.fetchError).toBeNull()
    })

    it('keeps existing decks when the API returns nothing', async () => {
      const store = useDeckStore()
      store.decks.push(deck('d1'))
      apiMock.mockResolvedValue(null)

      await store.fetchDecks()

      expect(store.decks).toEqual([deck('d1')])
      expect(store.loading).toBe(false)
    })

    it('sets fetchError, shows a toast and clears loading on failure', async () => {
      const store = useDeckStore()
      const failure = new Error('network down')
      apiMock.mockRejectedValue(failure)

      await store.fetchDecks()

      expect(store.fetchError).toBe(failure)
      expect(store.loading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to fetch decks', color: 'error' })
    })

    it('toggles loading while the request is in flight', async () => {
      const store = useDeckStore()
      let resolveApi!: (value: unknown) => void
      apiMock.mockReturnValue(new Promise((resolve) => { resolveApi = resolve }))

      const pending = store.fetchDecks()
      expect(store.loading).toBe(true)

      resolveApi({ decks: [] })
      await pending
      expect(store.loading).toBe(false)
    })
  })

  describe('submitForm', () => {
    it('adds the deck optimistically, then replaces it with the server response', async () => {
      const store = useDeckStore()
      apiMock.mockResolvedValue(deck('d9', 'work'))

      await store.submitForm({ name: 'work' })

      expect(apiMock).toHaveBeenCalledWith('decks', {
        method: 'POST',
        body: JSON.stringify({ name: 'work' })
      })
      expect(store.decks).toHaveLength(1)
      expect(store.decks[0]).toEqual(deck('d9', 'work'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Deck added successfully', color: 'success' })
    })

    it('rolls back the optimistic deck on failure', async () => {
      const store = useDeckStore()
      apiMock.mockRejectedValue(new Error('network down'))

      await store.submitForm({ name: 'work' })

      expect(store.decks).toHaveLength(0)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Deck not added', color: 'error' })
    })
  })

  describe('editDeck', () => {
    it('applies the edit optimistically, then reconciles with the server response', async () => {
      const store = useDeckStore()
      store.decks.push(deck('d1', 'old-name'))
      apiMock.mockResolvedValue(deck('d1', 'new-name'))

      await store.editDeck(deck('d1', 'new-name'))

      expect(apiMock).toHaveBeenCalledWith('decks/d1', {
        method: 'PATCH',
        body: JSON.stringify(deck('d1', 'new-name'))
      })
      expect(store.decks[0]).toEqual(deck('d1', 'new-name'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Edited deck succesfully', color: 'success' })
    })

    it('reverts to the original deck on failure', async () => {
      const store = useDeckStore()
      store.decks.push(deck('d1', 'original'))
      apiMock.mockRejectedValue(new Error('conflict'))

      await store.editDeck(deck('d1', 'changed'))

      expect(store.decks[0]).toEqual(deck('d1', 'original'))
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Editing deck failed', color: 'error' })
    })

    it('is a no-op for a deck that is not in the list', async () => {
      const store = useDeckStore()
      apiMock.mockResolvedValue(deck('d99', 'remote'))

      await store.editDeck(deck('d99', 'remote'))

      expect(store.decks).toHaveLength(0)
    })
  })

  describe('deleteDeck', () => {
    it('removes the deck optimistically and succeeds', async () => {
      const store = useDeckStore()
      store.decks.push(deck('d1'), deck('d2'))
      apiMock.mockResolvedValue(undefined)

      await store.deleteDeck('d1')

      expect(apiMock).toHaveBeenCalledWith('decks/d1', { method: 'DELETE' })
      expect(store.decks.map(d => d.deck_id)).toEqual(['d2'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Deck deleted successfully', color: 'success' })
    })

    it('shows a toast and does not call the API for an unknown id', async () => {
      const store = useDeckStore()

      await store.deleteDeck('d42')

      expect(apiMock).not.toHaveBeenCalled()
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Deck not found', color: 'error' })
    })

    it('restores the deck when deletion fails', async () => {
      const store = useDeckStore()
      store.decks.push(deck('d1'))
      apiMock.mockRejectedValue(new Error('server error'))

      await store.deleteDeck('d1')

      expect(store.decks.map(d => d.deck_id)).toEqual(['d1'])
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Deck not deleted', color: 'error' })
    })
  })

  describe('reset', () => {
    it('clears all decks', () => {
      const store = useDeckStore()
      store.decks.push(deck('d1'))

      store.reset()

      expect(store.decks).toHaveLength(0)
    })
  })

  describe('generateFlashcardsFromPDF', () => {
    it('posts the pdf payload with optional fields and refreshes decks', async () => {
      const store = useDeckStore()
      const response = {
        cards: [{ front: 'Q', back: 'A', category: 'definition' }],
        deck_id: 'd1',
        deck_name: 'Physics',
        message: 'Generated 1 flashcard'
      }
      apiMock
        .mockResolvedValueOnce(response) // /flashcards/from-pdf
        .mockResolvedValueOnce({ decks: [deck('d1', 'Physics')] }) // fetchDecks refresh

      const result = await store.generateFlashcardsFromPDF('base64-pdf', 5, 'd1', undefined, 'focus on math')

      expect(apiMock).toHaveBeenNthCalledWith(1, '/flashcards/from-pdf', {
        method: 'POST',
        body: JSON.stringify({
          pdf: 'base64-pdf',
          num_cards: 5,
          deck_id: 'd1',
          instruction: 'focus on math'
        })
      })
      expect(apiMock).toHaveBeenCalledTimes(2)
      expect(result).toEqual(response)
      expect(store.pdfToFlashcardLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Generated 1 flashcard', color: 'success' })
    })

    it('omits optional fields when they are not provided', async () => {
      const store = useDeckStore()
      apiMock
        .mockResolvedValueOnce({ cards: [], deck_id: 'd2', deck_name: 'New', message: 'ok' })
        .mockResolvedValueOnce({ decks: [] })

      await store.generateFlashcardsFromPDF('base64-pdf', 10, undefined, 'New deck')

      expect(apiMock).toHaveBeenNthCalledWith(1, '/flashcards/from-pdf', {
        method: 'POST',
        body: JSON.stringify({ pdf: 'base64-pdf', num_cards: 10, new_deck_name: 'New deck' })
      })
    })

    it('returns null and shows an error toast on failure', async () => {
      const store = useDeckStore()
      apiMock.mockRejectedValue(new Error('bad pdf'))

      const result = await store.generateFlashcardsFromPDF('base64-pdf', 5)

      expect(result).toBeNull()
      expect(store.pdfToFlashcardLoading).toBe(false)
      expect(toastAddMock).toHaveBeenCalledWith({ title: 'Failed to generate flashcards', color: 'error' })
    })
  })
})
