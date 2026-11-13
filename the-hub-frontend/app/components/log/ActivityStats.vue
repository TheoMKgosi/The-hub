<script setup lang="ts">
const activityStore = useActivityLogStore()

const rangeOptions = [
  { label: 'Last 7 days', value: 7 },
  { label: 'Last 30 days', value: 30 },
  { label: 'Last 90 days', value: 90 },
  { label: 'Last year', value: 365 }
]
const selectedRange = ref(30)

watch(selectedRange, (days) => {
  activityStore.fetchStats(days)
})

const stats = computed(() => activityStore.stats)

const summaryCards = computed(() => {
  if (!stats.value) return []
  const avgDaily = stats.value.days > 0 ? stats.value.total_minutes / stats.value.days : 0
  return [
    { label: 'Entries', value: stats.value.total_entries.toString(), icon: 'i-lucide-list' },
    {
      label: 'Total time',
      value: stats.value.total_minutes >= 60
        ? `${Math.floor(stats.value.total_minutes / 60)}h ${stats.value.total_minutes % 60}m`
        : `${stats.value.total_minutes}m`,
      icon: 'i-lucide-clock'
    },
    { label: 'Avg per day', value: `${Math.round(avgDaily)}m`, icon: 'i-lucide-trending-up' },
    { label: 'Current streak', value: `${stats.value.current_streak} ${stats.value.current_streak === 1 ? 'day' : 'days'}`, icon: 'i-lucide-flame' }
  ]
})

const maxCategoryMinutes = computed(() =>
  Math.max(1, ...(stats.value?.categories || []).map(c => c.minutes))
)

const dailyChart = computed(() => {
  const daily = stats.value?.daily_stats || []
  if (!daily.length) return []
  const max = Math.max(1, ...daily.map(d => d.minutes))
  return daily.slice(-30).map(d => ({
    label: d.date.slice(5),
    minutes: d.minutes,
    height: Math.round((d.minutes / max) * 100)
  }))
})

const weeklyChart = computed(() => {
  const weekly = stats.value?.weekly_stats || []
  if (!weekly.length) return []
  const max = Math.max(1, ...weekly.map(w => w.minutes))
  return weekly.slice(-12).map(w => ({
    label: w.week_start.slice(5),
    minutes: w.minutes,
    height: Math.round((w.minutes / max) * 100)
  }))
})

const formatMinutes = (minutes: number) => {
  if (!minutes) return '0m'
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  if (hours && mins) return `${hours}h ${mins}m`
  if (hours) return `${hours}h`
  return `${mins}m`
}

const formatDay = (date: string) =>
  new Date(`${date}T00:00:00`).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })

onMounted(() => {
  activityStore.fetchStats(selectedRange.value)
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold">Analysis</h2>
      <USelect v-model="selectedRange" :items="rangeOptions" class="w-40" />
    </div>

    <div v-if="!stats" class="py-8 text-center text-dimmed">Loading stats…</div>

    <template v-else>
      <!-- Summary cards -->
      <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
        <div
          v-for="card in summaryCards"
          :key="card.label"
          class="rounded-lg border border-surface-light dark:border-surface-dark bg-surface-light/50 dark:bg-surface-dark/50 p-4"
        >
          <div class="flex items-center gap-2 text-sm text-dimmed">
            <UIcon :name="card.icon" class="size-4" />
            {{ card.label }}
          </div>
          <p class="mt-1 text-2xl font-bold">{{ card.value }}</p>
        </div>
      </div>

      <!-- Category breakdown -->
      <div v-if="stats.categories!.length" class="rounded-lg border border-surface-light dark:border-surface-dark p-4">
        <h3 class="mb-3 font-semibold">Time by category</h3>
        <div class="space-y-2">
          <div v-for="category in stats.categories" :key="category.category">
            <div class="mb-1 flex justify-between text-sm">
              <span>{{ category.category }}</span>
              <span class="text-dimmed">{{ formatMinutes(category.minutes) }} · {{ category.entries }} {{ category.entries === 1 ? 'entry' : 'entries' }}</span>
            </div>
            <div class="h-2 w-full rounded-full bg-surface-light dark:bg-surface-dark overflow-hidden">
              <div
                class="h-full rounded-full bg-primary transition-all"
                :style="{ width: `${Math.max(2, Math.round((category.minutes / maxCategoryMinutes) * 100))}%` }"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Daily activity -->
      <div v-if="dailyChart.length" class="rounded-lg border border-surface-light dark:border-surface-dark p-4">
        <h3 class="mb-3 font-semibold">Daily activity</h3>
        <div class="flex items-end gap-1 h-40">
          <div
            v-for="day in dailyChart"
            :key="day.label"
            class="flex-1 flex flex-col items-center gap-1"
            :title="`${formatDay(day.label)}: ${formatMinutes(day.minutes)}`"
          >
            <div
              class="w-full rounded-t bg-primary/80 min-h-[2px]"
              :style="{ height: `${Math.max(2, day.height)}%` }"
            />
            <span class="text-[10px] text-dimmed">{{ day.label }}</span>
          </div>
        </div>
      </div>

      <!-- Weekly trends -->
      <div v-if="weeklyChart.length" class="rounded-lg border border-surface-light dark:border-surface-dark p-4">
        <h3 class="mb-3 font-semibold">Weekly trends</h3>
        <div class="flex items-end gap-2 h-32">
          <div
            v-for="week in weeklyChart"
            :key="week.label"
            class="flex-1 flex flex-col items-center gap-1"
            :title="`Week of ${formatDay(week.label)}: ${formatMinutes(week.minutes)}`"
          >
            <div
              class="w-full rounded-t bg-primary/60 min-h-[2px]"
              :style="{ height: `${Math.max(2, week.height)}%` }"
            />
            <span class="text-[10px] text-dimmed">{{ week.label }}</span>
          </div>
        </div>
      </div>

      <!-- Top tags -->
      <div v-if="stats.top_tags.length" class="rounded-lg border border-surface-light dark:border-surface-dark p-4">
        <h3 class="mb-3 font-semibold">Top tags</h3>
        <div class="flex flex-wrap gap-2">
          <UBadge v-for="tag in stats.top_tags" :key="tag.tag" variant="subtle">
            {{ tag.tag }} · {{ formatMinutes(tag.minutes) }} ({{ tag.count }})
          </UBadge>
        </div>
      </div>
    </template>
  </div>
</template>
