<script setup lang="ts">
import type { Theme, ThemeWeekday, ThemeSuggestion } from '~/types/theme'
import { WEEKDAYS, WEEKDAY_LABELS } from '~/types/theme'

const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ close: [] }>()

const themeStore = useThemeStore()
const taskStore = useTaskStore()
const toast = useToast()

const open = computed({
  get: () => props.show,
  set: (value) => { if (!value) emit('close') }
})

const DAYS = WEEKDAYS as ThemeWeekday[]

// Form state
const showForm = ref(false)
const editingTheme = ref<Theme | null>(null)
const form = reactive({
  name: '',
  description: '',
  color: '#3B82F6',
  days: [] as ThemeWeekday[]
})

const resetForm = () => {
  editingTheme.value = null
  form.name = ''
  form.description = ''
  form.color = '#3B82F6'
  form.days = []
  showForm.value = false
}

const startCreate = () => {
  resetForm()
  showForm.value = true
}

const startEdit = (theme: Theme) => {
  editingTheme.value = theme
  form.name = theme.name
  form.description = theme.description
  form.color = theme.color
  form.days = [...theme.days]
  showForm.value = true
}

const toggleDay = (day: ThemeWeekday) => {
  const index = form.days.indexOf(day)
  if (index === -1) {
    form.days.push(day)
  } else {
    form.days.splice(index, 1)
  }
}

const submitForm = async () => {
  if (!form.name.trim()) {
    toast.add({ title: 'Theme', description: 'Theme name is required', color: 'error' })
    return
  }
  if (form.days.length === 0) {
    toast.add({ title: 'Theme', description: 'Select at least one day', color: 'error' })
    return
  }

  const payload = {
    name: form.name.trim(),
    description: form.description,
    color: form.color,
    days: form.days
  }

  try {
    if (editingTheme.value) {
      await themeStore.updateTheme(editingTheme.value.theme_id, payload)
    } else {
      await themeStore.createTheme(payload)
    }
    resetForm()
  } catch (error) {
    // Toast handled in store
  }
}

const handleDelete = async (theme: Theme) => {
  await themeStore.deleteTheme(theme.theme_id)
}

// AI suggestions
const selectedSuggestions = ref<Set<string>>(new Set())

const runAISuggestions = async () => {
  await themeStore.getAISuggestions()
  selectedSuggestions.value = new Set(themeStore.suggestions.map(s => s.name))
}

const toggleSuggestion = (name: string) => {
  if (selectedSuggestions.value.has(name)) {
    selectedSuggestions.value.delete(name)
  } else {
    selectedSuggestions.value.add(name)
  }
}

const suggestionTaskCount = (suggestion: ThemeSuggestion) => {
  return suggestion.tasks.length
}

const applySuggestions = async () => {
  const selected = themeStore.suggestions.filter(s => selectedSuggestions.value.has(s.name))
  if (selected.length === 0) {
    toast.add({ title: 'Theme', description: 'Select at least one theme to apply', color: 'error' })
    return
  }

  const themes = selected.map(s => ({
    name: s.name,
    description: s.description,
    color: s.color,
    days: s.days
  }))

  const assignments = selected.flatMap(s =>
    s.tasks.map(taskId => ({ task_id: taskId, theme_name: s.name }))
  )

  const success = await themeStore.applyAISuggestions({ themes, assignments })
  if (success) {
    await taskStore.fetchTasks()
  }
}

onMounted(() => {
  if (themeStore.themes.length === 0) themeStore.fetchThemes()
  themeStore.fetchEnabled()
})
</script>

