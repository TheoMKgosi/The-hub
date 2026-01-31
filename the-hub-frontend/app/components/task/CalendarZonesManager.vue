<script setup lang="ts">

const zonesStore = useCalendarZonesStore()
const showManager = ref(false)
const showCreateForm = ref(false)
const editingZone = ref(null)

const newZone = reactive({
  name: '',
  description: '',
  category: 'work',
  color: '#3b82f6',
  start_time: '09:00',
  end_time: '17:00',
  days_of_week: ['monday', 'tuesday', 'wednesday', 'thursday', 'friday'] as string[],
  priority: 5,
  is_active: true,
  allow_scheduling: false,
  max_events_per_day: null as number | null,
  is_recurring: false,
  recurrence_start: null as string | null,
  recurrence_end: null as string | null,
})

const dayOptions = [
  { value: 'monday', label: 'Monday' },
  { value: 'tuesday', label: 'Tuesday' },
  { value: 'wednesday', label: 'Wednesday' },
  { value: 'thursday', label: 'Thursday' },
  { value: 'friday', label: 'Friday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' },
]

const priorityOptions = [
  { value: 1, label: '1 - Very Low' },
  { value: 2, label: '2 - Low' },
  { value: 3, label: '3 - Medium' },
  { value: 4, label: '4 - High' },
  { value: 5, label: '5 - Very High' },
]

const loadData = async () => {
  await Promise.all([
    zonesStore.fetchZones(),
    zonesStore.fetchCategories()
  ])
}

const openManager = async () => {
  showManager.value = true
  await loadData()
}

const closeManager = () => {
  showManager.value = false
  showCreateForm.value = false
  editingZone.value = null
  resetForm()
}

const resetForm = () => {
  Object.assign(newZone, {
    name: '',
    description: '',
    category: 'work',
    color: '#3b82f6',
    start_time: '09:00',
    end_time: '17:00',
    days_of_week: ['monday', 'tuesday', 'wednesday', 'thursday', 'friday'],
    priority: 5,
    is_active: true,
    allow_scheduling: false,
    max_events_per_day: null,
    is_recurring: false,
    recurrence_start: null,
    recurrence_end: null,
  })
}

const createZone = async () => {
  if (!newZone.name.trim()) return

  try {
    const zoneData = {
      ...newZone,
      days_of_week: newZone.days_of_week.join(','),
      start_time: new Date(`1970-01-01T${newZone.start_time}:00`),
      end_time: new Date(`1970-01-01T${newZone.end_time}:00`),
    }

    await zonesStore.createZone(zoneData)
    resetForm()
    showCreateForm.value = false
  } catch (error) {
    console.error('Failed to create zone:', error)
  }
}

const startEdit = (zone: any) => {
  editingZone.value = zone
  Object.assign(newZone, {
    name: zone.name,
    description: zone.description || '',
    category: zone.category,
    color: zone.color,
    start_time: new Date(zone.start_time).toTimeString().slice(0, 5),
    end_time: new Date(zone.end_time).toTimeString().slice(0, 5),
    days_of_week: zone.days_of_week ? zone.days_of_week.split(',').map((d: string) => d.trim()) : [],
    priority: zone.priority,
    is_active: zone.is_active,
    allow_scheduling: zone.allow_scheduling,
    max_events_per_day: zone.max_events_per_day,
    is_recurring: zone.is_recurring,
    recurrence_start: zone.recurrence_start,
    recurrence_end: zone.recurrence_end,
  })
  showCreateForm.value = true
}

const updateZone = async () => {
  if (!editingZone.value || !newZone.name.trim()) return

  try {
    const zoneData = {
      ...newZone,
      days_of_week: newZone.days_of_week.join(','),
      start_time: new Date(`1970-01-01T${newZone.start_time}:00`),
      end_time: new Date(`1970-01-01T${newZone.end_time}:00`),
    }

    await zonesStore.updateZone(editingZone.value.id, zoneData)
    resetForm()
    showCreateForm.value = false
    editingZone.value = null
  } catch (error) {
    console.error('Failed to update zone:', error)
  }
}

