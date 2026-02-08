<script setup lang="ts">
import ChevronDownIcon from '../ui/svg/ChevronDownIcon.vue'
import CheckIcon from '../ui/svg/CheckIcon.vue'

interface FilterOption {
  value: string
  label: string
  color?: string
}

interface Props {
  label: string
  options: FilterOption[]
  modelValue: string | string[]
  multiple?: boolean
  placeholder?: string
  clearable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  multiple: false,
  placeholder: 'Select...',
  clearable: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | string[]): void
}>()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement>()

// Close on click outside
onClickOutside(dropdownRef, () => {
  isOpen.value = false
})

// Get selected label(s)
const selectedLabel = computed(() => {
  if (props.multiple && Array.isArray(props.modelValue) && props.modelValue.length > 0) {
    if (props.modelValue.length === 1) {
      const option = props.options.find(o => o.value === props.modelValue[0])
      return option?.label || props.modelValue[0]
    }
    return `${props.modelValue.length} selected`
  }
  
  if (!props.multiple && props.modelValue) {
    const option = props.options.find(o => o.value === props.modelValue)
    return option?.label || props.modelValue
  }
  
  return props.placeholder
})

// Check if option is selected
const isSelected = (value: string): boolean => {
  if (props.multiple && Array.isArray(props.modelValue)) {
    return props.modelValue.includes(value)
  }
  return props.modelValue === value
}

// Toggle option selection
const toggleOption = (value: string) => {
  if (props.multiple) {
    const current = Array.isArray(props.modelValue) ? [...props.modelValue] : []
    const index = current.indexOf(value)
    
    if (index > -1) {
      current.splice(index, 1)
    } else {
      current.push(value)
    }
    
    emit('update:modelValue', current)
  } else {
    emit('update:modelValue', value)
    isOpen.value = false
  }
}

// Clear selection
const clearSelection = (e: Event) => {
  e.stopPropagation()
  if (props.multiple) {
    emit('update:modelValue', [])
  } else {
    emit('update:modelValue', '')
  }
}

// Check if has value
const hasValue = computed(() => {
  if (props.multiple && Array.isArray(props.modelValue)) {
    return props.modelValue.length > 0
  }
  return !!props.modelValue
})
</script>

<template>
  <div ref="dropdownRef" class="relative">
    <!-- Trigger Button -->
    <button
      @click="isOpen = !isOpen"
      class="inline-flex items-center justify-between gap-2 px-3 py-2 text-sm font-medium rounded-md border transition-all duration-200 min-w-[120px]"
      :class="[
        hasValue
          ? 'bg-primary/10 border-primary text-primary dark:text-primary'
          : 'bg-surface-light dark:bg-surface-dark border-surface-light/30 dark:border-surface-dark/30 text-text-light dark:text-text-dark hover:border-primary/50'
      ]"
    >
      <span class="truncate">{{ selectedLabel }}</span>
      <div class="flex items-center gap-1">
        <button
          v-if="clearable && hasValue"
          @click="clearSelection"
          class="p-0.5 hover:bg-primary/20 rounded-full transition-colors"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
        <ChevronDownIcon 
          class="w-4 h-4 transition-transform duration-200" 
          :class="{ 'rotate-180': isOpen }"
        />
      </div>
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
        class="absolute z-50 mt-1 min-w-[200px] max-h-[300px] overflow-auto bg-surface-light dark:bg-surface-dark rounded-md shadow-lg border border-surface-light/20 dark:border-surface-dark/20"
      >
        <div class="py-1">
          <!-- Options -->
          <button
            v-for="option in options"
            :key="option.value"
            @click="toggleOption(option.value)"
            class="w-full flex items-center justify-between px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 transition-colors"
          >
            <div class="flex items-center gap-2">
              <!-- Color indicator -->
              <span
                v-if="option.color"
                class="w-3 h-3 rounded-full"
                :style="{ backgroundColor: option.color }"
              ></span>
              <span>{{ option.label }}</span>
            </div>
            
            <!-- Checkmark -->
            <CheckIcon
              v-if="isSelected(option.value)"
              class="w-4 h-4 text-primary"
            />
          </button>
          
          <!-- Clear option -->
          <div v-if="clearable && hasValue" class="border-t border-surface-light/20 dark:border-surface-dark/20 mt-1 pt-1">
            <button
              @click="clearSelection"
              class="w-full px-4 py-2 text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors text-left"
            >
              Clear selection
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>
