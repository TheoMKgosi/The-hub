import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'



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
  const toast = useToast()

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
      toast.add({ title: 'Failed to fetch cards', color: 'error' })
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
      toast.add({ title: 'Failed to fetch due cards', color: 'error' })
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

      toast.add({ title: "Card added successfully", color: 'success' })
    } catch (err) {
      // Remove optimistic card on error
      cards.value = cards.value.filter(c => c.card_id !== optimisticCard.card_id)
      toast.add({ title: "Card not added", color: 'error' })
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

      toast.add({ title: "Card edited successfully", color: 'success' })
    } catch (error) {
      // Revert optimistic update on error
      if (originalCard && originalCardIndex !== -1) {
        cards.value[originalCardIndex] = originalCard
      }
      toast.add({ title: "Editing card failed", color: 'error' })
      console.error('Error editing card:', error)
    }
  }

  async function deleteCard(deckID: string, id: string) {
    // Store the card for potential rollback
    const cardToDelete = cards.value.find(c => c.card_id === id)
    if (!cardToDelete) {
      toast.add({ title: "Card not found", color: 'error' })
      return
    }

    // Optimistically remove from local state
    cards.value = cards.value.filter((card) => card.card_id !== id)

    try {
      const { $api } = useNuxtApp()
      await $api(`cards/${id}`, {
        method: 'DELETE'
      })

      toast.add({ title: "Card deleted successfully", color: 'success' })
    } catch (err) {
      // Restore the card on error
      cards.value.push(cardToDelete)
      toast.add({ title: "Card not deleted", color: 'error' })
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
      toast.add({ title: 'An error occurred while reviewing card', color: 'error' })
      console.error('Error reviewing card:', err)
    }
  }

  async function exportCards(deckID: string, format: 'json' | 'csv' = 'json') {
    try {
      const { $api } = useNuxtApp()
      const authStore = useAuthStore()
      const token = await authStore.getAccessToken()
      const response = await fetch(`${$api.defaults.baseURL}/decks/export/${deckID}/cards?format=${format}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Export failed: ${response.statusText}`)
      }

      // Get filename from content-disposition header
      const contentDisposition = response.headers.get('content-disposition')
      let filename = `cards_export.${format}`
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename=([^;]+)/)
        if (filenameMatch) {
          filename = filenameMatch[1].replace(/"/g, '')
        }
      }

      // Create blob and download
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)

      toast.add({ title: "Cards exported successfully", color: 'success' })
    } catch (error) {
      toast.add({ title: "Export failed", color: 'error' })
      console.error('Error exporting cards:', error)
      throw error
    }
  }

  async function importCards(deckID: string, file: File, format: 'json' | 'csv' = 'json') {
    const formData = new FormData()
    formData.append('file', file)

    try {
      const { $api } = useNuxtApp()
      await $api(`/decks/import/${deckID}/cards?format=${format}`, {
        method: 'POST',
        body: formData,
      })

      // Refresh cards list after successful import
      await fetchCards(deckID)

      toast.add({ title: 'Successfully imported cards' })

    } catch (error) {
      toast.add({ title: "Import failed", color: 'error' })
      console.error('Error importing cards:', error)
      throw error
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
    exportCards,
    importCards,
    reset,
  }
})
