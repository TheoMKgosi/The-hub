<script setup lang="ts">
import ChevronDownIcon from '../ui/svg/ChevronDownIcon.vue'
import ChevronUpIcon from '../ui/svg/ChevronUpIcon.vue'
import FilterIcon from '../ui/svg/FilterIcon.vue'
import type { TaskFilters } from '~/composables/useTaskFilters'

const props = defineProps<{
  filters: TaskFilters
  activeFilterCount: number
  activeFilterChips: { key: keyof TaskFilters; label: string; value: string }[]
  isAdvancedOpen: boolean
  isFilterBarOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'update:filters', filters: TaskFilters): void
  (e: 'update:isAdvancedOpen', value: boolean): void
  (e: 'update:isFilterBarOpen', value: boolean): void
  (e: 'clearFilter', key: keyof TaskFilters): void
  (e: 'clearAllFilters'): void
}>()

// Status options
const statusOptions = [
  { value: 'pending', label: 'Pending', color: '#F59E0B' },
  { value: 'in_progress', label: 'In Progress', color: '#3B82F6' },
  { value: 'completed', label: 'Completed', color: '#10B981' }
]

// Priority options
const priorityOptions = [
  { value: '1', label: '1 - Low', color: '#6B7280' },
  { value: '2', label: '2 - Low-Med', color: '#9CA3AF' },
  { value: '3', label: '3 - Medium', color: '#F59E0B' },
  { value: '4', label: '4 - High-Med', color: '#EF4444' },
  { value: '5', label: '5 - High', color: '#DC2626' }
]

// Due date options
const dueDateOptions = [
  { value: 'today', label: 'Today' },
  { value: 'tomorrow', label: 'Tomorrow' },
  { value: 'this_week', label: 'This Week' },
  { value: 'overdue', label: 'Overdue' },
  { value: 'none', label: 'No Due Date' }
]

// Linked options
const linkedOptions = [
  { value: 'true', label: 'Linked to Goals' },
  { value: 'false', label: 'Not Linked' }
]

// Duration options
const durationOptions = [
  { value: '0-30', label: '0-30 min' },
  { value: '30-60', label: '30-60 min' },
  { value: '60-120', label: '1-2 hours' },
  { value: '120+', label: '2+ hours' }
]

// Created options
const createdOptions = [
  { value: '7d', label: 'Last 7 days' },
  { value: '30d', label: 'Last 30 days' },
  { value: 'older', label: 'Older' }
]

// Schedule options
const scheduleOptions = [
  { value: 'true', label: 'Has Schedule' },
  { value: 'false', label: 'No Schedule' }
]

// Update filter helper
const updateFilter = (key: keyof TaskFilters, value: string | string[]) => {
  emit('update:filters', { ...props.filters, [key]: value })
}
</script>

