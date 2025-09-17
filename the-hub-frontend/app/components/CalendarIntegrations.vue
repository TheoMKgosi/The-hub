<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCalendarStore } from '~/stores/calendar'

const calendarStore = useCalendarStore()
const { integrations, loading, fetchError } = calendarStore

const googleIntegration = computed(() =>
  integrations.find(i => i.provider === 'google')
)

async function connectGoogle() {
  await calendarStore.initiateGoogleAuth()
}

async function syncCalendar(integrationId: string) {
  await calendarStore.syncCalendar(integrationId)
}

async function disconnectCalendar(integrationId: string) {
  if (confirm('Are you sure you want to disconnect this calendar integration?')) {
    await calendarStore.deleteIntegration(integrationId)
  }
}

onMounted(() => {
  calendarStore.fetchIntegrations()
})
</script>

<template>
  <div class="calendar-integrations">
    <div class="header">
      <p>Sync your scheduled tasks with external calendars</p>
    </div>

    <div class="integrations-grid">
      <!-- Google Calendar Integration -->
      <div class="integration-card">
        <div class="card-header">
          <div class="provider-info">
            <div class="provider-icon google">
              <svg viewBox="0 0 24 24" width="24" height="24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
            </div>
            <div class="provider-details">
              <h3>Google Calendar</h3>
              <p>Sync with your Google Calendar</p>
            </div>
          </div>
        </div>

        <div class="card-actions">
          <button
            v-if="!googleIntegration"
            @click="connectGoogle"
            class="btn-primary"
            :disabled="loading"
          >
            Connect Google Calendar
          </button>

          <div v-else class="integration-controls">
            <div class="integration-status">
              <span class="status-indicator connected"></span>
              <span>Connected</span>
            </div>

            <div class="integration-actions">
              <button
                @click="syncCalendar(googleIntegration.id)"
                class="btn-secondary"
                :disabled="loading"
              >
                Sync Now
              </button>

              <button
                @click="disconnectCalendar(googleIntegration.id)"
                class="btn-danger"
              >
                Disconnect
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Outlook Calendar (Coming Soon) -->
      <div class="integration-card coming-soon">
        <div class="card-header">
          <div class="provider-info">
            <div class="provider-icon outlook">
              <svg viewBox="0 0 24 24" width="24" height="24">
                <rect fill="#0078D4" x="2" y="2" width="20" height="20" rx="2"/>
                <rect fill="#FFFFFF" x="6" y="6" width="12" height="12"/>
                <rect fill="#0078D4" x="8" y="8" width="8" height="8"/>
              </svg>
            </div>
            <div class="provider-details">
              <h3>Outlook Calendar</h3>
              <p>Coming soon</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Apple Calendar (Coming Soon) -->
      <div class="integration-card coming-soon">
        <div class="card-header">
          <div class="provider-info">
            <div class="provider-icon apple">
              <svg viewBox="0 0 24 24" width="24" height="24">
                <path fill="#000000" d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
              </svg>
            </div>
            <div class="provider-details">
              <h3>Apple Calendar</h3>
              <p>Coming soon</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="fetchError" class="error-message">
      <p>{{ fetchError.message }}</p>
    </div>
  </div>
</template>

<style scoped>
.calendar-integrations {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  text-align: center;
  margin-bottom: 2rem;
}

.header h2 {
  color: #1f2937;
  font-size: 1.875rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.header p {
  color: #6b7280;
  font-size: 1rem;
}

.integrations-grid {
  display: grid;
  gap: 1.5rem;
}

.integration-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  padding: 1.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
}

.coming-soon {
  opacity: 0.6;
  position: relative;
}

.coming-soon::after {
  content: 'Coming Soon';
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: #f59e0b;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.card-header {
  margin-bottom: 1.5rem;
}

.provider-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.provider-icon {
  width: 48px;
  height: 48px;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
}

.provider-icon.google {
  background: #ffffff;
  border: 1px solid #e5e7eb;
}

.provider-icon.outlook {
  background: #ffffff;
  border: 1px solid #e5e7eb;
}

.provider-icon.apple {
  background: #ffffff;
  border: 1px solid #e5e7eb;
}

.provider-details h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.provider-details p {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0.25rem 0 0 0;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
}

.btn-primary {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-right: 0.5rem;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e7eb;
}

.btn-secondary:disabled {
  background: #f9fafb;
  cursor: not-allowed;
}

.btn-danger {
  background: #ef4444;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-danger:hover {
  background: #dc2626;
}

.integration-controls {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.integration-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.connected {
  background: #10b981;
}

.integration-actions {
  display: flex;
  gap: 0.5rem;
}

.error-message {
  margin-top: 1rem;
  padding: 1rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  color: #dc2626;
}

.error-message p {
  margin: 0;
  font-size: 0.875rem;
}
</style>
