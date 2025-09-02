<script setup>
definePageMeta({
  layout: false
})

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  token: '',
  password: '',
  confirmPassword: ''
})
const error = ref('')
const success = ref('')
const validationErrors = ref({})
const isLoading = ref(false)
const { validateObject, schemas } = useValidation()

// Get token from URL query parameter
const token = route.query.token
if (token && typeof token === 'string') {
  form.token = token
}

const handleResetPassword = async () => {
  try {
    error.value = ''
    success.value = ''
    validationErrors.value = {}
    isLoading.value = true

    const validation = validateObject(form, schemas.auth.resetPassword)

    if (!validation.isValid) {
      validationErrors.value = validation.errors
      return
    }

    if (form.password !== form.confirmPassword) {
      validationErrors.value.confirmPassword = 'Passwords do not match'
      return
    }

    await authStore.resetPassword({
      token: form.token,
      password: form.password
    })

    success.value = 'Password reset successfully! You can now log in with your new password.'

    // Redirect to login page after a delay
    setTimeout(() => {
      router.push('/login')
    }, 3000)

  } catch (err) {
    error.value = err?.message || 'Something went wrong. Please try again.'
  } finally {
    isLoading.value = false
  }
}

// Redirect if no token is provided
if (!token) {
  error.value = 'Invalid reset link. Please request a new password reset.'
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <h2 class="text-2xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Reset Password</h2>

      <div v-if="!form.token" class="text-center">
        <p class="text-red-500 dark:text-red-400 mb-4">Invalid reset link. Please request a new password reset.</p>
        <NuxtLink to="/forgot-password" class="text-primary hover:text-primary/80 underline">
          Request New Reset Link
        </NuxtLink>
      </div>

      <form v-else @submit.prevent="handleResetPassword" class="space-y-4">
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">New Password</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter new password"
            required
            :class="{ 'border-red-500 focus:ring-red-500': validationErrors.password }"
            :disabled="isLoading"
          />
          <p v-if="validationErrors.password" class="mt-1 text-sm text-red-500 dark:text-red-400">
            {{ validationErrors.password }}
          </p>
        </div>

        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Confirm New Password</label>
          <input
            v-model="form.confirmPassword"
            type="password"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Confirm new password"
            required
            :class="{ 'border-red-500 focus:ring-red-500': validationErrors.confirmPassword }"
            :disabled="isLoading"
          />
          <p v-if="validationErrors.confirmPassword" class="mt-1 text-sm text-red-500 dark:text-red-400">
            {{ validationErrors.confirmPassword }}
          </p>
        </div>

        <UiButton
          type="submit"
          variant="primary"
          size="md"
          class="w-full"
          :disabled="isLoading"
        >
          <span v-if="isLoading">Resetting...</span>
          <span v-else>Reset Password</span>
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