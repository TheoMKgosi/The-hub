<script setup lang="ts">
definePageMeta({
  layout: false
})

const error = ref('')
const authStore = useAuthStore()
const { schemas } = useValidation()
const { isOnline } = useOffline()

const fields = [
  {
    name: 'email',
    label: 'Email',
    type: 'email' as const,
    placeholder: 'Enter your email',
    required: true
  },
  {
    name: 'password',
    label: 'Password',
    type: 'password' as const,
    placeholder: 'Enter your password',
    required: true
  }
]

const handleLogin = async (formData: Record<string, any>) => {
  try {
    error.value = ''
    await authStore.login(formData)
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
      <UiFormInline
        :fields="fields"
        :validation-schema="schemas.auth.login"
        :loading="authStore.loading"
        :error="error"
        submit-label="Log in"
        @submit="handleLogin"
      />

      <div class="flex items-center justify-between mt-4">
        <NuxtLink to="/register" class="text-primary hover:text-primary/80 underline text-sm">
          Don't have an account?
        </NuxtLink>
        <NuxtLink to="/forgot-password" class="text-primary hover:text-primary/80 underline text-sm">
          Forgot Password?
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

