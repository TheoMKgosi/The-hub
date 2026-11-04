import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createPinia, setActivePinia, defineStore } from 'pinia'
import { ref, computed } from 'vue'

// The subscription store relies on Vue auto-imports (ref/computed) and Nuxt
// auto-imports (useToast); unlike the other stores it talks to the API through
// $fetch directly instead of $api via useNuxtApp. Stub them all before the
// store module loads; $fetch is stubbed with a plain mock.
const fetchMock = vi.fn()
const toastAddMock = vi.fn()

vi.stubGlobal('$fetch', fetchMock)
vi.stubGlobal('useToast', () => ({ add: toastAddMock }))
vi.stubGlobal('ref', ref)
vi.stubGlobal('computed', computed)
vi.stubGlobal('defineStore', defineStore)

const { useSubscriptionStore } = await import('../../../app/stores/subscription')

beforeEach(() => {
  setActivePinia(createPinia())
  fetchMock.mockReset()
  toastAddMock.mockClear()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

const plan = (planId = 'pro', price = 10) => ({
  plan_id: planId,
  name: 'Pro',
  description: 'Pro plan',
  price,
  currency: 'USD',
  interval: 'month',
  paypal_plan_id: `paypal-${planId}`,
  features: ['feature-1'],
  is_active: true
})

const subscription = (status: 'active' | 'cancelled' | 'expired' | 'pending' = 'active', overrides: Record<string, unknown> = {}) => ({
  subscription_id: 'sub-1',
  user_id: 'u1',
  plan_id: 'pro',
  plan: plan(),
  paypal_subscription_id: 'paypal-sub-1',
  status,
  start_date: '2026-01-01',
  auto_renew: true,
  ...overrides
})

describe('useSubscriptionStore', () => {
  it('starts with empty state and inactive subscription', () => {
    const store = useSubscriptionStore()

    expect(store.plans).toEqual([])
    expect(store.currentSubscription).toBeNull()
    expect(store.payments).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.isSubscribed).toBe(false)
    expect(store.isSubscriptionPending).toBe(false)
    expect(store.currentPlan).toBeUndefined()
  })

  describe('fetchSubscriptionPlans', () => {
    it('populates plans and clears loading', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockResolvedValue({ plans: [plan('free', 0), plan('pro', 10)] })

      await store.fetchSubscriptionPlans()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions/plans')
      expect(store.plans).toHaveLength(2)
      expect(store.loading).toBe(false)
    })

    it('shows an error toast on failure and clears loading', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockRejectedValue(new Error('server error'))

      await store.fetchSubscriptionPlans()

      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'error',
        title: 'Error',
        description: 'Failed to load subscription plans'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchUserSubscription', () => {
    it('stores the user subscription', async () => {
      const store = useSubscriptionStore()
      const userSubscription = subscription()
      fetchMock.mockResolvedValue({ subscription: userSubscription })

      await store.fetchUserSubscription()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions/user')
      expect(store.currentSubscription).toEqual(userSubscription)
      expect(store.loading).toBe(false)
    })

    it('clears the subscription when the user has none', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockResolvedValue({ subscription: null })

      await store.fetchUserSubscription()

      expect(store.currentSubscription).toBeNull()
    })

    it('shows an error toast on failure', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockRejectedValue(new Error('server error'))

      await store.fetchUserSubscription()

      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'error',
        title: 'Error',
        description: 'Failed to load subscription status'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('createSubscription', () => {
    it('stores the created subscription and shows a success toast', async () => {
      const store = useSubscriptionStore()
      const created = subscription('pending')
      fetchMock.mockResolvedValue({ subscription: created, paypal_approval_url: 'https://paypal.com/approve' })

      await store.createSubscription('pro')

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions', {
        method: 'POST',
        body: { plan_id: 'pro' }
      })
      expect(store.currentSubscription).toEqual(created)
      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'success',
        title: 'Subscription Created',
        description: 'Redirecting to PayPal for payment...'
      })
      expect(store.loading).toBe(false)
    })

    it('rethrows on failure and shows an error toast', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockRejectedValue(new Error('gateway down'))

      await expect(store.createSubscription('pro')).rejects.toThrow('gateway down')

      expect(store.currentSubscription).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'error',
        title: 'Error',
        description: 'Failed to create subscription'
      })
      expect(store.loading).toBe(false)
    })
  })

  describe('cancelSubscription', () => {
    it('marks the local subscription as cancelled and disables auto renew', async () => {
      const store = useSubscriptionStore()
      store.currentSubscription = subscription()
      fetchMock.mockResolvedValue(undefined)

      await store.cancelSubscription()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions', { method: 'DELETE' })
      expect(store.currentSubscription?.status).toBe('cancelled')
      expect(store.currentSubscription?.auto_renew).toBe(false)
      expect(store.currentSubscription?.cancelled_at).toEqual(expect.any(String))
      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'success',
        title: 'Subscription Cancelled',
        description: 'Your subscription has been cancelled'
      })
      expect(store.loading).toBe(false)
    })

    it('succeeds without state changes when there is no subscription', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockResolvedValue(undefined)

      await store.cancelSubscription()

      expect(store.currentSubscription).toBeNull()
      expect(toastAddMock).toHaveBeenCalledWith(expect.objectContaining({ color: 'success' }))
    })

    it('rethrows on failure, keeps state unchanged and shows an error toast', async () => {
      const store = useSubscriptionStore()
      store.currentSubscription = subscription()
      fetchMock.mockRejectedValue(new Error('server error'))

      await expect(store.cancelSubscription()).rejects.toThrow('server error')

      expect(store.currentSubscription?.status).toBe('active')
      expect(store.currentSubscription?.auto_renew).toBe(true)
      expect(store.currentSubscription?.cancelled_at).toBeUndefined()
      expect(toastAddMock).toHaveBeenCalledWith({
        color: 'error',
        title: 'Error',
        description: 'Failed to cancel subscription'
      })
      expect(store.loading).toBe(false)
    })
  })

  // NOTE: fetchPayments is currently a no-op stub (the request is commented
  // out in the store); these tests document that actual behavior.
  describe('fetchPayments', () => {
    it('is a no-op: no request is made and payments stay empty', async () => {
      const store = useSubscriptionStore()

      await store.fetchPayments()

      expect(fetchMock).not.toHaveBeenCalled()
      expect(store.payments).toEqual([])
    })
  })

  describe('getters', () => {
    it('isSubscribed is only true for an active subscription', () => {
      const store = useSubscriptionStore()
      store.currentSubscription = subscription('active')
      expect(store.isSubscribed).toBe(true)

      store.currentSubscription = subscription('pending')
      expect(store.isSubscribed).toBe(false)

      store.currentSubscription = subscription('cancelled')
      expect(store.isSubscribed).toBe(false)
    })

    it('isSubscriptionPending is only true for a pending subscription', () => {
      const store = useSubscriptionStore()
      store.currentSubscription = subscription('pending')
      expect(store.isSubscriptionPending).toBe(true)

      store.currentSubscription = subscription('active')
      expect(store.isSubscriptionPending).toBe(false)
    })

    it('currentPlan exposes the plan of the current subscription', () => {
      const store = useSubscriptionStore()
      store.currentSubscription = subscription('active', { plan: plan('basic', 5) })

      expect(store.currentPlan?.plan_id).toBe('basic')
      expect(store.currentPlan?.price).toBe(5)
    })
  })

  describe('init', () => {
    it('fetches plans and the user subscription together', async () => {
      const store = useSubscriptionStore()
      fetchMock.mockImplementation((url: string) => {
        if (url === '/api/v1/subscriptions/plans') return Promise.resolve({ plans: [plan()] })
        return Promise.resolve({ subscription: subscription() })
      })

      await store.init()

      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions/plans')
      expect(fetchMock).toHaveBeenCalledWith('/api/v1/subscriptions/user')
      expect(store.plans).toHaveLength(1)
      expect(store.currentSubscription).toEqual(subscription())
      expect(store.loading).toBe(false)
    })
  })
})
