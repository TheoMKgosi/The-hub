<script setup>
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()
</script>

<template>
  <div class="fixed top-4 right-4 z-50 flex flex-col space-y-3">
    <div v-for="toast in toasts" :key="toast.id" :class="[
      'flex items-center px-4 py-3 rounded shadow-lg transition-all duration-300',
      toast.type === 'success' && 'bg-green-100 text-green-800 border-l-4 border-green-500',
      toast.type === 'error' && 'bg-red-100 text-red-800 border-l-4 border-red-500',
      toast.type === 'warning' && 'bg-yellow-100 text-yellow-800 border-l-4 border-yellow-500',
      toast.type === 'info' && 'bg-blue-100 text-blue-800 border-l-4 border-blue-500'
    ]">
      <span class="flex-1">{{ toast.message }}</span>
      <button @click="removeToast(toast.id)" class="ml-4 text-xl font-bold focus:outline-none hover:text-gray-600"
        aria-label="Close">
        ×
      </button>
    </div>
  </div>
</template>
