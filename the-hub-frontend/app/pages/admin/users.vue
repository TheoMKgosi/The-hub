<script setup lang="ts">
definePageMeta({
  middleware: ['authenticated', 'admin']
})

const adminStore = useAdminStore()
const { users, loading } = storeToRefs(adminStore)
const { addToast } = useToast()

const searchQuery = ref('')
const selectedUser = ref<string | null>(null)
const showDeleteDialog = ref(false)
const showRoleDialog = ref(false)
const userToDelete = ref<string | null>(null)
const userToChangeRole = ref<string | null>(null)
const newRole = ref<'user' | 'admin'>('user')

onMounted(() => {
  adminStore.getAllUsers()
})

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user => 
    user.name.toLowerCase().includes(query) ||
    user.email.toLowerCase().includes(query) ||
    user.role.toLowerCase().includes(query)
  )
})

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const confirmDelete = (userId: string) => {
  userToDelete.value = userId
  showDeleteDialog.value = true
}

const confirmRoleChange = (userId: string, currentRole: string) => {
  userToChangeRole.value = userId
  newRole.value = currentRole === 'admin' ? 'user' : 'admin'
  showRoleDialog.value = true
}

const handleDelete = async () => {
  if (!userToDelete.value) return
  
  try {
    await adminStore.deleteUser(userToDelete.value)
    showDeleteDialog.value = false
    userToDelete.value = null
    addToast('User deleted successfully', 'success')
  } catch (error) {
    console.error('Delete failed:', error)
  }
}

const handleRoleChange = async () => {
  if (!userToChangeRole.value) return
  
  try {
    await adminStore.updateUserRole(userToChangeRole.value, newRole.value)
    showRoleDialog.value = false
    userToChangeRole.value = null
    addToast(`User role updated to ${newRole.value}`, 'success')
  } catch (error) {
    console.error('Role change failed:', error)
  }
}

const goToAdminDashboard = () => {
  $router.push('/admin')
}

const getRoleBadgeClass = (role: string) => {
  return role === 'admin' 
    ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
    : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
}
</script>

<template>
  <div class="min-h-screen bg-background-light dark:bg-background-dark">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-text-light dark:text-text-dark">
              User Management
            </h1>
            <p class="mt-2 text-text-light/70 dark:text-text-dark/70">
              Manage user accounts and permissions
            </p>
          </div>
          <div class="flex items-center space-x-3">
            <UiButton
              @click="adminStore.getAllUsers()"
              variant="secondary"
              size="sm"
              :loading="loading"
            >
              Refresh
            </UiButton>
            <UiButton
              @click="goToAdminDashboard"
              variant="secondary"
              size="sm"
            >
              ← Dashboard
            </UiButton>
          </div>
        </div>
      </div>

      <!-- Search -->
      <div class="mb-6">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search users by name, email, or role..."
            class="w-full px-4 py-2 pl-10 bg-surface-light dark:bg-surface-dark border border-surface-light/50 dark:border-surface-dark/50 rounded-lg text-text-light dark:text-text-dark placeholder-text-light/50 dark:placeholder-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary-light dark:focus:ring-primary-dark"
          />
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="h-5 w-5 text-text-light/40 dark:text-text-dark/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-light dark:border-primary-dark"></div>
      </div>

      <!-- Users Table -->
      <div v-else class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-sm border border-surface-light/50 dark:border-surface-dark/50 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-surface-light dark:divide-surface-dark">
            <thead class="bg-surface-light/50 dark:bg-surface-dark/50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-text-light/70 dark:text-text-dark/70 uppercase tracking-wider">
                  User
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-text-light/70 dark:text-text-dark/70 uppercase tracking-wider">
                  Role
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-text-light/70 dark:text-text-dark/70 uppercase tracking-wider">
                  Joined
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-text-light/70 dark:text-text-dark/70 uppercase tracking-wider">
                  Last Updated
                </th>
                <th class="px-6 py-3 text-right text-xs font-medium text-text-light/70 dark:text-text-dark/70 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-surface-light dark:divide-surface-dark">
              <tr v-for="user in filteredUsers" :key="user.user_id" class="hover:bg-surface-light/30 dark:hover:bg-surface-dark/30">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm font-medium text-text-light dark:text-text-dark">
                      {{ user.name }}
                    </div>
                    <div class="text-sm text-text-light/60 dark:text-text-dark/60">
                      {{ user.email }}
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="['inline-flex px-2 py-1 text-xs font-semibold rounded-full', getRoleBadgeClass(user.role)]">
                    {{ user.role }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light/70 dark:text-text-dark/70">
                  {{ formatDate(user.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-text-light/70 dark:text-text-dark/70">
                  {{ formatDate(user.updated_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div class="flex justify-end space-x-2">
                    <UiButton
                      @click="confirmRoleChange(user.user_id, user.role)"
                      variant="secondary"
                      size="sm"
                    >
                      {{ user.role === 'admin' ? 'Demote' : 'Promote' }}
                    </UiButton>
                    <UiButton
                      @click="confirmDelete(user.user_id)"
                      variant="danger"
                      size="sm"
                    >
                      Delete
                    </UiButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Empty State -->
        <div v-if="filteredUsers.length === 0" class="text-center py-12">
          <div class="text-text-light/50 dark:text-text-dark/50">
            <svg class="mx-auto h-12 w-12 text-text-light/30 dark:text-text-dark/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <h3 class="mt-2 text-sm font-medium text-text-light dark:text-text-dark">No users found</h3>
            <p class="mt-1 text-sm text-text-light/50 dark:text-text-dark/50">
              {{ searchQuery ? 'Try adjusting your search terms' : 'Get started by adding new users' }}
            </p>
          </div>
        </div>
      </div>

      <!-- Delete Confirmation Dialog -->
      <ConfirmDialog
        v-if="showDeleteDialog"
        title="Delete User"
        message="Are you sure you want to delete this user? This action cannot be undone."
        confirm-text="Delete"
        cancel-text="Cancel"
        @confirm="handleDelete"
        @cancel="showDeleteDialog = false"
      />

      <!-- Role Change Confirmation Dialog -->
      <ConfirmDialog
        v-if="showRoleDialog"
        :title="`${newRole === 'admin' ? 'Promote' : 'Demote'} User`"
        :message="`Are you sure you want to ${newRole === 'admin' ? 'promote' : 'demote'} this user to ${newRole}?`"
        :confirm-text="newRole === 'admin' ? 'Promote' : 'Demote'"
        cancel-text="Cancel"
        @confirm="handleRoleChange"
        @cancel="showRoleDialog = false"
      />
    </div>
  </div>
</template>