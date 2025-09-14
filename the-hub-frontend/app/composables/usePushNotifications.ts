import { ref, computed } from 'vue'

interface PushSubscriptionData {
  endpoint: string
  keys: {
    p256dh: string
    auth: string
  }
}

export const usePushNotifications = () => {
  const permission = ref<NotificationPermission>('default')
  const subscription = ref<PushSubscription | null>(null)
  const isSupported = ref(false)
  const isSubscribed = ref(false)

  // Check if push notifications are supported
  const checkSupport = () => {
    isSupported.value = 'serviceWorker' in navigator &&
                       'PushManager' in window &&
                       'Notification' in window
  }

  // Request permission for notifications
  const requestPermission = async (): Promise<NotificationPermission> => {
    if (!isSupported.value) {
      throw new Error('Push notifications are not supported in this browser')
    }

    try {
      const result = await Notification.requestPermission()
      permission.value = result
      return result
    } catch (error) {
      console.error('Error requesting notification permission:', error)
      throw error
    }
  }

  // Subscribe to push notifications
  const subscribe = async (vapidPublicKey?: string): Promise<PushSubscription | null> => {
    if (!isSupported.value) {
      throw new Error('Push notifications are not supported')
    }

    if (permission.value !== 'granted') {
      throw new Error('Notification permission not granted')
    }

    try {
      const registration = await navigator.serviceWorker.ready
      const existingSubscription = await registration.pushManager.getSubscription()

      if (existingSubscription) {
        subscription.value = existingSubscription
        isSubscribed.value = true
        return existingSubscription
      }

      const applicationServerKey = vapidPublicKey ?
        urlBase64ToUint8Array(vapidPublicKey) :
        undefined

      const newSubscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey
      })

      subscription.value = newSubscription
      isSubscribed.value = true

      return newSubscription
    } catch (error) {
      console.error('Error subscribing to push notifications:', error)
      throw error
    }
  }

  // Unsubscribe from push notifications
  const unsubscribe = async (): Promise<void> => {
    if (!subscription.value) return

    try {
      await subscription.value.unsubscribe()
      subscription.value = null
      isSubscribed.value = false
    } catch (error) {
      console.error('Error unsubscribing from push notifications:', error)
      throw error
    }
  }

  // Send subscription to backend
  const sendSubscriptionToBackend = async (subscriptionData: PushSubscriptionData): Promise<void> => {
    try {
      const { $api } = useNuxtApp()
      await $api('/push/subscription', {
        method: 'POST',
        body: JSON.stringify({
          subscription: subscriptionData,
          user_agent: navigator.userAgent
        })
      })
    } catch (error) {
      console.error('Error sending subscription to backend:', error)
      throw error
    }
  }

  // Get subscription data in the format expected by backend
  const getSubscriptionData = (sub: PushSubscription): PushSubscriptionData => {
    const keys = sub.getKey ? {
      p256dh: arrayBufferToBase64(sub.getKey('p256dh')!),
      auth: arrayBufferToBase64(sub.getKey('auth')!)
    } : { p256dh: '', auth: '' }

    return {
      endpoint: sub.endpoint,
      keys
    }
  }

  // Initialize push notifications
  const initialize = async (vapidPublicKey?: string): Promise<void> => {
    checkSupport()

    if (!isSupported.value) {
      console.warn('Push notifications are not supported in this browser')
      return
    }

    // Check current permission
    permission.value = Notification.permission

    // Get existing subscription
    try {
      const registration = await navigator.serviceWorker.ready
      const existingSubscription = await registration.pushManager.getSubscription()
      if (existingSubscription) {
        subscription.value = existingSubscription
        isSubscribed.value = true
      }
    } catch (error) {
      console.error('Error getting existing subscription:', error)
    }

    // Auto-subscribe if permission is granted and not already subscribed
    if (permission.value === 'granted' && !isSubscribed.value && vapidPublicKey) {
      try {
        await subscribe(vapidPublicKey)
      } catch (error) {
        console.error('Auto-subscription failed:', error)
      }
    }
  }

  // Test notification
  const testNotification = async (): Promise<void> => {
    if (!isSupported.value || permission.value !== 'granted') {
      throw new Error('Notifications not available')
    }

    const notification = new Notification('Test Notification', {
      body: 'This is a test push notification from The Hub!',
      icon: '/icon-192x192.png',
      badge: '/icon-192x192.png',
      tag: 'test-notification'
    })

    // Auto-close after 4 seconds
    setTimeout(() => {
      notification.close()
    }, 4000)
  }

  // Utility function to convert VAPID key
  const urlBase64ToUint8Array = (base64String: string): Uint8Array => {
    const padding = '='.repeat((4 - base64String.length % 4) % 4)
    const base64 = (base64String + padding)
      .replace(/-/g, '+')
      .replace(/_/g, '/')

    const rawData = window.atob(base64)
    const outputArray = new Uint8Array(rawData.length)

    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i)
    }
    return outputArray
  }

  // Utility function to convert ArrayBuffer to base64
  const arrayBufferToBase64 = (buffer: ArrayBuffer): string => {
    const bytes = new Uint8Array(buffer)
    let binary = ''
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i])
    }
    return window.btoa(binary)
  }

  return {
    permission: computed(() => permission.value),
    subscription: computed(() => subscription.value),
    isSupported: computed(() => isSupported.value),
    isSubscribed: computed(() => isSubscribed.value),
    requestPermission,
    subscribe,
    unsubscribe,
    sendSubscriptionToBackend,
    getSubscriptionData,
    initialize,
    testNotification
  }
}