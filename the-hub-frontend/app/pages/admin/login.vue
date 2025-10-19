<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'guest'
})

const auth = useAuthStore()
const router = useRouter()
const { addToast } = useToast()

// Redirect if already logged in as admin
if (auth.isLoggedIn && auth.user?.role === 'admin') {
  await router.push('/admin')
}

const formData = ref({
  email: '',
  password: ''
})

const loading = ref(false)
const showPassword = ref(false)

const handleLogin = async () => {
  if (!formData.value.email || !formData.value.password) {
    addToast('Please fill in all fields', 'error')
    return
  }

  loading.value = true
  
  try {
    await auth.login({
      email: formData.value.email,
      password: formData.value.password
    })
    
    // Check if user is admin after login
    if (auth.user?.role !== 'admin') {
      await auth.logout()
      addToast('Access denied. Admin privileges required.', 'error')
      return
    }
    
    addToast('Welcome to admin panel!', 'success')
    await router.push('/admin')
    
  } catch (error: any) {
    addToast(error?.message || 'Login failed. Please check your credentials.', 'error')
  } finally {
    loading.value = false
  }
}

const togglePassword = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-background-light via-surface-light to-background-light dark:from-background-dark dark:via-surface-dark dark:to-background-dark flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <!-- Admin Logo and Header -->
      <div class="text-center mb-8">
        <div class="mx-auto mb-6 flex justify-center">
          <div class="relative">
            <img src="/logo.svg" alt="The Hub Admin" class="h-16 w-16" />
            <div class="absolute -bottom-1 -right-1 bg-red-500 text-white text-xs font-bold rounded-full w-6 h-6 flex items-center justify-center">
              A
            </div>
          </div>
        </div>
        <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-2">
          Admin Portal
        </h1>
        <p class="text-text-light/70 dark:text-text-dark/70">
          Sign in to access admin dashboard
        </p>
      </div>

      <!-- Login Form -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-xl shadow-lg border border-surface-light/50 dark:border-surface-dark/50 p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Email Field -->
          <div>
            <label for="email" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Email Address
            </label>
            <div class="relative">
              <input
                id="email"
                v-model="formData.email"
                type="email"
                autocomplete="email"
                required
                class="w-full px-4 py-3 bg-background-light dark:bg-background-dark border border-surface-light dark:border-surface-dark rounded-lg text-text-light dark:text-text-dark placeholder-text-light/50 dark:placeholder-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary-light dark:focus:ring-primary-dark focus:border-transparent"
                placeholder="admin@example.com"
              />
              <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-text-light/40 dark:text-text-dark/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="formData.password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                class="w-full px-4 py-3 bg-background-light dark:bg-background-dark border border-surface-light dark:border-surface-dark rounded-lg text-text-light dark:text-text-dark placeholder-text-light/50 dark:placeholder-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary-light dark:focus:ring-primary-dark focus:border-transparent pr-12"
                placeholder="Enter your password"
              />
              <button
                type="button"
                @click="togglePassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-text-light/40 dark:text-text-dark/40 hover:text-text-light dark:hover:text-text-dark"
              >
                <svg v-if="showPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
                <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <div>
            <UiButton
              type="submit"
              :loading="loading"
              disabled={loading}
              variant="primary"
              size="lg"
              class="w-full"
            >
              <span v-if="!loading">Sign In to Admin</span>
              <span v-else>Signing in...</span>
            </UiButton>
          </div>
        </form>

        <!-- Security Notice -->
        <div class="mt-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <div class="flex items-start">
            <svg class="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <div class="text-sm text-red-800 dark:text-red-200">
              <p class="font-semibold mb-1">Admin Access Only</p>
              <p>This portal is restricted to authorized administrators only. Unauthorized access attempts will be logged.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center space-y-2">
        <p class="text-sm text-text-light/60 dark:text-text-dark/60">
          <NuxtLink 
            to="/login" 
            class="text-primary-light dark:text-primary-dark hover:underline"
          >
            ← Back to regular login
          </NuxtLink>
        </p>
        <p class="text-xs text-text-light/50 dark:text-text-dark/50">
          © 2024 The Hub Admin Portal
        </p>
      </div>
    </div>
  </div>
</template>