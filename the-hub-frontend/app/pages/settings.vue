<script setup lang="ts">
const auth = useAuthStore()
const { themeMode, setTheme } = useDarkMode()
const { addToast } = useToast()
const {
  isSupported,
  permission,
  isSubscribed,
  requestPermission,
  subscribe,
  unsubscribe,
  testNotification,
  initialize
} = usePushNotifications()

// User settings state
const userSettings = ref<any>({})
const name = ref(auth.user?.name || '')
const email = ref(auth.user?.email || '')
const userId = ref(auth.user?.user_id || '')
const isLoading = ref(false)
const isSaving = ref(false)

// Tab state
const activeTab = ref('profile')
const tabs = [
  { id: 'profile', label: 'Profile', icon: 'user' },
  { id: 'appearance', label: 'Appearance', icon: 'palette' },
  { id: 'calendar', label: 'Calendar', icon: 'calendar' },
  { id: 'tasks', label: 'Tasks', icon: 'checklist' },
  { id: 'notifications', label: 'Notifications', icon: 'bell' }
]

// Tasks settings state
const trimodalEnabled = ref(false)

// Load user settings on mount
onMounted(async () => {
  await loadUserSettings()
  // Initialize push notifications
  if (process.client) {
    await initialize()
  }
})

// Load settings from backend
const loadUserSettings = async () => {
  try {
    isLoading.value = true
    const { $api } = useNuxtApp()

    if (!userId.value) return

    const response = await $api(`/users/${userId.value}/settings`)
    userSettings.value = response.settings || {}

    // Sync theme with backend settings
    if (userSettings.value.theme?.mode) {
      setTheme(userSettings.value.theme.mode)
    }

    // Sync tri-modal setting with backend settings
    if (userSettings.value.tasks?.trimodalEnabled !== undefined) {
      trimodalEnabled.value = userSettings.value.tasks.trimodalEnabled
    }
  } catch (error) {
    addToast('Failed to load settings', 'error')
  } finally {
    isLoading.value = false
  }
}

// Save user profile and settings
const saveChanges = async () => {
  try {
    isSaving.value = true
    const { $api } = useNuxtApp()

    if (!userId.value) return

    // Update profile
    await $api(`/users/${userId.value}`, {
      method: 'PUT',
      body: JSON.stringify({
        name: name.value,
        email: email.value,
        settings: {
          ...userSettings.value,
          theme: {
            ...userSettings.value.theme,
            mode: themeMode.value
          },
          tasks: {
            ...userSettings.value.tasks,
            trimodalEnabled: trimodalEnabled.value
          }
        }
      })
    })

    addToast('Settings saved successfully!', 'success')
  } catch (error) {
    addToast('Failed to save settings', 'error')
  } finally {
    isSaving.value = false
  }
}

// Update theme and sync with backend
const updateTheme = async (newTheme: string) => {
  try {
    setTheme(newTheme)

    // Update local state
    userSettings.value = {
      ...userSettings.value,
      theme: {
        ...userSettings.value.theme,
        mode: newTheme
      }
    }

    // Auto-save theme to backend
    const { $api } = useNuxtApp()

    if (userId.value) {
      await $api(`/users/${userId.value}/settings`, {
        method: 'PATCH',
        body: JSON.stringify({
          theme: {
            ...userSettings.value.theme,
            mode: newTheme
          }
        })
      })
    }
  } catch (error) {
    addToast('Failed to save theme preference', 'error')
  }
}

// Push notification handlers
const enableNotifications = async () => {
  try {
    const perm = await requestPermission()
    if (perm === 'granted') {
      // Get VAPID key from runtime config
      const { $config } = useNuxtApp()
      const vapidKey = $config.public?.vapidPublicKey

      const subscription = await subscribe(vapidKey)
      if (subscription) {
        // Send subscription to backend
        const { sendSubscriptionToBackend, getSubscriptionData } = usePushNotifications()
        await sendSubscriptionToBackend(getSubscriptionData(subscription))
        addToast('Push notifications enabled!', 'success')
      }
    } else {
      addToast('Notification permission denied', 'error')
    }
  } catch (error) {
    addToast('Failed to enable notifications', 'error')
  }
}

