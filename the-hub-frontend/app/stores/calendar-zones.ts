import { defineStore } from 'pinia'
import { ref } from 'vue'

interface CalendarZone {
  id: string
  user_id: string
  name: string
  description: string
  category: string
  color: string
  start_time: string
  end_time: string
  days_of_week: string
  priority: number
  is_active: boolean
  allow_scheduling: boolean
  max_events_per_day?: number
  is_recurring: boolean
  recurrence_start?: string
  recurrence_end?: string
  created_at: string
  updated_at: string
}

interface ZoneCategory {
  id: string
  name: string
  description: string
  color: string
  icon: string
  is_default: boolean
}

export const useCalendarZonesStore = defineStore('calendarZones', () => {
  const zones = ref<CalendarZone[]>([])
  const categories = ref<ZoneCategory[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Fetch all zones for the current user
  async function fetchZones() {
    const { $api } = useNuxtApp()
    loading.value = true
    error.value = null

    try {
      const response = await $api<{ zones: CalendarZone[] }>('/calendar-zones')
      zones.value = response.zones || []
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch calendar zones'
      console.error('Error fetching calendar zones:', err)
    } finally {
      loading.value = false
    }
  }

  // Fetch available zone categories
  async function fetchCategories() {
    const { $api } = useNuxtApp()
    try {
      const response = await $api<{ categories: ZoneCategory[] }>('/calendar-zones/categories')
      categories.value = response.categories || []
    } catch (err: any) {
      console.error('Error fetching zone categories:', err)
    }
  }

  // Create a new calendar zone
  async function createZone(zoneData: Partial<CalendarZone>) {
    const { $api } = useNuxtApp()
    try {
      const newZone = await $api<CalendarZone>('/calendar-zones', {
        method: 'POST',
        body: JSON.stringify(zoneData)
      })

      zones.value.push(newZone)
      return newZone
    } catch (err: any) {
      error.value = err.message || 'Failed to create calendar zone'
      throw err
    }
  }

  // Update an existing calendar zone
  async function updateZone(zoneId: string, zoneData: Partial<CalendarZone>) {
    const { $api } = useNuxtApp()
    try {
      const updatedZone = await $api<CalendarZone>(`/calendar-zones/${zoneId}`, {
        method: 'PUT',
        body: JSON.stringify(zoneData)
      })

      const index = zones.value.findIndex(z => z.id === zoneId)
      if (index !== -1) {
        zones.value[index] = updatedZone
      }

      return updatedZone
    } catch (err: any) {
      error.value = err.message || 'Failed to update calendar zone'
      throw err
    }
  }

  // Delete a calendar zone
  async function deleteZone(zoneId: string) {
    const { $api } = useNuxtApp()
    try {
      await $api(`/calendar-zones/${zoneId}`, {
        method: 'DELETE'
      })

      zones.value = zones.value.filter(z => z.id !== zoneId)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete calendar zone'
      throw err
    }
  }

  // Get zones for a specific day
  function getZonesForDay(date: Date): CalendarZone[] {
    const dayName = date.toLocaleLowerCase('en-US', { weekday: 'long' })

    return zones.value.filter(zone => {
      if (!zone.is_active) return false

      // Check if zone applies to this day
      if (zone.days_of_week) {
        const daysArray = zone.days_of_week.toLowerCase().split(',')
        return daysArray.some(day => day.trim() === dayName)
      }

      return true // If no days specified, applies to all days
    })
  }

  // Get zones for a specific time range
  function getZonesForTimeRange(startTime: Date, endTime: Date): CalendarZone[] {
    return zones.value.filter(zone => {
      if (!zone.is_active) return false

      const zoneStart = new Date(startTime.toDateString() + ' ' + zone.start_time)
      const zoneEnd = new Date(startTime.toDateString() + ' ' + zone.end_time)

      return startTime >= zoneStart && endTime <= zoneEnd
    })
  }

  // Get the best zone for a specific time and category
  function getBestZoneForTime(time: Date, category?: string): CalendarZone | null {
    const applicableZones = zones.value.filter(zone => {
      if (!zone.is_active || !zone.allow_scheduling) return false

      // Check time
      const zoneStart = new Date(time.toDateString() + ' ' + zone.start_time)
      const zoneEnd = new Date(time.toDateString() + ' ' + zone.end_time)

      if (time < zoneStart || time > zoneEnd) return false

      // Check category match if specified
      if (category && zone.category !== category) return false

      return true
    })

    if (applicableZones.length === 0) return null

    // Return zone with highest priority
    return applicableZones.reduce((best, current) =>
      current.priority > best.priority ? current : best
    )
  }

  // Get zone statistics
  function getZoneStats() {
    const stats = {
      total: zones.value.length,
      active: zones.value.filter(z => z.is_active).length,
      byCategory: {} as Record<string, number>,
      averagePriority: 0
    }

    let totalPriority = 0
    zones.value.forEach(zone => {
      if (zone.is_active) {
        stats.byCategory[zone.category] = (stats.byCategory[zone.category] || 0) + 1
        totalPriority += zone.priority
      }
    })

    stats.averagePriority = stats.active > 0 ? totalPriority / stats.active : 0

    return stats
  }

  return {
    zones,
    categories,
    loading,
    error,
    fetchZones,
    fetchCategories,
    createZone,
    updateZone,
    deleteZone,
    getZonesForDay,
    getZonesForTimeRange,
    getBestZoneForTime,
    getZoneStats
  }
})