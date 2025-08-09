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
  const token = ref<string | null>('')
  const isLoggedIn = computed(() => !!token.value)


  if (process.client) {
    user.value = JSON.parse(localStorage.getItem('user') || 'null')
    token.value = localStorage.getItem('token')
  }

  const register = async (payload: { name: string; email: string; password: string }) => {
    try {
      const { $api } = useNuxtApp()
      const { token: fetchedToken, user: fetchedUser } = await $api<AuthResponse>('/register', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      token.value = fetchedToken
      user.value = fetchedUser

      router.push('/dashboard')

    } catch (err) {
      console.log(err)

    }
  }

  const login = async (payload: { email: string; password: string }) => {
    try {
      const { $api } = useNuxtApp()
      const { token: fetchedToken, user: fetchedUser } = await $api<AuthResponse>('/login', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      token.value = fetchedToken
      user.value = fetchedUser

      router.push('/dashboard')

    } catch (err) {
      console.log(err)
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

