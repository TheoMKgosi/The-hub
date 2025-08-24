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

  async function fetchCards(deckID: number) {
    const { $api } = useNuxtApp()
    loading.value = true
    fetchError.value = null
    try {
      const fetchedCards = await $api<CardResponse>(`decks/cards/${deckID}`)
      if (fetchedCards) {
        cards.value = fetchedCards.cards
      }
    } catch (error) {
      fetchError.value = error as Error
      addToast('Failed to fetch cards', 'error')
      console.error('Error fetching cards:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchDueCards(deckID: number) {
    const { $api } = useNuxtApp()
    loading.value = true
    fetchError.value = null
    try {
      const fetchedDueCards = await $api<CardResponse>(`cards/due/${deckID}`)
      if (fetchedDueCards) {
        reviewCards.value = fetchedDueCards.cards
      }
    } catch (error) {
      fetchError.value = error as Error
      addToast('Failed to fetch due cards', 'error')
      console.error('Error fetching due cards:', error)
    } finally {
      loading.value = false
    }
  }

  async function submitForm(payload: Card) {
    try {
      const { $api } = useNuxtApp()
      await $api<Card>('cards', {
        method: 'POST',
        body: payload
      })

      // Refresh cards for the deck
      fetchCards(payload.deck_id)
      addToast("Card added successfully", "success")
    } catch (err) {
      addToast("Card not added", "error")
      console.error('Error adding card:', err)
    }
  }

  async function editCard(card: Card) {
    try {
      const { $api } = useNuxtApp()
      await $api(`cards/${card.card_id}`, {
        method: 'PATCH',
        body: card
      })

      const index = cards.value.findIndex(c => c.card_id === card.card_id)
      if (index !== -1) {
        cards.value[index] = { ...cards.value[index], ...card }
        addToast("Card edited successfully", "success")
      } else {
        addToast("Card not found", "error")
      }
    } catch (error) {
      addToast("Editing card failed", "error")
      console.error('Error editing card:', error)
    }
  }

  async function deleteCard(deckID: string, id: string) {
    try {
      const { $api } = useNuxtApp()
      await $api(`cards/${id}`, {
        method: 'DELETE'
      })
      cards.value = cards.value.filter((card) => card.card_id !== id)
      addToast("Card deleted successfully", "success")
    } catch (err) {
      addToast("Card not deleted", "error")
      console.error('Error deleting card:', err)
    }
  }

  async function reviewCard(cardID: string, rating: number) {
    try {
      const { $api } = useNuxtApp()
      await $api(`cards/review/${cardID}`, {
        method: 'POST',
        body: { quality: rating }
      })
    } catch (err) {
      addToast('An error occurred while reviewing card', 'error')
      console.error('Error reviewing card:', err)
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
