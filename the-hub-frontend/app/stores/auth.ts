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
import { useToast } from "@/composables/useToast";

interface User {
  user_id: string
  name: string
  email: string
  role: string
}

interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()

  const { addToast } = useToast()

  const user = ref<User | null>(null)
  const accessToken = ref<string | null>('')
  const refreshToken = ref<string | null>('')
  const tokenExpiry = ref<number | null>(null)
  const isLoggedIn = computed(() => !!accessToken.value)

  // Token refresh state
  const isRefreshing = ref(false)
  const refreshPromise = ref<Promise<any> | null>(null)

  if (import.meta.client) {
    user.value = JSON.parse(localStorage.getItem('user') || 'null')
    accessToken.value = localStorage.getItem('accessToken')
    refreshToken.value = localStorage.getItem('refreshToken')
    const expiry = localStorage.getItem('tokenExpiry')
    tokenExpiry.value = expiry ? parseInt(expiry) : null
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
      const response = await $api<AuthResponse>('/register', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Store tokens and user data
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token
      tokenExpiry.value = Date.now() + (response.expires_in * 1000) // Convert to milliseconds
      user.value = response.user

      router.push('/dashboard')
      addToast('Account created successfully!', 'success')

    } catch (err) {
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
      const response = await $api<AuthResponse>('/login', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      // Store tokens and user data
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token
      tokenExpiry.value = Date.now() + (response.expires_in * 1000) // Convert to milliseconds
      user.value = response.user

      // Redirect based on user role
      if (response.user.role === 'admin') {
        router.push('/admin')
        addToast('Welcome to admin panel!', 'success')
      } else {
        router.push('/dashboard')
        addToast('Welcome back!', 'success')
      }

    } catch (err) {
      throw err
    }
  }

  const logout = async () => {
    try {
      // Attempt to revoke refresh token on server
      if (refreshToken.value) {
        const { $api } = useNuxtApp()
        await $api('/auth/logout', {
          method: 'POST',
          body: JSON.stringify({ refresh_token: refreshToken.value })
        })
      }
    } catch (error) {
      // Ignore logout errors - we still want to clear local state
      console.warn('Logout API call failed:', error)
    }

    // Clear all auth state
    accessToken.value = null
    refreshToken.value = null
    tokenExpiry.value = null
    user.value = null
    isRefreshing.value = false
    refreshPromise.value = null

    // Clear localStorage
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('tokenExpiry')
    localStorage.removeItem('user')

    resetStores()
    router.push('/login')
  }

  const forgotPassword = async (payload: { email: string }) => {
    try {
      const { validateObject, schemas } = useValidation()
      const validation = validateObject(payload, schemas.auth.forgotPassword)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      await $api('/auth/forgot-password', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      addToast('If an account with that email exists, a password reset link has been sent.', 'success')
      return { success: true }

    } catch (err) {
      addToast(err?.message || 'Failed to send reset email. Please try again.', 'error')
      throw err
    }
  }

  const resetPassword = async (payload: { token: string; password: string }) => {
    try {
      const { validateObject, schemas } = useValidation()
      const validation = validateObject(payload, schemas.auth.resetPassword)

      if (!validation.isValid) {
        const errorMessage = Object.values(validation.errors)[0]
        throw new Error(errorMessage)
      }

      const { $api } = useNuxtApp()
      await $api('/auth/reset-password', {
        method: 'POST',
        body: JSON.stringify(payload)
      })

      addToast('Password reset successfully!', 'success')
      return { success: true }

    } catch (err) {
      addToast(err?.message || 'Failed to reset password. Please try again.', 'error')
      throw err
    }
  }

  // Check if access token is expired or about to expire (within 5 minutes)
  const isTokenExpired = () => {
    if (!tokenExpiry.value) return true
    return Date.now() >= (tokenExpiry.value - (5 * 60 * 1000)) // 5 minutes buffer
  }

  // Refresh access token using refresh token
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    // If already refreshing, return the existing promise
    if (isRefreshing.value && refreshPromise.value) {
      return refreshPromise.value
    }

    isRefreshing.value = true

    try {
      refreshPromise.value = $fetch('/auth/refresh', {
        baseURL: useRuntimeConfig().public.apiBase,
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken.value })
      })

      const response = await refreshPromise.value

      // Update tokens
      accessToken.value = response.access_token
      tokenExpiry.value = Date.now() + (response.expires_in * 1000)

      return response
    } finally {
      isRefreshing.value = false
      refreshPromise.value = null
    }
  }

  // Get current access token, refreshing if necessary
  const getAccessToken = async () => {
    if (!accessToken.value) return null

    if (isTokenExpired()) {
      await refreshAccessToken()
    }

    return accessToken.value
  }

  // Persist auth state to localStorage when it changes
  if (process.client) {
    watch([accessToken, refreshToken, tokenExpiry, user], () => {
      if (accessToken.value) {
        localStorage.setItem('accessToken', accessToken.value)
      } else {
        localStorage.removeItem('accessToken')
      }

      if (refreshToken.value) {
        localStorage.setItem('refreshToken', refreshToken.value)
      } else {
        localStorage.removeItem('refreshToken')
      }

      if (tokenExpiry.value) {
        localStorage.setItem('tokenExpiry', tokenExpiry.value.toString())
      } else {
        localStorage.removeItem('tokenExpiry')
      }

      if (user.value) {
        localStorage.setItem('user', JSON.stringify(user.value))
      } else {
        localStorage.removeItem('user')
      }
    }, { immediate: true })
  }

  return {
    user,
    accessToken,
    refreshToken,
    isLoggedIn,
    login,
    register,
    logout,
    forgotPassword,
    resetPassword,
    refreshAccessToken,
    getAccessToken,
    isTokenExpired
  }
}, {
  persist: true
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

