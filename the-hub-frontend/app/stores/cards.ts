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
    fetchError.value = null
    try {
      const fetchedCards = await $api<CardResponse>(`/decks/cards/${deckID}`)
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

  async function fetchDueCards(deckID: string) {
    const { $api } = useNuxtApp()
    loading.value = true
    fetchError.value = null
    try {
      const fetchedDueCards = await $api<CardResponse>(`/cards/due/${deckID}`)
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

  async function submitForm(deckID: string, payload: Card) {
    // Create optimistic card
    const optimisticCard: Card = {
      card_id: `temp-${Date.now()}`,
      deck_id: payload.deck_id,
      question: payload.question,
      answer: payload.answer
    }

    // Optimistically add to local state
    cards.value.push(optimisticCard)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Card>('cards', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Replace optimistic card with real data
      const optimisticIndex = cards.value.findIndex(c => c.card_id === optimisticCard.card_id)
      if (optimisticIndex !== -1) {
        cards.value[optimisticIndex] = data
      }

      addToast("Card added successfully", "success")
    } catch (err) {
      // Remove optimistic card on error
      cards.value = cards.value.filter(c => c.card_id !== optimisticCard.card_id)
      addToast("Card not added", "error")
      console.error('Error adding card:', err)
    }
  }

  async function editCard(card: Card) {
    // Store original card for potential rollback
    const originalCardIndex = cards.value.findIndex(c => c.card_id === card.card_id)
    const originalCard = originalCardIndex !== -1 ? { ...cards.value[originalCardIndex] } : null

    // Optimistically update the card
    if (originalCardIndex !== -1) {
      cards.value[originalCardIndex] = { ...cards.value[originalCardIndex], ...card }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<Card>(`cards/${card.card_id}`, {
        method: 'PATCH',
        body: JSON.stringify(card)
      })

      // Update with server response to ensure consistency
      if (originalCardIndex !== -1 && data) {
        cards.value[originalCardIndex] = data
      }

      addToast("Card edited successfully", "success")
    } catch (error) {
      // Revert optimistic update on error
      if (originalCard && originalCardIndex !== -1) {
        cards.value[originalCardIndex] = originalCard
      }
      addToast("Editing card failed", "error")
      console.error('Error editing card:', error)
    }
  }

  async function deleteCard(deckID: string, id: string) {
    // Store the card for potential rollback
    const cardToDelete = cards.value.find(c => c.card_id === id)
    if (!cardToDelete) {
      addToast("Card not found", "error")
      return
    }

    // Optimistically remove from local state
    cards.value = cards.value.filter((card) => card.card_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`cards/${id}`, {
        method: 'DELETE'
      })

      addToast("Card deleted successfully", "success")
    } catch (err) {
      // Restore the card on error
      cards.value.push(cardToDelete)
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
