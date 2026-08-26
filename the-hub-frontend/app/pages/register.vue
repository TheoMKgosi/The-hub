<script setup lang="ts">
definePageMeta({
  layout: false
})

const auth = useAuthStore()
const error = ref('')

const formData = reactive({
  name: '',
  email: '',
  password: ''
})

async function submit() {
  try {
    error.value = ''
    await auth.register(formData)
  } catch (err) {
    error.value = err?.message || 'Something went wrong.'
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
        <h2 class="text-2xl font-bold mb-6 text-center">Register</h2>
        <UForm title="Register" :state="formData" @sumbit="submit">
          <UFormField label="Name">
            <UInput v-model="formData.name" />
          </UFormField>
          <UFormField label="Email">
            <UInput v-model="formData.email" />
          </UFormField>
          <UFormField label="Password">
            <UInput v-model="formData.password" type="password" />
          </UFormField>
          <UButton label="Create Account" type="submit" class="w-full" />
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
