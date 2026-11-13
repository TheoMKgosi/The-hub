<script setup lang="ts">
import type { ActivityEntry } from '~/stores/activity-log'

const activityStore = useActivityLogStore()

// Reka UI (behind USelect) reserves '' as the "clear selection" sentinel, so
// SelectItems must not use an empty-string value. This sentinel stands in for
// "no category filter" and is stripped before hitting the API.
const ALL_CATEGORIES = '__all__'

const filters = reactive({
  search: '',
  category: ALL_CATEGORIES,
  date_from: '',
  date_to: ''
})

let searchDebounce: ReturnType<typeof setTimeout> | null = null

const applyFilters = () => {
  activityStore.fetchEntries({
    search: filters.search || undefined,
    category: filters.category === ALL_CATEGORIES ? undefined : filters.category || undefined,
    date_from: filters.date_from || undefined,
    date_to: filters.date_to || undefined,
    limit: 100
  })
}

watch(() => filters.search, () => {
  if (searchDebounce) clearTimeout(searchDebounce)
  searchDebounce = setTimeout(applyFilters, 300)
})

watch([() => filters.category, () => filters.date_from, () => filters.date_to], applyFilters)

const clearFilters = () => {
  filters.search = ''
  filters.category = ALL_CATEGORIES
  filters.date_from = ''
  filters.date_to = ''
  applyFilters()
}

const hasActiveFilters = computed(() =>
  filters.search || filters.category || filters.date_from || filters.date_to
)