const disableNotifications = async () => {
  try {
    await unsubscribe()
    addToast('Push notifications disabled', 'info')
  } catch (error) {
    console.error('Failed to disable notifications:', error)
    addToast('Failed to disable notifications', 'error')
  }
}

const sendTestNotification = async () => {
  try {
    await testNotification()
    addToast('Test notification sent!', 'success')
  } catch (error) {
    console.error('Failed to send test notification:', error)
    addToast('Failed to send test notification', 'error')
  }
}

// Save tri-modal setting immediately when toggled
const saveTrimodalSetting = async () => {
  try {
    const { $api } = useNuxtApp()

    if (!userId.value) return

    // Update local state
    userSettings.value = {
      ...userSettings.value,
      task: {
        ...userSettings.value.tasks,
        tri_modal: trimodalEnabled.value
      }
    }

    // Save to backend
    await $api(`/users/${userId.value}/settings`, {
      method: 'PATCH',
      body: JSON.stringify({
        task: {
          tri_modal: trimodalEnabled.value
        }
      })
    })

    addToast(`Tri-modal tasks ${trimodalEnabled.value ? 'enabled' : 'disabled'}!`, 'success')
  } catch (error) {
    addToast('Failed to save tri-modal setting', 'error')
    // Revert state on error
    trimodalEnabled.value = !trimodalEnabled.value
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto p-6 space-y-8">
    <h1 class="text-3xl font-bold text-text-light dark:text-text-dark mb-8">Settings</h1>

    <!-- Tab Navigation -->
    <div class="border-b border-gray-200 dark:border-gray-700 mb-8">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id" :class="[
          'py-2 px-4 text-sm font-medium transition-all duration-200',
          activeTab === tab.id
            ? 'border-2 border-primary text-primary rounded-t-lg bg-primary/5 shadow-sm'
            : 'border-0 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 border-b-2 border-gray-300 dark:border-gray-600'
        ]">
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- Tab Content -->
    <div class="space-y-8">
      <!-- Profile Tab -->
      <div v-show="activeTab === 'profile'" class="space-y-8">
        <!-- Profile Information -->
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 shadow-lg border border-gray-200 dark:border-gray-700">
          <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6">Profile Information</h2>
          <div class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Name:</label>
              <input type="text" id="name" v-model="name"
                class="w-full px-4 py-2 bg-background-light dark:bg-background-dark border border-gray-300 dark:border-gray-600 rounded-lg text-text-light dark:text-text-dark focus:ring-2 focus:ring-primary focus:border-transparent transition-colors">
            </div>
            <div>
              <label for="email"
                class="block text-sm font-medium text-text-light dark:text-text-dark mb-2">Email:</label>
              <input type="email" id="email" v-model="email"
                class="w-full px-4 py-2 bg-background-light dark:bg-background-dark border border-gray-300 dark:border-gray-600 rounded-lg text-text-light dark:text-text-dark focus:ring-2 focus:ring-primary focus:border-transparent transition-colors">
            </div>
            <button @click="saveChanges" :disabled="isSaving"
              class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-orange-600 focus:ring-2 focus:ring-primary focus:ring-offset-2 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed">
              {{ isSaving ? 'Saving...' : 'Save Changes' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Appearance Tab -->
    <div v-show="activeTab === 'appearance'" class="space-y-8">
      <!-- Theme Settings -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 shadow-lg border border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6">Appearance</h2>
        <div class="space-y-4">
          <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-3">Theme Mode:</label>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <button @click="updateTheme('light')" :class="[
              'p-4 rounded-lg border-2 transition-all duration-200 text-left',
              themeMode === 'light'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
            ]">
              <div class="flex items-center space-x-3">
                <div
                  class="w-8 h-8 rounded-full bg-linear-to-br from-yellow-400 to-orange-500 flex items-center justify-center">
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z"
                      clip-rule="evenodd" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-text-light dark:text-text-dark">Light</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">Use light theme</div>
                </div>
              </div>
            </button>

            <button @click="updateTheme('dark')" :class="[
              'p-4 rounded-lg border-2 transition-all duration-200 text-left',
              themeMode === 'dark'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
            ]">
              <div class="flex items-center space-x-3">
                <div
                  class="w-8 h-8 rounded-full bg-linear-to-br from-gray-800 to-gray-900 flex items-center justify-center">
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-text-light dark:text-text-dark">Dark</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">Use dark theme</div>
                </div>
              </div>
            </button>

            <button @click="updateTheme('system')" :class="[
              'p-4 rounded-lg border-2 transition-all duration-200 text-left',
              themeMode === 'system'
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
            ]">
              <div class="flex items-center space-x-3">
                <div
                  class="w-8 h-8 rounded-full bg-linear-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zm0 4a1 1 0 011-1h12a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1V8z"
                      clip-rule="evenodd" />
                  </svg>
                </div>
                <div>
                  <div class="font-medium text-text-light dark:text-text-dark">System</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">Follow system preference</div>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Calendar Tab -->
    <div v-show="activeTab === 'calendar'" class="space-y-8">
      <!-- Calendar Integration Settings -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 shadow-lg border border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6">Calendar Integration</h2>
        <CalendarIntegrations />
      </div>
    </div>

    <!-- Tasks Tab -->
    <div v-show="activeTab === 'tasks'" class="space-y-8">
      <!-- Task Settings -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 shadow-lg border border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6">Task Settings</h2>
        <div class="space-y-6">
          <!-- Tri-modal Setting -->
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div class="flex-1">
              <h3 class="text-lg font-medium text-text-light dark:text-text-dark mb-2">Enable Tri-modal Tasks</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                Enable tri-modal task organization to categorize tasks as high, medium, or low priority for better focus
                and productivity.
              </p>
            </div>
            <div class="ml-4">
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="trimodalEnabled" @change="saveTrimodalSetting" class="sr-only peer">
                <div
                  class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 dark:peer-focus:ring-primary/30 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary">
                </div>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Notifications Tab -->
    <div v-show="activeTab === 'notifications'" class="space-y-8">
      <!-- Push Notifications Settings -->
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-xl p-6 shadow-lg border border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-6">Push Notifications</h2>

        <div v-if="!isSupported" class="text-center py-8">
          <div class="text-gray-500 dark:text-gray-400 mb-4">
            <svg class="w-12 h-12 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <p class="text-lg font-medium">Push notifications not supported</p>
            <p class="text-sm">Your browser doesn't support push notifications.</p>
          </div>
        </div>

        <div v-else class="space-y-6">
          <!-- Permission Status -->
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div>
              <h3 class="font-medium text-text-light dark:text-text-dark">Notification Permission</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                {{ permission === 'granted' ? 'Allowed' : permission === 'denied' ? 'Blocked' : 'Not requested' }}
              </p>
            </div>
            <div class="flex items-center">
              <div :class="[
                'w-3 h-3 rounded-full',
                permission === 'granted' ? 'bg-green-500' :
                  permission === 'denied' ? 'bg-red-500' : 'bg-yellow-500'
              ]"></div>
            </div>
          </div>

          <!-- Subscription Status -->
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div>
              <h3 class="font-medium text-text-light dark:text-text-dark">Push Subscription</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                {{ isSubscribed ? 'Active' : 'Inactive' }}
              </p>
            </div>
            <div class="flex items-center">
              <div :class="[
                'w-3 h-3 rounded-full',
                isSubscribed ? 'bg-green-500' : 'bg-gray-400'
              ]"></div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex flex-wrap gap-3">
            <button v-if="!isSubscribed && permission !== 'denied'" @click="enableNotifications"
              class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-orange-600 focus:ring-2 focus:ring-primary focus:ring-offset-2 transition-colors font-medium">
              Enable Notifications
            </button>

            <button v-if="isSubscribed" @click="disableNotifications"
              class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors font-medium">
              Disable Notifications
            </button>

            <button v-if="isSubscribed" @click="sendTestNotification"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors font-medium">
              Test Notification
            </button>
          </div>

          <!-- Help Text -->
          <div class="text-sm text-gray-600 dark:text-gray-400">
            <p class="mb-2">
              <strong>What you'll receive:</strong>
            </p>
            <ul class="list-disc list-inside space-y-1 ml-4">
              <li>Reminders for upcoming tasks and deadlines</li>
              <li>Goal progress updates</li>
              <li>Study session suggestions</li>
              <li>Budget alerts and notifications</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
