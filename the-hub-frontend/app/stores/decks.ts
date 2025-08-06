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
    loading.value = true
    const { data, error } = await useMyFetch('decks').json<DeckResponse>()

    if (data.value) decks.value = data.value.decks
    fetchError.value = error.value

    loading.value = false
  }

  async function submitForm(formData: Deck) {
    const { data, error } = await useMyFetch('decks').post(formData).json()
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Deck not added", "error")
    } else {
      decks.value.push(data.value)
      addToast("Deck added succesfully", "success")
    }
  }

  async function editDeck(deck: Deck) {
    const { error } = await useMyFetch(`deck/${deck.deck_id}`).patch(deck).json()

    if (!error.value) {
      const index = decks.value.findIndex(d => d.deck_id === deck.deck_id)
      if (index !== -1) {
        decks.value[index] = { ...decks.value[index], ...decks }
        addToast("Edited deck succesfully", "success")
      } else {
        addToast("Editing deck failed", "error")
      }
    } else {
      addToast("Editing deck failed", "error")
    }
  }

  async function deleteDeck(id: number) {
    await useMyFetch(`decks/${id}`).delete().json()
    decks.value = decks.value.filter((t) => t.deck_id !== id)
    addToast("Deck deleted succesfully", "success")
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
