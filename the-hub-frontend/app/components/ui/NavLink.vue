<script setup lang="ts">
interface NavLinkProps {
  to?: string
  href?: string
  active?: boolean
  variant?: 'nav' | 'tab'
}

const props = withDefaults(defineProps<NavLinkProps>(), {
  variant: 'nav'
})

const isActive = computed(() => props.active)

const baseClasses = 'transition-all duration-300 font-medium'

const variantClasses = {
  nav: [
    'flex items-center px-4 py-3 rounded-xl w-full',
    'text-text-light dark:text-text-dark',
    'hover:shadow-md hover:text-gray-900 hover:bg-gray-400/30 dark:hover:text-gray-100 dark:hover:bg-gray-700/30',
    isActive.value ? 'text-gray-400 font-semibold dark:text-gray-300' : ''
  ],
  tab: [
    'px-4 py-2',
    isActive.value
      ? 'border-b-2 border-primary text-primary dark:text-primary'
      : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
  ]
}

const classes = computed(() => [
  baseClasses,
  ...variantClasses[props.variant]
])
</script>

<template>
  <NuxtLink v-if="to" :to="to" :class="classes">
    <slot />
  </NuxtLink>
  <a v-else-if="href" :href="href" :class="classes">
    <slot />
  </a>
  <button v-else :class="classes" @click="$emit('click')">
    <slot />
  </button>
</template>