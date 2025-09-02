<script setup>
definePageMeta({
  layout: false
})

const form = reactive({ email: '', password: '' })
const error = ref('')
const validationErrors = ref({})
const authStore = useAuthStore()
const { validateObject, schemas } = useValidation()

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
    error.value = err?.message || 'Something went wrong.'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <h2 class="text-2xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Login</h2>
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

