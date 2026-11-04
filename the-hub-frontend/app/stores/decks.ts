import { defineStore } from 'pinia'
import { ref } from 'vue'



interface Deck {
  deck_id: string
  name: string
}

export interface DeckResponse {
  decks: Deck[]
}

export const useDeckStore = defineStore('deck', () => {
  const decks = ref<Deck[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const toast = useToast()

  async function fetchDecks() {
    const { $api } = useNuxtApp()
    loading.value = true
    fetchError.value = null
    try {
      const fetchedDecks = await $api<DeckResponse>('decks')
      if (fetchedDecks) {
        decks.value = fetchedDecks.decks
      }
    } catch (error) {
      fetchError.value = error as Error
      toast.add({ title: 'Failed to fetch decks', color: 'error' })
      console.error('Error fetching decks:', error)
    } finally {
      loading.value = false
    }
  }

  async function submitForm(payload: Omit<Deck, 'deck_id'>) {
    // Create optimistic deck
    const optimisticDeck: Deck = {
      deck_id: `temp-${Date.now()}`,
      name: payload.name
    }

    // Optimistically add to local state
    decks.value.push(optimisticDeck)

    try {
      const { $api } = useNuxtApp()
      const newDeck = await $api<Deck>('decks', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic deck with real data
      const optimisticIndex = decks.value.findIndex(d => d.deck_id === optimisticDeck.deck_id)
      if (optimisticIndex !== -1) {
        decks.value[optimisticIndex] = newDeck
      }

      toast.add({ title: "Deck added successfully", color: 'success' })
    } catch (err) {
      // Remove optimistic deck on error
      decks.value = decks.value.filter(d => d.deck_id !== optimisticDeck.deck_id)
      toast.add({ title: "Deck not added", color: 'error' })
      console.error('Error adding deck:', err)
    }
  }

  async function editDeck(payload: Deck) {
    // Store original deck for potential rollback
    const originalDeckIndex = decks.value.findIndex(d => d.deck_id === payload.deck_id)
    const originalDeck = originalDeckIndex !== -1 ? { ...decks.value[originalDeckIndex] } : null

    // Optimistically update the deck
    if (originalDeckIndex !== -1) {
      decks.value[originalDeckIndex] = { ...decks.value[originalDeckIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Deck>(`decks/${payload.deck_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload)
      })

      // Update with server response to ensure consistency
      if (originalDeckIndex !== -1 && data) {
        decks.value[originalDeckIndex] = data
      }

      toast.add({ title: "Edited deck succesfully", color: 'success' })
    } catch (err) {
      // Revert optimistic update on error
      if (originalDeck && originalDeckIndex !== -1) {
        decks.value[originalDeckIndex] = originalDeck
      }
      toast.add({ title: "Editing deck failed", color: 'error' })
    }
  }

  async function deleteDeck(id: string) {
    // Store the deck for potential rollback
    const deckToDelete = decks.value.find(d => d.deck_id === id)
    if (!deckToDelete) {
      toast.add({ title: "Deck not found", color: 'error' })
      return
    }

    // Optimistically remove from local state
    decks.value = decks.value.filter((d) => d.deck_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`decks/${id}`, {
        method: 'DELETE'
      })

      toast.add({ title: "Deck deleted successfully", color: 'success' })
    } catch (err) {
      // Restore the deck on error
      decks.value.push(deckToDelete)
      toast.add({ title: "Deck not deleted", color: 'error' })
      console.error('Error deleting deck:', err)
    }
  }

  function reset() {
    decks.value = []
  }

  interface FlashcardFromPDF {
    front: string
    back: string
    category: string
  }

  interface GenerateFlashcardsResponse {
    cards: FlashcardFromPDF[]
    deck_id: string
    deck_name: string
    message: string
  }

  const pdfToFlashcardLoading = ref(false)

  async function generateFlashcardsFromPDF(
    pdfBase64: string,
    numCards: number,
    deckId?: string,
    newDeckName?: string,
    instruction?: string
  ): Promise<GenerateFlashcardsResponse | null> {
    const { $api } = useNuxtApp()
    pdfToFlashcardLoading.value = true

    const body: Record<string, unknown> = {
      pdf: pdfBase64,
      num_cards: numCards,
    }

    if (deckId) {
      body.deck_id = deckId
    }
    if (newDeckName) {
      body.new_deck_name = newDeckName
    }
    if (instruction) {
      body.instruction = instruction
    }

    try {
      const response = await $api<GenerateFlashcardsResponse>('/flashcards/from-pdf', {
        method: 'POST',
        body: JSON.stringify(body),
      })

      if (response) {
        toast.add({ title: response.message, color: 'success' })
        await fetchDecks()
        return response
      }
      return null
    } catch (error) {
      console.error('Failed to generate flashcards from PDF:', error)
      toast.add({ title: 'Failed to generate flashcards', color: 'error' })
      return null
    } finally {
      pdfToFlashcardLoading.value = false
    }
  }

  return {
    decks,
    loading,
    fetchError,
    fetchDecks,
    editDeck,
    deleteDeck,
    submitForm,
    reset,
    pdfToFlashcardLoading,
    generateFlashcardsFromPDF,
  }
})
