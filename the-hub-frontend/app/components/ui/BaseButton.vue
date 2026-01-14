<script setup lang="ts">
import { computed } from 'vue'

interface ButtonProps {
  text?: string
  variant?: 'primary' | 'secondary' | 'danger' | 'default'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
  icon?: string // Icon name or path
  iconPosition?: 'left' | 'right' // Icon position relative to text
  iconOnly?: boolean // Show only icon, no text
}

const props = withDefaults(defineProps<ButtonProps>(), {
  text: 'Button',
  variant: 'default',
  size: 'md',
  disabled: false,
  type: 'button',
  icon: undefined,
  iconPosition: 'left',
  iconOnly: false
})

const baseClasses = 'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'

const variantClasses = {
  primary: 'bg-primary text-white hover:bg-primary/90 focus:ring-primary',
  secondary: 'bg-secondary text-white hover:bg-secondary/90 focus:ring-secondary',
  danger: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500',
  default: 'bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600'
}

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm gap-1',
  md: 'px-4 py-2 text-base gap-2',
  lg: 'px-6 py-3 text-lg gap-2'
}

const iconSizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6'
}

const classes = computed(() => [
  baseClasses,
  variantClasses[props.variant],
  sizeClasses[props.size]
])

const iconClasses = computed(() => [
  iconSizeClasses[props.size],
  props.iconOnly ? '' : (props.iconPosition === 'left' ? 'mr-1' : 'ml-1')
])
</script>

<template>
  <button :class="classes" :disabled="disabled" :type="type">
    <template v-if="icon">
      <span v-if="iconPosition === 'left' || iconOnly" :class="iconClasses">
        <!-- Replace with your icon component or SVG -->
        <component :is="icon" v-if="typeof icon === 'object'" />
        <i v-else :class="icon" class="inline-block"></i>
      </span>
    </template>

    <span v-if="!iconOnly">{{ text }}</span>

    <template v-if="icon && iconPosition === 'right' && !iconOnly">
      <span :class="iconClasses">
        <!-- Replace with your icon component or SVG -->
        <component :is="icon" v-if="typeof icon === 'object'" />
        <i v-else :class="icon" class="inline-block"></i>
      </span>
    </template>
  </button>
</template>

