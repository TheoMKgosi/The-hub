<script setup lang="ts">
import ChevronDownIcon from '../ui/svg/ChevronDownIcon.vue'
import ArrowUpIcon from '../ui/svg/ArrowUpIcon.vue'
import ArrowDownIcon from '../ui/svg/ArrowDownIcon.vue'

interface SortOption {
  value: string
  label: string
}

interface Props {
  modelValue: string
  order: 'asc' | 'desc'
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'update:order', value: 'asc' | 'desc'): void
}>()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement>()

onClickOutside(dropdownRef, () => {
  isOpen.value = false
})

const sortOptions: SortOption[] = [
  { value: 'due_date', label: 'Due Date' },
  { value: 'priority', label: 'Priority' },
  { value: 'created', label: 'Created Date' },
  { value: 'title', label: 'Title' }
]

const selectedLabel = computed(() => {
  const option = sortOptions.find(o => o.value === props.modelValue)
  return option?.label || 'Sort by'
})

const toggleOrder = () => {
  emit('update:order', props.order === 'asc' ? 'desc' : 'asc')
}

const selectOption = (value: string) => {
  emit('update:modelValue', value)
  isOpen.value = false
}
</script>

<template>
  <div ref="dropdownRef" class="relative flex items-center gap-1">
    <!-- Sort Dropdown -->
    <button
      @click="isOpen = !isOpen"
      class="inline-flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-md border transition-all duration-200"
      :class="[
        modelValue
          ? 'bg-primary/10 border-primary text-primary dark:text-primary'
          : 'bg-surface-light dark:bg-surface-dark border-surface-light/30 dark:border-surface-dark/30 text-text-light dark:text-text-dark hover:border-primary/50'
      ]"
    >
      <span>{{ selectedLabel }}</span>
      <ChevronDownIcon 
        class="w-4 h-4 transition-transform duration-200" 
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <!-- Order Toggle -->
    <button
      @click="toggleOrder"
      class="p-2 rounded-md border transition-all duration-200"
      :class="[
        'bg-surface-light dark:bg-surface-dark border-surface-light/30 dark:border-surface-dark/30 text-text-light dark:text-text-dark hover:border-primary/50'
      ]"
      :title="order === 'asc' ? 'Ascending' : 'Descending'"
    >
      <ArrowUpIcon v-if="order === 'asc'" class="w-4 h-4" />
      <ArrowDownIcon v-else class="w-4 h-4" />
    </button>

    <!-- Dropdown Menu -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 top-full mt-1 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-lg border border-surface-light/20 dark:border-surface-dark/20 z-50"
      >
        <div class="py-1">
          <button
            v-for="option in sortOptions"
            :key="option.value"
            @click="selectOption(option.value)"
            class="w-full flex items-center justify-between px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 transition-colors"
            :class="{ 'bg-primary/10': modelValue === option.value }"
          >
            <span>{{ option.label }}</span>
            <span v-if="modelValue === option.value" class="text-primary">✓</span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