const deleteZone = async (zoneId: string) => {
  if (!confirm('Are you sure you want to delete this calendar zone?')) return

  try {
    await zonesStore.deleteZone(zoneId)
  } catch (error) {
    console.error('Failed to delete zone:', error)
  }
}

const toggleZone = async (zone: any) => {
  try {
    await zonesStore.updateZone(zone.id, { is_active: !zone.is_active })
  } catch (error) {
    console.error('Failed to toggle zone:', error)
  }
}

const getCategoryInfo = (categoryName: string) => {
  return zonesStore.categories.find(cat => cat.name === categoryName) || {
    name: categoryName,
    color: '#6b7280',
    icon: 'calendar'
  }
}

const formatDaysOfWeek = (daysInput: string | string[]) => {
  let days: string[]
  if (Array.isArray(daysInput)) {
    days = daysInput
  } else {
    if (!daysInput) return 'All days'
    days = daysInput.split(',').map(d => d.trim())
  }

  if (days.length === 7) return 'All days'
  if (days.length === 5 && days.includes('monday') && days.includes('friday')) return 'Weekdays'
  if (days.length === 2 && days.includes('saturday') && days.includes('sunday')) return 'Weekends'
  return days.map(d => d.charAt(0).toUpperCase() + d.slice(1)).join(', ')
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div>
    <!-- Trigger Button -->
    <BaseButton @click="openManager" variant="outline" size="sm">
      <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd"
          d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
          clip-rule="evenodd" />
      </svg>
      Calendar Zones
    </BaseButton>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showManager"
        class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50"
        @click="closeManager">
        <div
          class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-6xl max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light/30 dark:border-surface-dark/30"
          @click.stop>

          <!-- Header -->
          <div
            class="flex items-center justify-between p-6 border-b border-surface-light/20 dark:border-surface-dark/20">
            <div>
              <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Calendar Zones Manager</h2>
              <p class="text-sm text-text-light/60 dark:text-text-dark/60 mt-1">
                Define time zones for different activities to help AI make better scheduling suggestions
              </p>
            </div>
            <BaseButton @click="closeManager" variant="default" size="sm" class="p-2">
              ×
            </BaseButton>
          </div>

          <div class="p-6">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

              <!-- Zone List -->
              <div class="lg:col-span-2">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="text-lg font-medium text-text-light dark:text-text-dark">Your Zones</h3>
                  <BaseButton @click="showCreateForm = true" variant="primary" size="sm">
                    <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd"
                        d="M10 3a1 1 0 011 1v5h5a1 1 0 102 0v-5a1 1 0 01-1-1H10zM9 3a1 1 0 00-1 1v5H3a1 1 0 00-1 1v5a1 1 0 001 1h5v5a1 1 0 102 0v-5a1 1 0 001-1v-5H9V4a1 1 0 00-1-1z"
                        clip-rule="evenodd" />
                    </svg>
                    New Zone
                  </BaseButton>
                </div>

                <div v-if="zonesStore.loading" class="text-center py-8">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
                  <p class="text-sm text-text-light/60 dark:text-text-dark/60 mt-2">Loading zones...</p>
                </div>

                <div v-else-if="zonesStore.zones.length === 0"
                  class="text-center py-8 text-text-light/60 dark:text-text-dark/60">
                  <svg class="w-12 h-12 mx-auto mb-4 opacity-50" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd"
                      d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                      clip-rule="evenodd" />
                  </svg>
                  <p class="text-lg mb-2">No calendar zones yet</p>
                  <p class="text-sm">Create your first zone to get started</p>
                </div>

                <div v-else class="space-y-3">
                  <div v-for="zone in zonesStore.zones" :key="zone.id"
                    class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/30 dark:border-surface-dark/30">

                    <div class="flex items-center justify-between mb-3">
                      <div class="flex items-center gap-3">
                        <div class="w-4 h-4 rounded" :style="{ backgroundColor: zone.color }"></div>
                        <div>
                          <h4 class="font-semibold text-text-light dark:text-text-dark">{{ zone.name }}</h4>
                          <p class="text-sm text-text-light/60 dark:text-text-dark/60">{{
                            getCategoryInfo(zone.category).description }}</p>
                        </div>
                      </div>
                      <div class="flex items-center gap-2">
                        <span :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          zone.is_active ? 'bg-green-100 dark:bg-green-900/20 text-green-700 dark:text-green-300' : 'bg-gray-100 dark:bg-gray-900/20 text-gray-600 dark:text-gray-400'
                        ]">
                          {{ zone.is_active ? 'Active' : 'Inactive' }}
                        </span>
                        <BaseButton @click="toggleZone(zone)" variant="default" size="md"
                          :title="zone.is_active ? 'Deactivate' : 'Activate'">
                          <svg v-if="zone.is_active" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                              clip-rule="evenodd" />
                          </svg>
                          <svg v-else class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                              clip-rule="evenodd" />
                          </svg>
                        </BaseButton>
                      </div>
                    </div>

                    <div class="grid grid-cols-2 gap-4 mb-3">
                      <div>
                        <p class="text-sm font-medium text-text-light dark:text-text-dark">Time Range</p>
                        <p class="text-sm text-text-light/80 dark:text-text-dark/80">
                          {{ new Date(zone.start_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
                          -
                          {{ new Date(zone.end_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}
                        </p>
                      </div>
                      <div>
                        <p class="text-sm font-medium text-text-light dark:text-text-dark">Days</p>
                        <p class="text-sm text-text-light/80 dark:text-text-dark/80">{{
                          formatDaysOfWeek(zone.days_of_week) }}</p>
                      </div>
                    </div>

                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-4">
                        <div class="flex items-center gap-1">
                          <span class="text-sm text-text-light/60 dark:text-text-dark/60">Priority:</span>
                          <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ zone.priority
                            }}/10</span>
                        </div>
                        <div v-if="zone.allow_scheduling" class="flex items-center gap-1">
                          <svg class="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                              clip-rule="evenodd" />
                          </svg>
                          <span class="text-sm text-green-600 dark:text-green-400">AI Scheduling</span>
                        </div>
                      </div>
                      <div class="flex gap-2">
                        <BaseButton @click="startEdit(zone)" text="Edit" variant="default" size="md">Edit</BaseButton>
                        <BaseButton @click="deleteZone(zone.id)" text="Delete" variant="danger" size="md">Delete
                        </BaseButton>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Create/Edit Form -->
              <div class="lg:col-span-1">
                <div v-if="showCreateForm"
                  class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/30 dark:border-surface-dark/30">
                  <h3 class="text-lg font-medium mb-4 text-text-light dark:text-text-dark">
                    {{ editingZone ? 'Edit Zone' : 'Create New Zone' }}
                  </h3>

                  <form @submit.prevent="editingZone ? updateZone() : createZone()" class="space-y-4">
                    <div>
                      <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Zone
                        Name</label>
                      <input v-model="newZone.name" type="text" placeholder="e.g., Work Hours, Class Time"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                        required />
                    </div>

                    <div>
                      <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Category</label>
                      <select v-model="newZone.category"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                        <option v-for="category in zonesStore.categories" :key="category.name" :value="category.name">
                          {{ category.name.charAt(0).toUpperCase() + category.name.slice(1) }}
                        </option>
                      </select>
                    </div>

                    <div>
                      <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Color</label>
                      <input v-model="newZone.color" type="color"
                        class="w-full h-10 border border-surface-light/30 dark:border-surface-dark/30 rounded-md cursor-pointer" />
                    </div>

                    <div class="grid grid-cols-2 gap-3">
                      <div>
                        <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Start
                          Time</label>
                        <input v-model="newZone.start_time" type="time"
                          class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                          required />
                      </div>
                      <div>
                        <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">End
                          Time</label>
                        <input v-model="newZone.end_time" type="time"
                          class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                          required />
                      </div>
                    </div>

                    <div>
                      <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Days of
                        Week</label>
                      <div class="grid grid-cols-2 gap-2">
                        <label v-for="day in dayOptions" :key="day.value" class="flex items-center">
                          <input v-model="newZone.days_of_week" type="checkbox" :value="day.value" class="mr-2" />
                          <span class="text-sm text-text-light dark:text-text-dark">{{ day.label }}</span>
                        </label>
                      </div>
                    </div>

                    <div>
                      <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Priority
                        (1-10)</label>
                      <select v-model="newZone.priority"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                        <option v-for="option in priorityOptions" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </option>
                      </select>
                    </div>

                    <div class="flex items-center">
                      <input v-model="newZone.allow_scheduling" type="checkbox"
                        :id="'allow-scheduling-' + (editingZone?.id || 'new')"
                        class="w-4 h-4 mr-3 rounded border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-primary focus:ring-2 focus:ring-primary focus:ring-offset-0 focus:outline-none cursor-pointer" />
                      <label :for="'allow-scheduling-' + (editingZone?.id || 'new')"
                        class="text-sm text-text-light dark:text-text-dark cursor-pointer">Allow AI to schedule events
                        in this zone</label>
                    </div>

                    <div class="flex gap-2">
                      <BaseButton type="submit" :text="editingZone ? 'Update Zone' : 'Create Zone'" variant="primary"
                        size="sm" class="flex-1">
                      </BaseButton>
                      <BaseButton @click="showCreateForm = false; editingZone = null; resetForm()" text="Cancel"
                        variant="default" size="sm">
                      </BaseButton>
                    </div>
                  </form>
                </div>

                <!-- Zone Statistics -->
                <div v-else
                  class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4 border border-surface-light/30 dark:border-surface-dark/30">
                  <h3 class="text-lg font-medium mb-4 text-text-light dark:text-text-dark">Zone Statistics</h3>

                  <div class="space-y-3">
                    <div class="flex justify-between items-center">
                      <span class="text-sm text-text-light dark:text-text-dark">Total Zones</span>
                      <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ zonesStore.zones.length
                        }}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-sm text-text-light dark:text-text-dark">Active Zones</span>
                      <span class="text-sm font-medium text-green-600 dark:text-green-400">
                        {{zonesStore.zones.filter(z => z.is_active).length}}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-sm text-text-light dark:text-text-dark">AI Scheduling Enabled</span>
                      <span class="text-sm font-medium text-blue-600 dark:text-blue-400">{{zonesStore.zones.filter(z =>
                        z.allow_scheduling).length}}</span>
                    </div>
                    <div class="flex justify-between items-center">
                      <span class="text-sm text-text-light dark:text-text-dark">Average Priority</span>
                      <span class="text-sm font-medium text-text-light dark:text-text-dark">
                        {{zonesStore.zones.length > 0 ? (zonesStore.zones.reduce((sum, z) => sum + z.priority, 0) /
                          zonesStore.zones.length).toFixed(1) : '0'}}/10
                      </span>
                    </div>
                  </div>

                  <div class="mt-4 pt-4 border-t border-surface-light/20 dark:border-surface-dark/20">
                    <h4 class="text-sm font-medium mb-2 text-text-light dark:text-text-dark">Categories</h4>
                    <div class="space-y-1">
                      <div v-for="category in zonesStore.categories.slice(0, 5)" :key="category.name"
                        class="flex justify-between items-center">
                        <span class="text-xs text-text-light/80 dark:text-text-dark/80">{{ category.name }}</span>
                        <span class="text-xs font-medium text-text-light dark:text-text-dark">
                          {{zonesStore.zones.filter(z => z.category === category.name).length}}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
