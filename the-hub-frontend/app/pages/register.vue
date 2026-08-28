<script setup lang="ts">
import * as v from 'valibot'
import type { FormSubmitEvent } from '@nuxt/ui'

definePageMeta({
  layout: false
})

const auth = useAuthStore()

const schema = v.object({
  name: v.pipe(v.string()),
  email: v.pipe(v.string(), v.email('Invalid email')),
  password: v.pipe(v.string(), v.minLength(6, 'Must be at least 6 characters'))
})

type Schema = v.InferOutput<typeof schema>


const formData = reactive({
  name: '',
  email: '',
  password: ''
})

async function onSubmit(event: FormSubmitEvent<Schema>) {
  await auth.register(event.data)
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
        <h2 class="text-2xl font-bold mb-6 text-center">Register</h2>
        <UForm :schema="schema" :state="formData" @submit="onSubmit" class="space-y-3">
          <UFormField label="Name" name="name">
            <UInput v-model="formData.name" class="w-full" />
          </UFormField>
          <UFormField label="Email" name="email">
            <UInput v-model="formData.email" class="w-full" />
          </UFormField>
          <UFormField label="Password" name="password">
            <UInput v-model="formData.password" type="password" class="w-full" />
          </UFormField>
          <UButton label="Create Account" type="submit" />
        </UForm>

        <div class="flex items-center justify-center mt-4">
          <NuxtLink to="/login" class="text-primary hover:text-primary/80 underline text-sm">
            Already have an account?
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
