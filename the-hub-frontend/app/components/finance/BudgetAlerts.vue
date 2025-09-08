<script setup>
const budgetStore = useBudgetStore()
const { addToast } = useToast()

// Alert settings
const alertSettings = reactive({
  warningThreshold: 80,
  dangerThreshold: 90,
  overspendThreshold: 100,
  periodEndWarning: 3, // days before period end
  enableEmailAlerts: true,
  enablePushAlerts: false
})

// Save alert settings
const saveAlertSettings = async () => {
  try {
    // In a real app, this would save to backend
    addToast('Alert settings saved successfully', 'success')
  } catch (error) {
    addToast('Failed to save alert settings', 'error')
  }
}

// Mark alert as read
const markAlertAsRead = (alertId) => {
  // In a real app, this would update the alert status in backend
  addToast('Alert marked as read', 'success')
}

// Dismiss alert
const dismissAlert = (alertId) => {
  // In a real app, this would remove the alert from backend
  addToast('Alert dismissed', 'success')
}

onMounted(() => {
  if (budgetStore.alerts.length === 0) {
    budgetStore.fetchBudgetAlerts()
  }
})
</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 text-text-light dark:text-text-dark">Budget Alerts & Notifications</h2>

    <!-- Alert Settings -->
    <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Alert Settings</h3>
      <form @submit.prevent="saveAlertSettings" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
              Warning Threshold (%)
            </label>
            <input v-model.number="alertSettings.warningThreshold" type="number" min="1" max="99"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <p class="text-xs text-text-light dark:text-text-dark/60 mt-1">Alert when budget usage reaches this percentage</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
              Danger Threshold (%)
            </label>
            <input v-model.number="alertSettings.dangerThreshold" type="number" min="1" max="99"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <p class="text-xs text-text-light dark:text-text-dark/60 mt-1">Critical alert when budget usage reaches this percentage</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
              Period End Warning (days)
            </label>
            <input v-model.number="alertSettings.periodEndWarning" type="number" min="1" max="30"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <p class="text-xs text-text-light dark:text-text-dark/60 mt-1">Days before budget period ends to send warning</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">
              Overspend Threshold (%)
            </label>
            <input v-model.number="alertSettings.overspendThreshold" type="number" min="100" max="200"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
            <p class="text-xs text-text-light dark:text-text-dark/60 mt-1">Alert when spending exceeds budget by this percentage</p>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <label class="flex items-center">
            <input v-model="alertSettings.enableEmailAlerts" type="checkbox"
              class="rounded border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-primary focus:ring-primary" />
            <span class="ml-2 text-sm text-text-light dark:text-text-dark">Email notifications</span>
          </label>

          <label class="flex items-center">
            <input v-model="alertSettings.enablePushAlerts" type="checkbox"
              class="rounded border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-primary focus:ring-primary" />
            <span class="ml-2 text-sm text-text-light dark:text-text-dark">Push notifications</span>
          </label>
        </div>

        <div class="flex justify-end">
          <UiButton type="submit" variant="primary" size="md">
            Save Settings
          </UiButton>
        </div>
      </form>
    </div>

    <!-- Active Alerts -->
    <div v-if="budgetStore.alerts.length > 0" class="mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Active Alerts</h3>
      <div class="space-y-3">
        <div v-for="alert in budgetStore.alerts" :key="alert.budget_id"
          :class="[
            'p-4 rounded-lg border-l-4 flex items-start justify-between',
            alert.alert_type === 'over_budget' ? 'bg-red-50 dark:bg-red-900/20 border-red-500' :
            alert.alert_type === 'danger' ? 'bg-orange-50 dark:bg-orange-900/20 border-orange-500' :
            'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-500'
          ]">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-1">
              <h4 class="font-medium text-text-light dark:text-text-dark">{{ alert.category_name }}</h4>
              <span :class="[
                'px-2 py-1 text-xs font-medium rounded-full',
                alert.alert_type === 'over_budget' ? 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100' :
                alert.alert_type === 'danger' ? 'bg-orange-100 text-orange-800 dark:bg-orange-800 dark:text-orange-100' :
                'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100'
              ]">
                {{ alert.alert_type.replace('_', ' ').toUpperCase() }}
              </span>
            </div>
            <p class="text-sm text-text-light dark:text-text-dark/80">{{ alert.message }}</p>
            <div class="flex items-center gap-4 mt-2 text-xs text-text-light dark:text-text-dark/60">
              <span>Usage: {{ alert.current_usage.toFixed(1) }}%</span>
              <span>Budget: ${{ alert.budget_amount.toFixed(2) }}</span>
              <span>Days left: {{ alert.days_remaining }}</span>
            </div>
          </div>
          <div class="flex gap-2 ml-4">
            <UiButton variant="default" size="sm" @click="markAlertAsRead(alert.budget_id)">
              Mark Read
            </UiButton>
            <UiButton variant="danger" size="sm" @click="dismissAlert(alert.budget_id)">
              Dismiss
            </UiButton>
          </div>
        </div>
      </div>
    </div>

    <!-- Alert History -->
    <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Alert History</h3>
      <div v-if="budgetStore.alerts.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p>No alerts yet. Alerts will appear here when your budgets need attention.</p>
      </div>
      <div v-else class="space-y-2">
        <div v-for="alert in budgetStore.alerts.slice().reverse()" :key="`history-${alert.budget_id}`"
          class="flex items-center justify-between p-3 bg-surface-light/50 dark:bg-surface-dark/50 rounded-lg">
          <div>
            <p class="text-sm font-medium text-text-light dark:text-text-dark">{{ alert.category_name }}</p>
            <p class="text-xs text-text-light dark:text-text-dark/60">{{ alert.message }}</p>
          </div>
          <div class="text-right">
            <p class="text-xs text-text-light dark:text-text-dark/60">{{ new Date().toLocaleDateString() }}</p>
            <span :class="[
              'px-2 py-1 text-xs font-medium rounded-full',
              alert.alert_type === 'over_budget' ? 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100' :
              alert.alert_type === 'danger' ? 'bg-orange-100 text-orange-800 dark:bg-orange-800 dark:text-orange-100' :
              'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100'
            ]">
              {{ alert.alert_type.replace('_', ' ').toUpperCase() }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>