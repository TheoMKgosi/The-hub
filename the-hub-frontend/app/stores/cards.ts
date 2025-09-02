import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToast } from '@/composables/useToast'



interface Card {
  card_id: string
  deck_id: string
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

  async function fetchCards(deckID: string) {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedCards = await $api<CardResponse>(`decks/cards/${deckID}`)
    if (fetchedCards) cards.value = fetchedCards.cards
    loading.value = false
  }

  async function fetchDueCards(deckID: string) {
    const { $api } = useNuxtApp()
    loading.value = true
    const fetchedDueCards = await $api<CardResponse>(`cards/due/${deckID}`)

    if (fetchedDueCards) reviewCards.value = fetchedDueCards.cards

    loading.value = false
  }

  async function submitForm(deckID: string, payload: Card) {
    try {
      const { $api } = useNuxtApp()
      await $api<Card>('cards', {
        method: 'POST',
        body: payload
      })

      fetchCards(deckID)
      addToast("Card added succesfully", "success")
    } catch (err) {
      addToast("Card not added", "error")
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

  async function deleteCard(deckID: number, id: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`cards/${id}`, {
        method: 'DELETE'
      })
      cards.value = cards.value.filter((t) => t.card_id !== id)
      addToast("Card deleted succesfully", "success")
    } catch (err) {
      addToast("Card not deleted", "error")

    }
  }

  async function reviewCard(cardID: number, rating: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`cards/review/${cardID}`, {
        method: 'POST',
        body: JSON.stringify({ quality: rating })
      })
    } catch (err) {
      addToast('An error occured', 'error')
    }
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
