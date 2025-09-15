<script setup>
definePageMeta({
  layout: false
})

const form = reactive({ email: '', password: '' })
const error = ref('')
const validationErrors = ref({})
const authStore = useAuthStore()
const { validateObject, schemas } = useValidation()
const { isOnline } = useOffline()

const handleLogin = async () => {
  try {
    error.value = ''
    validationErrors.value = {}

    const validation = validateObject(form, schemas.auth.login)

    if (!validation.isValid) {
      validationErrors.value = validation.errors
      return
    }

    await authStore.login(form)
  } catch (err) {
    const errorMessage = err?.message || 'Something went wrong.'

    // Provide more user-friendly error messages
    if (errorMessage.includes('Failed to fetch') || errorMessage.includes('NetworkError')) {
      if (!isOnline.value) {
        error.value = 'No internet connection. Please check your network and try again.'
      } else {
        error.value = 'Unable to connect to the server. Please check your connection and try again.'
      }
    } else if (errorMessage.includes('timeout') || errorMessage.includes('TimeoutError')) {
      error.value = 'Connection timed out. Please check your network and try again.'
    } else if (errorMessage.includes('401') || errorMessage.includes('Invalid credentials')) {
      error.value = 'Invalid email or password. Please check your credentials and try again.'
    } else if (errorMessage.includes('429')) {
      error.value = 'Too many login attempts. Please wait a few minutes and try again.'
    } else if (errorMessage.includes('500') || errorMessage.includes('Server error')) {
      error.value = 'Server is temporarily unavailable. Please try again in a few minutes.'
    } else {
      error.value = errorMessage
    }
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Login</h2>
      </div>
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter your email"
            required
            :class="{ 'border-red-500 focus:ring-red-500': validationErrors.email }"
          />
          <p v-if="validationErrors.email" class="mt-1 text-sm text-red-500 dark:text-red-400">
            {{ validationErrors.email }}
          </p>
        </div>
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Password</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter your password"
            required
            :class="{ 'border-red-500 focus:ring-red-500': validationErrors.password }"
          />
          <p v-if="validationErrors.password" class="mt-1 text-sm text-red-500 dark:text-red-400">
            {{ validationErrors.password }}
          </p>
        </div>
        <div class="flex items-center justify-between">
           <NuxtLink to="/register" class="text-primary hover:text-primary/80 underline text-sm">
             Don't have an account?
           </NuxtLink>
           <NuxtLink to="/forgot-password" class="text-primary hover:text-primary/80 underline text-sm">
             Forgot Password?
           </NuxtLink>
         </div>
        <UiButton
          type="submit"
          variant="primary"
          size="md"
          class="w-full"
        >
          Log in
        </UiButton>
      </form>
      <p v-if="error" class="text-red-500 dark:text-red-400 mt-4 text-center text-sm">{{ error }}</p>
    </div>
  </div>
</template>

