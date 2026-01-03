<script setup>
import { onErrorCaptured } from 'vue'

const { addToast } = useToast()

onErrorCaptured((error, instance, info) => {
  console.error('Global error caught:', error, info)

  // Only show toast for unexpected errors, not for handled API errors
  if (!error?.message?.includes('Invalid email or password') &&
      !error?.message?.includes('Your session has expired')) {
    addToast('An unexpected error occurred. Please try again.', 'error')
  }

  // Return false to prevent the error from propagating further
  return false
})
</script>

<template>
  <slot />
</template>