<template>
  <UModal v-model:open="open" title="Theme Days">
    <template #body>
      <div class="space-y-6">
        <!-- Enable toggle -->
        <div class="flex items-center justify-between gap-4 p-4 rounded-lg bg-surface-light/50 dark:bg-surface-dark/50 border border-surface-light/20 dark:border-surface-dark/20">
          <div>
            <p class="text-sm font-medium text-text-light dark:text-text-dark">Enable Theme Days</p>
            <p class="text-xs text-text-light/60 dark:text-text-dark/60 mt-0.5">
              Only show tasks belonging to today's theme
            </p>
          </div>
          <ThemeDaysToggle />
        </div>

        <!-- AI Suggestions -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-text-light dark:text-text-dark">AI Theme Suggestions</h3>
            <UButton
              label="Suggest Themes"
              icon="i-lucide-bot"
              variant="outline"
              size="sm"
              :loading="themeStore.suggestLoading"
              @click="runAISuggestions"
            />
          </div>

          <p v-if="themeStore.suggestions.length === 0" class="text-xs text-text-light/60 dark:text-text-dark/60">
            Let the AI group your tasks into themes with weekday assignments.
          </p>

          <div v-else class="space-y-3">
            <div
              v-for="suggestion in themeStore.suggestions"
              :key="suggestion.name"
              class="border border-surface-light/20 dark:border-surface-dark/20 rounded-lg p-3"
              :class="selectedSuggestions.has(suggestion.name) ? 'border-primary bg-primary/5' : 'opacity-70'"
            >
              <div class="flex items-start gap-3">
                <input
                  type="checkbox"
                  :checked="selectedSuggestions.has(suggestion.name)"
                  @change="toggleSuggestion(suggestion.name)"
                  class="h-4 w-4 mt-1 rounded border-surface-light dark:border-surface-dark text-primary focus:ring-primary"
                />
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: suggestion.color }" />
                    <span class="text-sm font-medium text-text-light dark:text-text-dark">{{ suggestion.name }}</span>
                    <span class="text-xs text-text-light/50 dark:text-text-dark/50">{{ suggestionTaskCount(suggestion) }} tasks</span>
                  </div>
                  <p class="text-xs text-text-light/70 dark:text-text-dark/70 mt-1">{{ suggestion.description }}</p>
                  <div class="flex flex-wrap gap-1 mt-2">
                    <span
                      v-for="day in suggestion.days"
                      :key="day"
                      class="px-2 py-0.5 text-xs rounded bg-secondary/10 dark:bg-secondary/20 text-secondary"
                    >
                      {{ WEEKDAY_LABELS[day] }}
                    </span>
                  </div>
                  <p class="text-xs italic text-text-light/50 dark:text-text-dark/50 mt-2">{{ suggestion.reasoning }}</p>
                </div>
              </div>
            </div>

            <UButton
              label="Apply Selected Themes"
              color="primary"
              size="sm"
              :loading="themeStore.suggestLoading"
              :disabled="selectedSuggestions.size === 0"
              @click="applySuggestions"
            />
          </div>
        </div>

        <div class="border-t border-surface-light/20 dark:border-surface-dark/20" />

        <!-- Theme list -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-text-light dark:text-text-dark">My Themes</h3>
            <UButton label="Add Theme" icon="i-lucide-plus" variant="outline" size="sm" @click="startCreate" />
          </div>

          <p v-if="themeStore.themes.length === 0" class="text-xs text-text-light/60 dark:text-text-dark/60">
            No themes yet. Create one or let the AI suggest themes from your tasks.
          </p>

          <div v-else class="space-y-2">
            <div
              v-for="theme in themeStore.themes"
              :key="theme.theme_id"
              class="flex items-center justify-between gap-3 p-3 rounded-lg border border-surface-light/20 dark:border-surface-dark/20"
            >
              <div class="flex items-center gap-3 min-w-0">
                <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: theme.color }" />
                <div class="min-w-0">
                  <p class="text-sm font-medium text-text-light dark:text-text-dark truncate">{{ theme.name }}</p>
                  <div class="flex flex-wrap gap-1 mt-1">
                    <span
                      v-for="day in theme.days"
                      :key="day"
                      class="px-1.5 py-0.5 text-[10px] rounded bg-surface-light/50 dark:bg-surface-dark/50 text-text-light/70 dark:text-text-dark/70"
                    >
                      {{ WEEKDAY_LABELS[day] }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-1 shrink-0">
                <UButton icon="i-lucide-pencil" variant="ghost" size="sm" color="neutral" @click="startEdit(theme)" />
                <UButton icon="i-lucide-trash" variant="ghost" size="sm" color="error" @click="handleDelete(theme)" />
              </div>
            </div>
          </div>
        </div>

        <!-- Create/Edit form -->
        <div v-if="showForm" class="p-4 rounded-lg bg-surface-light/50 dark:bg-surface-dark/50 border border-surface-light/20 dark:border-surface-dark/20 space-y-3">
          <h3 class="text-sm font-semibold text-text-light dark:text-text-dark">
            {{ editingTheme ? 'Edit Theme' : 'New Theme' }}
          </h3>

          <div class="flex flex-col">
            <label class="mb-1 text-xs font-medium text-text-light dark:text-text-dark">Name</label>
            <input
              v-model="form.name"
              type="text"
              placeholder="e.g. Deep Work"
              class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div class="flex flex-col">
            <label class="mb-1 text-xs font-medium text-text-light dark:text-text-dark">Description</label>
            <textarea
              v-model="form.description"
              rows="2"
              placeholder="Optional description"
              class="px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none"
            />
          </div>

          <div class="flex items-center gap-3">
            <div class="flex flex-col">
              <label class="mb-1 text-xs font-medium text-text-light dark:text-text-dark">Color</label>
              <input v-model="form.color" type="color" class="h-9 w-14 rounded border border-surface-light/30 dark:border-surface-dark/30" />
            </div>

            <div class="flex flex-col">
              <label class="mb-1 text-xs font-medium text-text-light dark:text-text-dark">Days</label>
              <div class="flex flex-wrap gap-1.5">
                <button
                  v-for="day in DAYS"
                  :key="day"
                  type="button"
                  class="px-2 py-1 text-xs rounded border"
                  :class="form.days.includes(day)
                    ? 'border-primary bg-primary/10 text-primary'
                    : 'border-surface-light/30 dark:border-surface-dark/30 text-text-light/70 dark:text-text-dark/70'"
                  @click="toggleDay(day)"
                >
                  {{ WEEKDAY_LABELS[day] }}
                </button>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2 pt-1">
            <UButton :label="editingTheme ? 'Save Changes' : 'Create Theme'" color="primary" size="sm" @click="submitForm" />
            <UButton label="Cancel" variant="ghost" color="neutral" size="sm" @click="resetForm" />
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>
