<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'
definePageMeta({
  layout: false
})

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

const handleLogin = async (event: FormSubmitEvent<Schema>) => {
  await authStore.login(event.data)
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
        <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
        <UForm :schema="schema" :state="state" class="space-y-3" @submit="handleLogin">
          <UFormField label="Email" name="email" require>
            <UInput v-model="state.email" placeholder="Enter your email" class="w-full" />
          </UFormField>
          <UFormField label="Password" name="password" require>
            <UInput v-model="state.password" placeholder="*******" type="password" class="w-full" />
          </UFormField>
          <UButton label="Login" type="submit" />
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
