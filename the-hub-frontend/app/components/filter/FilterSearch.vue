<script setup lang="ts">
import SearchIcon from '../ui/svg/SearchIcon.vue'
import CrossIcon from '../ui/svg/CrossIcon.vue'

interface Props {
  modelValue: string
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Search tasks...'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const inputRef = ref<HTMLInputElement>()

// Debounced input
const debouncedValue = ref(props.modelValue)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const handleInput = (e: Event) => {
  const value = (e.target as HTMLInputElement).value
  debouncedValue.value = value
  
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  
  debounceTimer = setTimeout(() => {
    emit('update:modelValue', value)
  }, 300)
}

// Sync with parent
watch(() => props.modelValue, (newVal) => {
  debouncedValue.value = newVal
})

// Clear search
const clearSearch = () => {
  debouncedValue.value = ''
  emit('update:modelValue', '')
  inputRef.value?.focus()
}

// Focus on mount
onMounted(() => {
  inputRef.value?.focus()
})
</script>

<template>
  <div class="relative flex-1 min-w-[200px] max-w-[400px]">
    <!-- Search Icon -->
    <div class="absolute left-3 top-1/2 -translate-y-1/2 text-text-light/50 dark:text-text-dark/50">
      <SearchIcon class="w-5 h-5" />
    </div>
    
    <!-- Input -->
    <input
      ref="inputRef"
      :value="debouncedValue"
      @input="handleInput"
      :placeholder="placeholder"
      class="w-full pl-10 pr-10 py-2 text-sm rounded-md border transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
      :class="[
        debouncedValue
          ? 'bg-primary/5 border-primary/30 text-text-light dark:text-text-dark'
          : 'bg-surface-light dark:bg-surface-dark border-surface-light/30 dark:border-surface-dark/30 text-text-light dark:text-text-dark'
      ]"
    />
    
    <!-- Clear Button -->
    <button
      v-if="debouncedValue"
      @click="clearSearch"
      class="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-text-light/50 dark:text-text-dark/50 hover:text-text-light dark:hover:text-text-dark hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 rounded-full transition-colors"
    >
      <CrossIcon class="w-4 h-4" />
    </button>
  </div>
</template>