<template>
  <div class="sticky top-0 z-30 bg-surface-light/95 dark:bg-surface-dark/95 backdrop-blur-md border-b border-surface-light/20 dark:border-surface-dark/20 shadow-sm">
    <!-- Main Filter Bar -->
    <div class="p-4">
      <!-- Header with toggle -->
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-2">
          <FilterIcon class="w-5 h-5 text-primary" />
          <span class="font-semibold text-text-light dark:text-text-dark">Filters</span>
          <span
            v-if="activeFilterCount > 0"
            class="px-2 py-0.5 text-xs font-medium bg-primary text-white rounded-full"
          >
            {{ activeFilterCount }}
          </span>
        </div>
        <button
          @click="emit('update:isFilterBarOpen', !isFilterBarOpen)"
          class="p-1.5 rounded-md hover:bg-surface-light/50 dark:hover:bg-surface-dark/50 transition-colors"
        >
          <ChevronUpIcon v-if="isFilterBarOpen" class="w-5 h-5 text-text-light/60 dark:text-text-dark/60" />
          <ChevronDownIcon v-else class="w-5 h-5 text-text-light/60 dark:text-text-dark/60" />
        </button>
      </div>

      <!-- Filter Controls -->
      <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 max-h-0"
        enter-to-class="opacity-100 max-h-[500px]"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 max-h-[500px]"
        leave-to-class="opacity-0 max-h-0"
      >
        <div v-show="isFilterBarOpen" class="overflow-hidden">
          <!-- Primary Filters Row -->
          <div class="flex flex-wrap items-center gap-3 mb-3">
            <!-- Search -->
            <FilterSearch
              :model-value="filters.search || ''"
              @update:model-value="updateFilter('search', $event)"
            />

            <!-- Status Dropdown -->
            <FilterDropdown
              label="Status"
              :options="statusOptions"
              :model-value="filters.status || ''"
              @update:model-value="updateFilter('status', $event)"
            />

            <!-- Priority Dropdown (Multi-select) -->
            <FilterDropdown
              label="Priority"
              :options="priorityOptions"
              :model-value="filters.priority ? filters.priority.split(',') : []"
              :multiple="true"
              @update:model-value="updateFilter('priority', Array.isArray($event) ? $event.join(',') : $event)"
            />

            <!-- Due Date Dropdown -->
            <FilterDropdown
              label="Due Date"
              :options="dueDateOptions"
              :model-value="filters.dueDate || ''"
              @update:model-value="updateFilter('dueDate', $event)"
            />

            <!-- Linked Dropdown -->
            <FilterDropdown
              label="Goals"
              :options="linkedOptions"
              :model-value="filters.linked || ''"
              @update:model-value="updateFilter('linked', $event)"
            />

            <!-- Sort -->
            <FilterSort
              :model-value="filters.sortBy || 'due_date'"
              :order="filters.sortOrder || 'asc'"
              @update:model-value="updateFilter('sortBy', $event)"
              @update:order="updateFilter('sortOrder', $event)"
            />

            <!-- Advanced Toggle -->
            <button
              @click="emit('update:isAdvancedOpen', !isAdvancedOpen)"
              class="inline-flex items-center gap-1 px-3 py-2 text-sm font-medium rounded-md border transition-all duration-200"
              :class="[
                isAdvancedOpen
                  ? 'bg-primary/10 border-primary text-primary'
                  : 'bg-surface-light dark:bg-surface-dark border-surface-light/30 dark:border-surface-dark/30 text-text-light dark:text-text-dark hover:border-primary/50'
              ]"
            >
              <span>Advanced</span>
              <ChevronDownIcon
                class="w-4 h-4 transition-transform duration-200"
                :class="{ 'rotate-180': isAdvancedOpen }"
              />
            </button>
          </div>

          <!-- Advanced Filters -->
          <Transition
            enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="opacity-0 max-h-0"
            enter-to-class="opacity-100 max-h-[200px]"
            leave-active-class="transition-all duration-200 ease-in"
            leave-from-class="opacity-100 max-h-[200px]"
            leave-to-class="opacity-0 max-h-0"
          >
            <div v-show="isAdvancedOpen" class="overflow-hidden">
              <div class="pt-3 border-t border-surface-light/20 dark:border-surface-dark/20">
                <div class="flex flex-wrap items-center gap-3">
                  <!-- Duration -->
                  <FilterDropdown
                    label="Duration"
                    :options="durationOptions"
                    :model-value="filters.duration || ''"
                    @update:model-value="updateFilter('duration', $event)"
                  />

                  <!-- Created Date -->
                  <FilterDropdown
                    label="Created"
                    :options="createdOptions"
                    :model-value="filters.created || ''"
                    @update:model-value="updateFilter('created', $event)"
                  />

                  <!-- Schedule Status -->
                  <FilterDropdown
                    label="Schedule"
                    :options="scheduleOptions"
                    :model-value="filters.hasStartTime || ''"
                    @update:model-value="updateFilter('hasStartTime', $event)"
                  />
                </div>
              </div>
            </div>
          </Transition>

          <!-- Active Filter Chips -->
          <div v-if="activeFilterChips.length > 0" class="mt-3 pt-3 border-t border-surface-light/20 dark:border-surface-dark/20">
            <FilterChips
              :chips="activeFilterChips"
              @remove="emit('clearFilter', $event)"
              @clear-all="emit('clearAllFilters')"
            />
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>
