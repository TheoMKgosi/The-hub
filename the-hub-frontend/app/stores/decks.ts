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

  // TODO: Change payload parameters to object
  async function submitForm(payload: Omit<Deck, 'deck_id'>) {
    try {
      const { $api } = useNuxtApp()
      const newDeck = await $api<Deck>('decks', {
        method: 'POST',
        body: payload
      })

      if (newDeck) {
        decks.value.push(newDeck)
        addToast("Deck added successfully", "success")
      }
    } catch (err) {
      addToast("Deck not added", "error")
      console.error('Error adding deck:', err)
    }
  }

  async function editDeck(payload: Deck) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api(`deck/${payload.deck_id}`, {
        method: 'PATCH',
        body: payload
      })

      // TODO: change implementation to locally filter decks instead my extra fetches
      fetchDecks()
      addToast("Edited deck succesfully", "success")
    } catch (err) {
      addToast("Editing deck failed", "error")
    }
  }

  async function deleteDeck(id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`decks/${id}`, {
        method: 'DELETE'
      })
      decks.value = decks.value.filter((d) => d.deck_id !== id)
      addToast("Deck deleted successfully", "success")
    } catch (err) {
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
