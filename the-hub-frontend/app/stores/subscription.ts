import { defineStore } from 'pinia'
import { useToast } from '@/composables/useToast'

interface SubscriptionPlan {
  plan_id: string
  name: string
  description: string
  price: number
  currency: string
  interval: string
  paypal_plan_id: string
  features: string[]
  is_active: boolean
}

interface Subscription {
  subscription_id: string
  user_id: string
  plan_id: string
  plan: SubscriptionPlan
  paypal_subscription_id: string
  status: 'active' | 'cancelled' | 'expired' | 'pending'
  start_date: string
  end_date?: string
  cancelled_at?: string
  auto_renew: boolean
}

interface Payment {
  payment_id: string
  subscription_id: string
  user_id: string
  paypal_payment_id: string
  amount: number
  currency: string
  status: 'completed' | 'pending' | 'failed' | 'refunded'
  payment_date: string
  description: string
}

export const useSubscriptionStore = defineStore('subscription', () => {
  const { addToast } = useToast()

  // State
  const plans = ref<SubscriptionPlan[]>([])
  const currentSubscription = ref<Subscription | null>(null)
  const payments = ref<Payment[]>([])
  const loading = ref(false)

  // Getters
  const isSubscribed = computed(() => {
    return currentSubscription.value?.status === 'active'
  })

  const isSubscriptionPending = computed(() => {
    return currentSubscription.value?.status === 'pending'
  })

  const currentPlan = computed(() => {
    return currentSubscription.value?.plan
  })

  // Actions
  const fetchSubscriptionPlans = async () => {
    try {
      loading.value = true
      const response = await $fetch<{
        plans: SubscriptionPlan[]
      }>('/api/v1/subscriptions/plans')

      plans.value = response.plans
    } catch (error) {
      console.error('Failed to fetch subscription plans:', error)
      addToast({
        type: 'error',
        title: 'Error',
        description: 'Failed to load subscription plans'
      })
    } finally {
      loading.value = false
    }
  }

  const fetchUserSubscription = async () => {
    try {
      loading.value = true
      const response = await $fetch<{
        subscription: Subscription | null
      }>('/api/v1/subscriptions/user')

      currentSubscription.value = response.subscription
    } catch (error) {
      console.error('Failed to fetch user subscription:', error)
      addToast({
        type: 'error',
        title: 'Error',
        description: 'Failed to load subscription status'
      })
    } finally {
      loading.value = false
    }
  }

  const createSubscription = async (planId: string) => {
    try {
      loading.value = true
      const response = await $fetch<{
        subscription: Subscription
        paypal_approval_url: string
      }>('/api/v1/subscriptions', {
        method: 'POST',
        body: { plan_id: planId }
      })

      currentSubscription.value = response.subscription

      // Redirect to PayPal approval URL
      if (import.meta.client) {
        window.location.href = response.paypal_approval_url
      }

      addToast({
        type: 'success',
        title: 'Subscription Created',
        description: 'Redirecting to PayPal for payment...'
      })
    } catch (error) {
      console.error('Failed to create subscription:', error)
      addToast({
        type: 'error',
        title: 'Error',
        description: 'Failed to create subscription'
      })
      throw error
    } finally {
      loading.value = false
    }
  }

  const cancelSubscription = async () => {
    try {
      loading.value = true
      await $fetch('/api/v1/subscriptions', {
        method: 'DELETE'
      })

      // Update local state
      if (currentSubscription.value) {
        currentSubscription.value.status = 'cancelled'
        currentSubscription.value.cancelled_at = new Date().toISOString()
        currentSubscription.value.auto_renew = false
      }

      addToast({
        type: 'success',
        title: 'Subscription Cancelled',
        description: 'Your subscription has been cancelled'
      })
    } catch (error) {
      console.error('Failed to cancel subscription:', error)
      addToast({
        type: 'error',
        title: 'Error',
        description: 'Failed to cancel subscription'
      })
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchPayments = async () => {
    // This would fetch payment history - implement when backend supports it
    try {
      // payments.value = await $fetch<Payment[]>('/api/v1/subscriptions/payments')
    } catch (error) {
      console.error('Failed to fetch payments:', error)
    }
  }

  // Initialize store
  const init = async () => {
    await Promise.all([
      fetchSubscriptionPlans(),
      fetchUserSubscription()
    ])
  }

  return {
    // State
    plans,
    currentSubscription,
    payments,
    loading,

    // Getters
    isSubscribed,
    isSubscriptionPending,
    currentPlan,

    // Actions
    fetchSubscriptionPlans,
    fetchUserSubscription,
    createSubscription,
    cancelSubscription,
    fetchPayments,
    init
  }
})