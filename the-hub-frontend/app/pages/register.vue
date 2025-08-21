<script setup>
definePageMeta({
  layout: false
})

const auth = useAuthStore()
const form = reactive({ name: '', email: '', password: '' })
const error = ref('')

const submit = async () => {
  try {
    error.value = ''
    await auth.register(form)
  } catch (err) {
    error.value = err?.message || 'Something went wrong.'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <h2 class="text-2xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Register</h2>
      <form @submit.prevent="submit" class="space-y-4">
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Name</label>
          <input
            v-model="form.name"
            type="text"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter your full name"
            required
          />
        </div>
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Enter your email"
            required
          />
        </div>
        <div>
          <label class="block mb-2 font-medium text-text-light dark:text-text-dark">Password</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full px-3 py-2 rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
            placeholder="Create a password"
            required
          />
        </div>
        <div class="flex items-center justify-center">
          <NuxtLink to="/login" class="text-primary hover:text-primary/80 underline text-sm">
            Already have an account?
          </NuxtLink>
        </div>
        <UiButton
          type="submit"
          variant="primary"
          size="md"
          class="w-full"
        >
          Register
        </UiButton>
      </form>
      <p v-if="error" class="text-red-500 dark:text-red-400 mt-4 text-center text-sm">{{ error }}</p>
    </div>
  </div>
</template>

