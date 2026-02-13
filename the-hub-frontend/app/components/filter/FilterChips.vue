<script setup lang="ts">
import CrossIcon from '../ui/svg/CrossIcon.vue'

interface FilterChip {
  key: string
  label: string
  value: string
}

interface Props {
  chips: FilterChip[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'remove', key: string): void
  (e: 'clearAll'): void
}>()
</script>

<template>
  <div class="flex flex-wrap items-center gap-2">
    <!-- Filter Chips -->
    <TransitionGroup
      enter-active-class="transition ease-out duration-200"
      enter-from-class="transform opacity-0 scale-90"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-90"
    >
      <span
        v-for="chip in chips"
        :key="chip.key"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm rounded-full bg-primary/10 text-primary border border-primary/20"
      >
        <span class="font-medium">{{ chip.label }}:</span>
        <span>{{ chip.value }}</span>
        <button
          @click="emit('remove', chip.key)"
          class="p-0.5 hover:bg-primary/20 rounded-full transition-colors"
          title="Remove filter"
        >
          <CrossIcon class="w-3.5 h-3.5" />
        </button>
      </span>
    </TransitionGroup>
    
    <!-- Clear All Button -->
    <button
      v-if="chips.length > 1"
      @click="emit('clearAll')"
      class="text-sm text-red-500 hover:text-red-600 hover:underline transition-colors"
    >
      Clear all
    </button>
  </div>
</template>
