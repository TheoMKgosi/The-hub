<script setup>
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { RouterLink } from 'vue-router'

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
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-2xl shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center">Register</h2>
      <form @submit.prevent="submit">
        <div class="mb-4">
          <label class="block mb-1 font-medium">Name</label>
          <input
            v-model="form.name"
            type="text"
            class="w-full p-3 rounded-xl border border-gray-300 focus:outline-none focus:ring"
            required
          />
        </div>
        <div class="mb-6">
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
        <RouterLink to="login">Already have an account</RouterLink>
        <button
          type="submit"
          class="w-full bg-blue-600 text-white p-3 rounded-xl hover:bg-blue-700"
        >
          Register
        </button>
      </form>
      <p v-if="error" class="text-red-500 mt-4 text-center">{{ error }}</p>
    </div>
  </div>
</template>

