import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'
import { useMyFetch } from '@/config/fetch'
import { ref } from 'vue'

interface User {
  id: number
  name: string
  email: string
  // add more fields if needed
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
    router.push('/login')
  }

  return { user, token, login, register, logout }
})

