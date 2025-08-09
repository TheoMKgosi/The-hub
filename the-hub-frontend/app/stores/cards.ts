import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'



interface Card {
  card_id: number
  deck_id: number
  question: string
  answer: string
}

export interface CardResponse {
  cards: Card[]
}

export const useCardStore = defineStore('card', () => {
  const cards = ref<Card[]>([])
  const reviewCards = ref<Card[]>([])
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchCards(deckID: number) {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedCards = await $api<CardResponse>(`decks/cards/${deckID}`)
    if (fetchedCards) cards.value = fetchedCards.cards
    loading.value = false
  }

  async function fetchDueCards(deckID: number) {
    const { $api } = useNuxtApp()
    loading.value = true
    const { data, error } = await $api(`cards/due/${deckID}`).json<CardResponse>()

    if (data.value) reviewCards.value = data.value.cards
    fetchError.value = error.value

    loading.value = false
  }

  async function submitForm(formData: Card) {
    const { $api } = useNuxtApp()
    const { data, error } = await $api('cards').post(formData).json()
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Card not added", "error")
    } else {
      cards.value.push(data.value)
      addToast("Card added succesfully", "success")
    }
  }

  async function editCard(card: Card) {
    const { $api } = useNuxtApp()
    const { error } = await $api(`cards/${card.card_id}`).patch(card).json()

    if (!error.value) {
      const index = cards.value.findIndex(c => c.card_id === card.card_id)
      if (index !== -1) {
        cards.value[index] = { ...cards.value[index], ...cards }
        addToast("Edited card succesfully", "success")
      } else {
        addToast("Editing card failed", "error")
      }
    } else {
      addToast("Editing card failed", "error")
    }
  }

  async function deleteCard(id: number) {
    const { $api } = useNuxtApp()
    await $api(`cards/${id}`).delete().json()
    cards.value = cards.value.filter((t) => t.card_id !== id)
    addToast("Card deleted succesfully", "success")
  }

  async function reviewCard(cardID: number, rating: number) {
    const { $api } = useNuxtApp()
    await $api(`cards/review/${cardID}`).post({ quality: rating}).json()
  }

  function reset() {
    cards.value = []
  }

  return {
    cards,
    reviewCards,
    loading,
    fetchError,
    fetchCards,
    fetchDueCards,
    editCard,
    deleteCard,
    submitForm,
    reviewCard,
    reset,
  }
})
