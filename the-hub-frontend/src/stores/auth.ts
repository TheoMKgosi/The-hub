import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'
import { useMyFetch } from '@/config/fetch'
import { ref } from 'vue'
import { useTopicStore } from './topics'
import { useTaskStore } from './tasks'
import { useTagStore } from './tags'
import { useGoalStore } from './goals'
import { useDeckStore } from './decks'
import { useCardStore } from './cards'
import { useIncomeStore } from './income'
import { useScheduleStore } from './schedule'
import { useBudgetStore, useCategoryStore } from './finance'

interface User {
  id: number
  name: string
  email: string
}

interface AuthResponse {
  token: string
  user: User
}

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref<string | null>(localStorage.getItem('token'))

  const register = async (payload: { name: string; email: string; password: string }) => {
    const { data, error } = await useMyFetch<AuthResponse>('/register')
      .post(payload)
      .json()

    if (error.value) throw new Error(error.value.message || 'Registration failed')

    token.value = data.value?.token || ''
    user.value = data.value?.user || null

    // Save to localStorage
    localStorage.setItem('token', token.value)
    localStorage.setItem('user', JSON.stringify(user.value))

    router.push('/dashboard')
  }

  const login = async (payload: { email: string; password: string }) => {
    const { data, error } = await useMyFetch<AuthResponse>('/login')
      .post(payload)
      .json()

    if (error.value) throw new Error(error.value.message || 'Login failed')

    token.value = data.value?.token || ''
    user.value = data.value?.user || null

    // Save to localStorage
    localStorage.setItem('token', token.value)
    localStorage.setItem('user', JSON.stringify(user.value))

    router.push('/dashboard')
  }

  const logout = () => {
    // TODO: Clean up resource
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    resetStores()
    router.push('/login')
  }

  return { user, token, login, register, logout }
})

// Used to clear stores when logging out
const resetStores = () => {
  const taskStore = useTaskStore()
  const goalStore = useGoalStore()
  const scheduleStore = useScheduleStore()
  const tagStore = useTagStore()
  const topicStore = useTopicStore()
  const deckStore = useDeckStore()
  const cardStore = useCardStore()
  const incomeStore = useIncomeStore()
  const budgetStore = useBudgetStore()
  const categoryStore = useCategoryStore()

  taskStore.reset()
  goalStore.reset()
  scheduleStore.reset()
  tagStore.reset()
  topicStore.reset()
  deckStore.reset
  cardStore.reset()
  incomeStore.reset()
  budgetStore.reset()
  categoryStore.reset()
}
