import { defineStore } from 'pinia'
import { useToast } from '@/composables/useToast'

interface AdminUser {
  user_id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

interface AdminStats {
  total_users: number
  active_users: number
  total_tasks: number
  total_goals: number
  total_decks: number
  total_study_sessions: number
}

export const useAdminStore = defineStore('admin', () => {
  const { addToast } = useToast()
  
  const users = ref<AdminUser[]>([])
  const stats = ref<AdminStats | null>(null)
  const loading = ref(false)

  const getAllUsers = async () => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<{ users: AdminUser[] }>('/admin/users')
      users.value = response.users
    } catch (error) {
      addToast('Failed to fetch users', 'error')
      console.error('Error fetching users:', error)
    } finally {
      loading.value = false
    }
  }

  const getSystemStats = async () => {
    loading.value = true
    try {
      const { $api } = useNuxtApp()
      const response = await $api<AdminStats>('/admin/stats')
      stats.value = response
    } catch (error) {
      addToast('Failed to fetch system statistics', 'error')
      console.error('Error fetching stats:', error)
    } finally {
      loading.value = false
    }
  }

  const updateUserRole = async (userId: string, role: string) => {
    try {
      const { $api } = useNuxtApp()
      const response = await $api<AdminUser>(`/admin/users/${userId}/role`, {
        method: 'PUT',
        body: JSON.stringify({ role })
      })
      
      // Update user in local state
      const userIndex = users.value.findIndex(u => u.user_id === userId)
      if (userIndex !== -1) {
        users.value[userIndex] = response
      }
      
      addToast(`User role updated to ${role}`, 'success')
      return response
    } catch (error) {
      addToast('Failed to update user role', 'error')
      console.error('Error updating user role:', error)
      throw error
    }
  }

  const promoteToAdmin = async (userId: string) => {
    try {
      const { $api } = useNuxtApp()
      const response = await $api<AdminUser>(`/admin/users/${userId}/promote`, {
        method: 'POST'
      })
      
      // Update user in local state
      const userIndex = users.value.findIndex(u => u.user_id === userId)
      if (userIndex !== -1) {
        users.value[userIndex] = response
      }
      
      addToast('User promoted to admin successfully', 'success')
      return response
    } catch (error) {
      addToast('Failed to promote user to admin', 'error')
      console.error('Error promoting user:', error)
      throw error
    }
  }

  const deleteUser = async (userId: string) => {
    try {
      const { $api } = useNuxtApp()
      await $api(`/admin/users/${userId}`, {
        method: 'DELETE'
      })
      
      // Remove user from local state
      users.value = users.value.filter(u => u.user_id !== userId)
      
      addToast('User deleted successfully', 'success')
    } catch (error) {
      addToast('Failed to delete user', 'error')
      console.error('Error deleting user:', error)
      throw error
    }
  }

  return {
    users,
    stats,
    loading,
    getAllUsers,
    getSystemStats,
    updateUserRole,
    promoteToAdmin,
    deleteUser
  }
})