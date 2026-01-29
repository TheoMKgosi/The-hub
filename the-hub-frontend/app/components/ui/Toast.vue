<script setup>
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()
</script>

<template>
  <div class="fixed top-4 right-4 z-[99999999] flex flex-col space-y-3">
    <div v-for="toast in toasts" :key="toast.id" :class="[
      'flex items-center px-4 py-3 rounded-lg shadow-lg transition-all duration-300 backdrop-blur-sm',
      'bg-surface-light/90 dark:bg-surface-dark/90 border border-surface-light/20 dark:border-surface-dark/20',
      toast.type === 'success' && 'border-l-success bg-success/10 dark:bg-success/20',
      toast.type === 'error' && 'border-l-red-500 bg-red-50/90 dark:bg-red-900/20',
      toast.type === 'warning' && 'border-l-warning bg-warning/10 dark:bg-warning/20',
      toast.type === 'info' && 'border-l-secondary bg-secondary/10 dark:bg-secondary/20'
    ]">
      <div class="flex-1">
        <div v-if="toast.type === 'success'" class="text-success dark:text-success font-medium">✓ {{ toast.message }}</div>
        <div v-else-if="toast.type === 'error'" class="text-red-700 dark:text-red-300 font-medium">✗ {{ toast.message }}</div>
        <div v-else-if="toast.type === 'warning'" class="text-warning dark:text-warning font-medium">⚠ {{ toast.message }}</div>
        <div v-else class="text-secondary dark:text-secondary font-medium">ℹ {{ toast.message }}</div>
      </div>
      <button @click="removeToast(toast.id)"
        class="ml-4 text-text-light dark:text-text-dark hover:text-text-light/70 dark:hover:text-text-dark/70 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-primary rounded p-1"
        aria-label="Close">
        ×
      </button>
    </div>
  </div>
</template>
