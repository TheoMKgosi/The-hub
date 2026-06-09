<script setup lang="ts">
import { useIntersectionObserver } from '@vueuse/core';
import { ref } from 'vue'

const pricing = [
  {
    title: 'Basic',
    description: 'Get started with the basic plan for free',
    price: '$0',
    features: [
      'Unlimited tasks',
      'Unlimited finance',
      'Unlimited learning',
      'Basic analytics'
    ],
    color: 'primary'
  },
  {
    title: 'Pro',
    description: 'Get flexible with the Pro plan',
    price: '$5',
    features: [
      'Everything in the basic plan',
      'Theme days',
      'Budgeting',
      'Advanced analytics'
    ],
    color: 'secondary'
  },
  {
    title: 'AI',
    description: 'Enhance with the AI plan',
    price: '$25',
    features: [
      'Everything in the Pro plan',
      'Automatic schedule',
      'Generated quizzes',
    ],
    color: 'accent'
  }
]

const pricingRef = ref<HTMLElement | null>(null)

const pricingVisible = ref(true)

useIntersectionObserver(pricingRef, ([{ isIntersecting }]) => {
  if (isIntersecting) pricingVisible.value = true
}, { threshold: 0.1 })

</script>

<template>
  <div ref="pricingRef" class="relative">
    <!-- Pricing Stack -->
    <div class="flex flex-col md:flex-row items-center justify-center gap-3">
      <div v-for="(price, _) in pricing" :key="price.title"
        class="flex-1 w-full max-w-72 inset-0 transition-all duration-500 ease-[cubic-bezier(0.16,1,0.3,1)]"
        :class="pricingVisible ? 'opacity-100' : 'opacity-0 translate-y-8'"
        <!-- Card -->
        <div
          class="group relative h-full w-full bg-background-light dark:bg-background-dark rounded-2xl p-8 border border-surface-light dark:border-surface-dark overflow-hidden shadow-lg">
          <!-- Background Gradient -->
          <div class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500" :class="{
            'bg-linear-to-br from-primary/5 to-transparent': price.color === 'accent',
            'bg-linear-to-br from-secondary/5 to-transparent': price.color === 'secondary',
            'bg-linear-to-br from-accent/5 to-transparent': price.color === 'primary'
          }" />

          <!-- Content -->
          <div class="relative">
            <!-- Title & Price -->
            <div class="flex items-baseline justify-between mb-3">
              <h3 class="text-xl font-bold text-text-light dark:text-text-dark">
                {{ price.title }}
              </h3>
              <span class="text-2xl font-bold" :class="{
                'text-primary/80': price.color === 'accent',
                'text-secondary/80': price.color === 'secondary',
                'text-accent/80': price.color === 'primary'
              }">
                {{ price.price }}
              </span>
            </div>

            <!-- Description -->
            <p class="text-text-light/70 dark:text-text-dark/70 leading-relaxed mb-4">
              {{ price.description }}
            </p>
            <ul class="list-decimal list-inside text-text-light/70 dark:text-text-dark/70 leading-relaxed space-y-1">
              <li v-for="feature in price.features" :key="feature">
                {{ feature }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
