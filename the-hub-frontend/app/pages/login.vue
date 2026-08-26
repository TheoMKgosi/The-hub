<script setup lang="ts">
import * as v from 'valibot'
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
definePageMeta({
  layout: false
})


const error = ref('')
const loading = ref(false)
const authStore = useAuthStore()

const state = reactive({
  email: '',
  password: ''
})


const schema = v.object({
  email: v.pipe(v.string(), v.email('Invalid email')),
  password: v.pipe(v.string(), v.minLength(6, 'Must be at least 6 characters'))
})

type Schema = v.InferOutput<typeof schema>

const handleLogin = async (formData: FormSubmitEvent<Schema>) => {
  try {
    error.value = ''
    await authStore.login(formData.data)
  } catch (err) {
    const errorMessage = err?.message || 'Something went wrong.'
    // Provide more user-friendly error messages
    if (errorMessage.includes('Failed to fetch') || errorMessage.includes('NetworkError')) {
      error.value = 'No internet connection. Please check your network and try again.'
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
  <div class="min-h-screen bg-background-light dark:bg-background-dark py-8">
    <div class="max-w-md mx-auto px-4">
      <Banner />
    </div>
    <div class="flex items-center justify-center">
      <div
        class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Login</h2>
        </div>
        <UForm :schema="schema" :state="state" class="space-y-4" @submit="handleLogin">
          <UFormField label="Email" name="email" require>
            <UInput v-model="state.email" placeholder="Enter your email" />
          </UFormField>
          <UFormField label="Password" name="password" require>
            <UInput v-model="state.password" placeholder="*******" type="password" />
          </UFormField>
          <UButton label="Login" :loading="loading" type="submit" color="secondary" />
        </UForm>

        <div class="flex items-center justify-between mt-4">
          <NuxtLink to="/register" class="text-neutral-400 hover:text-primary/80 underline text-sm">
            Don't have an account?
          </NuxtLink>
          <NuxtLink to="/forgot-password" class="text-neutral-400 hover:text-primary/80 underline text-sm">
            Forgot Password?
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
