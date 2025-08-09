import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'



interface Deck {
  deck_id: number
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
    const fetchedDecks = await $api<DeckResponse>('decks')

    if (fetchedDecks) decks.value = fetchedDecks.decks

    loading.value = false
  }

  // TODO: Change payload parameters to object
  async function submitForm(payload: Deck) {
    try {
      const { $api } = useNuxtApp()
      const data = await $api<Deck>('decks', {
        method: 'POST',
        body: payload
      })
      decks.value.push(data)
      addToast("Deck added succesfully", "success")

    } catch (err) {
      addToast("Deck not added", "error")
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

  async function deleteDeck(id: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`decks/${id}`, {
        method: 'DELETE'
      })
      decks.value = decks.value.filter((t) => t.deck_id !== id)
      addToast("Deck deleted succesfully", "success")

    } catch (err) {
      addToast("Deck did not delete", "error")
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
