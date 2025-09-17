<script setup lang="ts">
definePageMeta({
  layout: false
})

const auth = useAuthStore()
const error = ref('')
const { schemas } = useValidation()

// Import FormInline component
import FormInline from '@/components/ui/FormInline.vue'

const fields = [
  {
    name: 'name',
    label: 'Name',
    type: 'text' as const,
    placeholder: 'Enter your full name',
    required: true
  },
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
    placeholder: 'Create a password',
    required: true
  }
]

const submit = async (formData: Record<string, any>) => {
  try {
    error.value = ''
    await auth.register(formData)
  } catch (err) {
    error.value = err?.message || 'Something went wrong.'
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background-light dark:bg-background-dark">
    <div class="bg-surface-light dark:bg-surface-dark p-8 rounded-2xl shadow-lg w-full max-w-md border border-surface-light dark:border-surface-dark">
      <h2 class="text-2xl font-bold mb-6 text-center text-text-light dark:text-text-dark">Register</h2>
      <FormInline
        :fields="fields"
        :validation-schema="schemas.auth.register"
        :loading="auth.loading"
        :error="error"
        submit-label="Register"
        @submit="submit"
      />

      <div class="flex items-center justify-center mt-4">
        <NuxtLink to="/login" class="text-primary hover:text-primary/80 underline text-sm">
          Already have an account?
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