const formatDayLabel = (dateKey: string) => {
  const today = new Date().toISOString().slice(0, 10)
  const yesterday = new Date(Date.now() - 86400000).toISOString().slice(0, 10)
  if (dateKey === today) return 'Today'
  if (dateKey === yesterday) return 'Yesterday'
  return new Date(`${dateKey}T00:00:00`).toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const dayTotal = (dayEntries: ActivityEntry[]) =>
  dayEntries.reduce((sum, entry) => sum + (entry.duration_min || 0), 0)

const formatDuration = (minutes: number) => {
  if (!minutes) return ''
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  if (hours && mins) return `${hours}h ${mins}m`
  if (hours) return `${hours}h`
  return `${mins}m`
}

const groupedEntries = computed(() => {
  const groups = new Map<string, ActivityEntry[]>()
  for (const entry of activityStore.entries) {
    const dateKey = entry.activity_date.slice(0, 10)
    if (!groups.has(dateKey)) groups.set(dateKey, [])
    groups.get(dateKey)!.push(entry)
  }
  return [...groups.entries()].map(([date, dayEntries]) => ({ date, entries: dayEntries }))
})

// Edit modal
const isEditOpen = ref(false)
const editingEntry = ref<ActivityEntry | null>(null)
const editForm = reactive({
  title: '',
  notes: '',
  category: '',
  tagsInput: '',
  activity_date: '',
  duration_min: 0 as number | null
})
const editErrors = ref<Record<string, string>>({})

const openEdit = (entry: ActivityEntry) => {
  editingEntry.value = entry
  editForm.title = entry.title
  editForm.notes = entry.notes
  editForm.category = entry.category
  editForm.tagsInput = (entry.tags || []).join(', ')
  editForm.activity_date = entry.activity_date.slice(0, 10)
  editForm.duration_min = entry.duration_min
  editErrors.value = {}
  isEditOpen.value = true
}

const saveEdit = async () => {
  if (!editingEntry.value) return
  const { validateObject } = useValidation()
  const validation = validateObject(editForm, { title: { required: true, maxLength: 255 } })
  editErrors.value = validation.errors
  if (!validation.isValid) return

  await activityStore.updateEntry(editingEntry.value.activity_id, {
    title: editForm.title.trim(),
    notes: editForm.notes.trim(),
    category: editForm.category.trim(),
    tags: editForm.tagsInput.split(',').map(tag => tag.trim()).filter(Boolean),
    activity_date: editForm.activity_date,
    duration_min: editForm.duration_min || 0
  })
  isEditOpen.value = false
}

const confirmDelete = async (entry: ActivityEntry) => {
  await activityStore.deleteEntry(entry.activity_id)
}

onMounted(applyFilters)
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <UInput
        v-model="filters.search"
        icon="i-lucide-search"
        placeholder="Search activities…"
        class="w-full sm:w-64"
      />
      <USelect
        v-model="filters.category"
        placeholder="All categories"
        :items="[{ label: 'All categories', value: ALL_CATEGORIES }, ...activityStore.categorySuggestions.map(c => ({ label: c, value: c }))]"
        class="w-44"
      />
      <UInput v-model="filters.date_from" type="date" class="w-40" />
      <span class="text-sm text-dimmed">to</span>
      <UInput v-model="filters.date_to" type="date" class="w-40" />
      <UButton
        v-if="hasActiveFilters"
        label="Clear"
        variant="ghost"
        icon="i-lucide-x"
        @click="clearFilters"
      />
    </div>

    <div v-if="activityStore.loading && !groupedEntries.length" class="py-8 text-center text-dimmed">
      Loading activities…
    </div>

    <div v-else-if="!groupedEntries.length" class="py-8 text-center text-dimmed">
      No activities logged yet. Add your first one above!
    </div>

    <div v-else class="space-y-6">
      <div v-for="group in groupedEntries" :key="group.date">
        <div class="mb-2 flex items-baseline gap-2">
          <h3 class="font-semibold">{{ formatDayLabel(group.date) }}</h3>
          <span v-if="dayTotal(group.entries)" class="text-sm text-dimmed">
            {{ formatDuration(dayTotal(group.entries)) }}
          </span>
        </div>

        <div class="space-y-2">
          <div
            v-for="entry in group.entries"
            :key="entry.activity_id"
            class="rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light/50 dark:bg-surface-dark/50 p-3 flex flex-wrap items-start justify-between gap-2"
          >
            <div class="min-w-0">
              <p class="font-medium truncate">{{ entry.title }}</p>
              <div class="mt-1 flex flex-wrap items-center gap-1.5">
                <UBadge v-if="entry.category" variant="subtle">{{ entry.category }}</UBadge>
                <UBadge v-for="tag in entry.tags" :key="tag" variant="subtle" color="neutral">{{ tag }}</UBadge>
                <span v-if="entry.duration_min" class="text-sm text-dimmed">
                  {{ formatDuration(entry.duration_min) }}
                </span>
              </div>
              <p v-if="entry.notes" class="mt-1.5 text-sm text-dimmed">{{ entry.notes }}</p>
            </div>

            <div class="flex gap-1">
              <UButton
                icon="i-lucide-pencil"
                variant="ghost"
                size="sm"
                aria-label="Edit activity"
                @click="openEdit(entry)"
              />
              <UButton
                icon="i-lucide-trash-2"
                variant="ghost"
                color="error"
                size="sm"
                aria-label="Delete activity"
                @click="confirmDelete(entry)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <UModal v-model:open="isEditOpen" title="Edit activity">
      <template #body>
        <UForm :state="editForm" class="space-y-3" @submit="saveEdit">
          <UFormField label="Title" required :error="editErrors.title">
            <UInput v-model="editForm.title" class="w-full" />
          </UFormField>
          <div class="grid gap-3 sm:grid-cols-2">
            <UFormField label="Category">
              <UInput v-model="editForm.category" class="w-full" />
            </UFormField>
            <UFormField label="Date">
              <UInput v-model="editForm.activity_date" type="date" class="w-full" />
            </UFormField>
          </div>
          <UFormField label="Tags (comma separated)">
            <UInput v-model="editForm.tagsInput" class="w-full" />
          </UFormField>
          <UFormField label="Duration (minutes)">
            <UInputNumber v-model="editForm.duration_min" :min="0" class="w-full" />
          </UFormField>
          <UFormField label="Notes">
            <UTextarea v-model="editForm.notes" class="w-full" :rows="3" />
          </UFormField>
        </UForm>
      </template>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UButton label="Cancel" variant="ghost" @click="isEditOpen = false" />
          <UButton label="Save" @click="saveEdit" />
        </div>
      </template>
    </UModal>
  </div>
</template>
