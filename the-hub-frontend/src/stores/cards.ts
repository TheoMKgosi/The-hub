import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMyFetch } from '@/config/fetch'
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
  const loading = ref(false)
  const fetchError = ref<Error | null>(null)
  const { addToast } = useToast()

  async function fetchCards(deckID: number) {
    loading.value = true
    const { data, error } = await useMyFetch(`decks/cards/${deckID}`).json<CardResponse>()

    if (data.value) cards.value = data.value.cards
    fetchError.value = error.value

    loading.value = false
  }

  async function submitForm(formData: Card) {
    console.log(formData)
    const { data, error } = await useMyFetch('cards').post(formData).json()
    fetchError.value = error.value
    if (fetchError.value) {
      addToast("Card not added", "error")
    } else {
      cards.value.push(data.value)
      addToast("Card added succesfully", "success")
    }
  }

  async function editCard(card: Card) {
    const { error } = await useMyFetch(`cards/${card.card_id}`).patch(card).json()

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
    await useMyFetch(`cards/${id}`).delete().json()
    cards.value = cards.value.filter((t) => t.card_id !== id)
    addToast("Card deleted succesfully", "success")
  }

  return {
    cards,
    loading,
    fetchError,
    fetchCards,
    editCard,
    deleteCard,
    submitForm,
  }
})
