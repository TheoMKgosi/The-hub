<template>
  <div class="subscription-manager">
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
        Subscription Management
      </h2>
      <p class="text-gray-600 dark:text-gray-400">
        Manage your subscription and billing information
      </p>
    </div>

    <!-- Current Subscription Status -->
    <div v-if="subscriptionStore.currentSubscription" class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 mb-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          Current Subscription
        </h3>
        <span
          :class="[
            'px-3 py-1 rounded-full text-sm font-medium',
            subscriptionStore.currentSubscription.status === 'active'
              ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
              : subscriptionStore.currentSubscription.status === 'pending'
              ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
              : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
          ]"
        >
          {{ subscriptionStore.currentSubscription.status }}
        </span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">Plan</p>
          <p class="font-medium text-gray-900 dark:text-white">
            {{ subscriptionStore.currentSubscription.plan.name }}
          </p>
        </div>
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">Price</p>
          <p class="font-medium text-gray-900 dark:text-white">
            ${{ subscriptionStore.currentSubscription.plan.price }}/{{ subscriptionStore.currentSubscription.plan.interval }}
          </p>
        </div>
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">Start Date</p>
          <p class="font-medium text-gray-900 dark:text-white">
            {{ formatDate(subscriptionStore.currentSubscription.start_date) }}
          </p>
        </div>
        <div v-if="subscriptionStore.currentSubscription.end_date">
          <p class="text-sm text-gray-500 dark:text-gray-400">End Date</p>
          <p class="font-medium text-gray-900 dark:text-white">
            {{ formatDate(subscriptionStore.currentSubscription.end_date) }}
          </p>
        </div>
      </div>

      <div v-if="subscriptionStore.currentSubscription.status === 'active'" class="flex gap-3">
        <button
          @click="showCancelDialog = true"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
        >
          Cancel Subscription
        </button>
      </div>
    </div>

    <!-- Subscription Plans -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
        Available Plans
      </h3>

      <div v-if="subscriptionStore.loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="plan in subscriptionStore.plans"
          :key="plan.plan_id"
          :class="[
            'border rounded-lg p-6 transition-all',
            subscriptionStore.currentSubscription?.plan_id === plan.plan_id
              ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
              : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
          ]"
        >
          <div class="flex items-center justify-between mb-4">
            <h4 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ plan.name }}
            </h4>
            <span v-if="subscriptionStore.currentSubscription?.plan_id === plan.plan_id" class="text-sm text-blue-600 dark:text-blue-400 font-medium">
              Current Plan
            </span>
          </div>

          <div class="mb-4">
            <div class="text-3xl font-bold text-gray-900 dark:text-white mb-1">
              ${{ plan.price }}
              <span class="text-lg font-normal text-gray-500 dark:text-gray-400">
                /{{ plan.interval }}
              </span>
            </div>
            <p class="text-gray-600 dark:text-gray-400 text-sm">
              {{ plan.description }}
            </p>
          </div>

          <ul v-if="plan.features && plan.features.length > 0" class="mb-6 space-y-2">
            <li
              v-for="feature in plan.features"
              :key="feature"
              class="flex items-center text-sm text-gray-600 dark:text-gray-400"
            >
              <svg class="w-4 h-4 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              {{ feature }}
            </li>
          </ul>

          <button
            v-if="!subscriptionStore.currentSubscription || subscriptionStore.currentSubscription.plan_id !== plan.plan_id"
            @click="subscribeToPlan(plan.plan_id)"
            :disabled="subscriptionStore.loading"
            class="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {{ subscriptionStore.loading ? 'Processing...' : 'Subscribe' }}
          </button>

          <button
            v-else
            disabled
            class="w-full px-4 py-2 bg-gray-100 text-gray-500 rounded-lg cursor-not-allowed"
          >
            Current Plan
          </button>
        </div>
      </div>
    </div>

    <!-- Cancel Subscription Dialog -->
    <div
      v-if="showCancelDialog"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="showCancelDialog = false"
    >
      <div
        class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4"
        @click.stop
      >
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Cancel Subscription
        </h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          Are you sure you want to cancel your subscription? You'll still have access until the end of your current billing period.
        </p>
        <div class="flex gap-3">
          <button
            @click="showCancelDialog = false"
            class="flex-1 px-4 py-2 bg-gray-200 text-gray-800 dark:bg-gray-700 dark:text-gray-200 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            Keep Subscription
          </button>
          <button
            @click="confirmCancel"
            :disabled="subscriptionStore.loading"
            class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {{ subscriptionStore.loading ? 'Cancelling...' : 'Cancel Subscription' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSubscriptionStore } from '@/stores/subscription'

const subscriptionStore = useSubscriptionStore()
const showCancelDialog = ref(false)

// Initialize store data
onMounted(async () => {
  await subscriptionStore.init()
})

const subscribeToPlan = async (planId: string) => {
  try {
    await subscriptionStore.createSubscription(planId)
  } catch (error) {
    console.error('Failed to subscribe:', error)
  }
}

const confirmCancel = async () => {
  try {
    await subscriptionStore.cancelSubscription()
    showCancelDialog.value = false
  } catch (error) {
    console.error('Failed to cancel subscription:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.subscription-manager {
  @apply max-w-6xl mx-auto p-6;
}
</style>