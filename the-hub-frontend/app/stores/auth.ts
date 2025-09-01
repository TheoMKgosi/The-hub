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
import { useValidation } from '@/composables/useValidation'

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
  const token = ref<string | null>('')
  const isLoggedIn = computed(() => !!token.value)


  if (process.client) {
    user.value = JSON.parse(localStorage.getItem('user') || 'null')
    token.value = localStorage.getItem('token')
  }

  const register = async (payload: { name: string; email: string; password: string }) => {
    try {
      const { validateObject, schemas } = useValidation()
      const validation = validateObject(payload, schemas.auth.register)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      const { token: fetchedToken, user: fetchedUser } = await $api<AuthResponse>('/register', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      token.value = fetchedToken
      user.value = fetchedUser

      router.push('/dashboard')
      addToast('Account created successfully!', 'success')

    } catch (err) {
      console.log(err)
      addToast(err?.message || 'Registration failed. Please try again.', 'error')
    }
  }

  const login = async (payload: { email: string; password: string }) => {
    try {
      const { validateObject, schemas } = useValidation()
      const validation = validateObject(payload, schemas.auth.login)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      const { token: fetchedToken, user: fetchedUser } = await $api<AuthResponse>('/login', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      token.value = fetchedToken
      user.value = fetchedUser

      router.push('/dashboard')
      addToast('Welcome back!', 'success')

    } catch (err) {
      console.log(err)
      throw err
    }
  }

  const logout = () => {
    token.value = null
    user.value = null

    localStorage.removeItem('token')
    localStorage.removeItem('user')

    resetStores()
    router.push('/login')
  }

  return { user, token, isLoggedIn, login, register, logout }
}, {
  persist: {
    storage: piniaPluginPersistedstate.localStorage()
  }
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

