<script setup>
definePageMeta({
  layout: false
})

const form = reactive({ email: '' })
const error = ref('')
const success = ref('')
const validationErrors = ref({})
const isLoading = ref(false)
const { validateObject, schemas } = useValidation()

const authStore = useAuthStore()

const handleForgotPassword = async () => {
  try {
    error.value = ''
    success.value = ''
    validationErrors.value = {}
    isLoading.value = true

    const validation = validateObject(form, schemas.auth.forgotPassword)

    if (!validation.isValid) {
      validationErrors.value = validation.errors
      return
    }

    await authStore.forgotPassword(form)

    success.value = 'If an account with that email exists, a password reset link has been sent.'
    form.email = '' // Clear the form
  } catch (err) {
    error.value = err?.message || 'Something went wrong. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <h2 class="text-2xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Forgot Password</h2>
      <p class="text-sm text-text-light/70 dark:text-text-dark/70 mb-6 text-center">
        Enter your email address and we'll send you a link to reset your password.
      </p>

      <form @submit.prevent="handleForgotPassword" class="space-y-4">
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter your email"
            required
            :class="{ 'border-red-500 focus:ring-red-500': validationErrors.email }"
            :disabled="isLoading"
          />
          <p v-if="validationErrors.email" class="mt-1 text-sm text-red-500 dark:text-red-400">
            {{ validationErrors.email }}
          </p>
        </div>

        <UiButton
          type="submit"
          variant="primary"
          size="md"
          class="w-full"
          :disabled="isLoading"
        >
          <span v-if="isLoading">Sending...</span>
          <span v-else>Send Reset Link</span>
        </UiButton>
      </form>

      <div class="mt-6 text-center">
        <NuxtLink to="/login" class="text-primary hover:text-primary/80 underline text-sm">
          Back to Login
        </NuxtLink>
      </div>

      <p v-if="error" class="text-red-500 dark:text-red-400 mt-4 text-center text-sm">{{ error }}</p>
      <p v-if="success" class="text-green-600 dark:text-green-400 mt-4 text-center text-sm">{{ success }}</p>
    </div>
  </div>
</template>
