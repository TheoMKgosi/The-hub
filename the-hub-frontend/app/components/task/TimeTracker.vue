<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useTaskStore } from '@/stores/tasks'

interface Props {
  taskId: string
  taskTitle: string
}

const props = defineProps<Props>()

const taskStore = useTaskStore()
const isTracking = ref(false)
const elapsedTime = ref(0)
const startTime = ref<Date | null>(null)
const interval = ref<NodeJS.Timeout | null>(null)
const description = ref('')

const formatTime = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}`
}

const startTracking = async () => {
  try {
    await taskStore.startTimeTracking(props.taskId, description.value)
    isTracking.value = true
    startTime.value = new Date()
    elapsedTime.value = 0
    startInterval()
  } catch (error) {
    console.error('Failed to start time tracking:', error)
  }
}

const stopTracking = async () => {
  try {
    await taskStore.stopTimeTracking(props.taskId)
    isTracking.value = false
    stopInterval()
    elapsedTime.value = 0
    description.value = ''
  } catch (error) {
    console.error('Failed to stop time tracking:', error)
  }
}

const startInterval = () => {
  interval.value = setInterval(() => {
    if (startTime.value) {
      elapsedTime.value = Math.floor((new Date().getTime() - startTime.value.getTime()) / 60000)
    }
  }, 1000)
}

const stopInterval = () => {
  if (interval.value) {
    clearInterval(interval.value)
    interval.value = null
  }
}

onUnmounted(() => {
  stopInterval()
})
</script>

<template>
  <div class="bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 rounded-lg p-4">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-medium text-text-light dark:text-text-dark">Time Tracker</h3>
      <div class="text-xs text-text-light dark:text-text-dark/60">
        {{ props.taskTitle }}
      </div>
    </div>

    <div v-if="!isTracking" class="space-y-3">
      <input
        v-model="description"
        type="text"
        placeholder="What are you working on?"
        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-sm"
      />

      <UiButton
        @click="startTracking"
        variant="primary"
        size="sm"
        class="w-full"
        :disabled="!description.trim()"
      >
        <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z" clip-rule="evenodd" />
        </svg>
        Start Tracking
      </UiButton>
    </div>

    <div v-else class="space-y-3">
      <div class="flex items-center justify-center">
        <div class="text-2xl font-mono font-bold text-primary">
          {{ formatTime(elapsedTime) }}
        </div>
      </div>

      <div class="text-xs text-text-light dark:text-text-dark/60 text-center">
        {{ description }}
      </div>

      <UiButton
        @click="stopTracking"
        variant="danger"
        size="sm"
        class="w-full"
      >
        <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8 7a.75.75 0 00-1.5 0v4.5a.75.75 0 001.5 0V7z" clip-rule="evenodd" />
        </svg>
        Stop Tracking
      </UiButton>
    </div>
  </div>
</template>