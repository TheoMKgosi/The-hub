<script setup lang="ts">
interface StarRatingProps {
  modelValue: number | null
  maxRating?: number
  readonly?: boolean
  size?: 'sm' | 'md' | 'lg'
}

interface StarRatingEmits {
  (e: 'update:modelValue', value: number | null): void
}

const props = withDefaults(defineProps<StarRatingProps>(), {
  maxRating: 5,
  readonly: false,
  size: 'md'
})

const emit = defineEmits<StarRatingEmits>()

const hoverRating = ref<number | null>(null)

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6'
}

const setRating = (rating: number) => {
  if (props.readonly) return
  emit('update:modelValue', rating)
}

const clearRating = () => {
  if (props.readonly) return
  emit('update:modelValue', null)
}

const getStarClass = (starIndex: number) => {
  const currentRating = hoverRating.value ?? props.modelValue ?? 0
  const isFilled = starIndex <= currentRating

  return [
    sizeClasses[props.size],
    'transition-colors duration-150',
    props.readonly ? 'cursor-default' : 'cursor-pointer',
    isFilled ? 'text-yellow-400 fill-current' : 'text-gray-300 dark:text-gray-600 hover:text-yellow-300'
  ].join(' ')
}
</script>

<template>
  <div class="flex items-center gap-1">
    <button
      v-for="star in maxRating"
      :key="star"
      type="button"
      :class="getStarClass(star)"
      @click="setRating(star)"
      @mouseenter="hoverRating = star"
      @mouseleave="hoverRating = null"
      :disabled="readonly"
      :aria-label="`Rate ${star} star${star > 1 ? 's' : ''}`"
    >
      <svg viewBox="0 0 24 24" class="transition-transform duration-150 hover:scale-110">
        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
      </svg>
    </button>

    <button
      v-if="!readonly && modelValue"
      type="button"
      @click="clearRating"
      class="ml-2 text-xs text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors"
      aria-label="Clear rating"
    >
      Clear
    </button>
  </div>
</template>