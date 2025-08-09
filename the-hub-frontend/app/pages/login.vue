<script setup>

const form = reactive({ email: '', password: '' })
const error = ref('')
const authStore = useAuthStore()

const handleLogin = async () => {
  try {
    error.value = ''
    await authStore.login(form)
  } catch (err) {
    error.value = err?.message || 'Something went wrong.'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-2xl shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block mb-1 font-medium">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="w-full p-3 rounded-xl border border-gray-300 focus:outline-none focus:ring"
            required
          />
        </div>
        <div class="mb-6">
          <label class="block mb-1 font-medium">Password</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full p-3 rounded-xl border border-gray-300 focus:outline-none focus:ring"
            required
          />
        </div>
        <NuxtLink to="/register" class="text-blue-600 underline text-sm block mb-4">
          Don’t have an account?
        </NuxtLink>
        <button
          type="submit"
          class="w-full bg-blue-600 text-white p-3 rounded-xl hover:bg-blue-700"
        >
          Log in
        </button>
      </form>
      <p v-if="error" class="text-red-500 mt-4 text-center">{{ error }}</p>
    </div>
  </div>
</template>

