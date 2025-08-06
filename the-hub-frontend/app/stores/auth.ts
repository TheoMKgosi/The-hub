import { defineStore } from 'pinia'
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

  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  if (process.client) {
    user.value = JSON.parse(localStorage.getItem('user') || 'null')
    token.value = localStorage.getItem('token')
  }

  const register = async (payload: { name: string; email: string; password: string }) => {
    const { data, error } = await useMyFetch<AuthResponse>('/register').post(payload).json()

    if (error.value) throw new Error(error.value.message || 'Registration failed')

    token.value = data.value?.token || ''
    user.value = data.value?.user || null

    if (process.client) {
      localStorage.setItem('token', token.value)
      localStorage.setItem('user', JSON.stringify(user.value))
    }

    router.push('/dashboard')
  }

  const login = async (payload: { email: string; password: string }) => {
    const { data, error } = await useMyFetch<AuthResponse>('/login', {
      method: 'POST',
      body: payload
    })


    if (error.value) throw new Error(error.value.message || 'Login failed')

    token.value = data.value?.token || ''
    user.value = data.value?.user || null

    if (process.client) {
      localStorage.setItem('token', token.value)
      localStorage.setItem('user', JSON.stringify(user.value))
    }

    router.push('/dashboard')
  }

  const logout = () => {
    token.value = null
    user.value = null

    if (process.client) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }

    resetStores()
    router.push('/login')
  }

  return { user, token, login, register, logout }
})

const resetStores = () => {
  useTaskStore().reset()
  useGoalStore().reset()
  useScheduleStore().reset()
  useTagStore().reset()
  useTopicStore().reset()
  useDeckStore().reset()
  useCardStore().reset()
  useIncomeStore().reset()
  useBudgetStore().reset()
  useCategoryStore().reset()
}

