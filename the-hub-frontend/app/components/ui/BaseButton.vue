<script setup lang="ts">

interface Props {
  text?: string
  variant?: 'primary' | 'secondary' | 'danger' | 'default' | 'clear'
  size?: 'sm' | 'md' | 'lg' | 'full'
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
  icon?: string | object // Icon name or path
  iconPosition?: 'left' | 'right'
  iconOnly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  text: 'Button',
  variant: 'default',
  size: 'md',
  disabled: false,
  type: 'button',
  icon: undefined,
  iconPosition: 'left',
  iconOnly: false
})

const baseClasses = 'inline-flex items-center  font-medium rounded-lg transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed'

const variantClasses = {
  primary: 'bg-primary text-white hover:bg-primary/90 focus:ring-primary justify-center',
  secondary: 'bg-secondary text-white hover:bg-secondary/90 focus:ring-secondary justify-center',
  danger: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500 justify-center',
  default: 'bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600 justify-center',
  clear: 'text-gray-900 hover:bg-gray-300 focus:ring-gray-500 dark:text-gray-100 dark:hover:bg-gray-600 p-2'
}

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm gap-1 ',
  md: 'px-4 py-2 text-base gap-2',
  lg: 'px-6 py-3 text-lg gap-2',
  full: 'w-full text-lg'
}

const iconSizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6',
  full: 'w-6 h-6'
}

const classes = computed(() => [
  baseClasses,
  variantClasses[props.variant],
  sizeClasses[props.size]
])

const iconClasses = computed(() => [
  'flex justify-center items-center shrink-0',
  iconSizeClasses[props.size],
  props.iconOnly ? '' : (props.iconPosition === 'left' ? 'mr-1' : 'ml-1')
])
</script>

<template>
  <button :class="classes" :disabled="disabled" :type="type">
    <template v-if="icon">
      <span v-if="iconPosition === 'left' || iconOnly" :class="iconClasses">
        <component :is="icon" v-if="typeof icon === 'object'" />
        <i v-else :class="icon"></i>
      </span>
    </template>

    <span v-if="!iconOnly">{{ text }}</span>

    <template v-if="icon && iconPosition === 'right' && !iconOnly">
      <span :class="iconClasses">
        <!-- Replace with your icon component or SVG -->
        <component :is="icon" v-if="typeof icon === 'object'" />
        <i v-else :class="icon"></i>
      </span>
    </template>
  </button>
</template>
