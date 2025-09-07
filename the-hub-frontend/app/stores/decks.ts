import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'



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
  const { addToast } = useToast()

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
      addToast('Failed to fetch decks', 'error')
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

      addToast("Deck added successfully", "success")
    } catch (err) {
      // Remove optimistic deck on error
      decks.value = decks.value.filter(d => d.deck_id !== optimisticDeck.deck_id)
      addToast("Deck not added", "error")
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

      addToast("Edited deck succesfully", "success")
    } catch (err) {
      // Revert optimistic update on error
      if (originalDeck && originalDeckIndex !== -1) {
        decks.value[originalDeckIndex] = originalDeck
      }
      addToast("Editing deck failed", "error")
    }
  }

  async function deleteDeck(id: string) {
    // Store the deck for potential rollback
    const deckToDelete = decks.value.find(d => d.deck_id === id)
    if (!deckToDelete) {
      addToast("Deck not found", "error")
      return
    }

    // Optimistically remove from local state
    decks.value = decks.value.filter((d) => d.deck_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`decks/${id}`, {
        method: 'DELETE'
      })

      addToast("Deck deleted successfully", "success")
    } catch (err) {
      // Restore the deck on error
      decks.value.push(deckToDelete)
      addToast("Deck not deleted", "error")
      console.error('Error deleting deck:', err)
    }
  }

  function reset() {
    decks.value = []
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
  }
})
