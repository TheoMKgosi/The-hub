<script setup lang="ts">
const focusStore = useFocusSessionStore()

onMounted(async () => {
  if (!focusStore.stats) {
    await focusStore.fetchStats(7)
  }
})

const dailyGoalMinutes = 120

const weeklyTotal = computed(() => {
  if (!focusStore.stats) return 0
  return focusStore.stats.weekly_stats?.reduce((sum, w) => sum + w.minutes, 0)
})

const todayMinutes = computed(() => focusStore.getTodayMinutes)

const todayProgressPercent = computed(() => {
  if (!todayMinutes.value) return 0
  return Math.min(100, Math.round((todayMinutes.value / dailyGoalMinutes) * 100))
})

const todayProgressFraction = computed(() => {
  if (dailyGoalMinutes <= 0) return 0
  return Math.min(1, todayMinutes.value / dailyGoalMinutes)
})

const circumference = 2 * Math.PI * 36

const maxDailyMinutes = computed(() => {
  if (!focusStore.stats?.daily_stats.length) return 60
  return Math.max(...focusStore.stats.daily_stats.map(d => d.minutes), 60)
})

const maxTaskMinutes = computed(() => {
  if (!focusStore.stats?.task_breakdown.length) return 1
  return Math.max(...focusStore.stats.task_breakdown.map(t => t.minutes), 1)
})
</script>

<template>
  <div v-if="focusStore.stats">
    <div class="space-y-4">
      <!-- Circular daily progress -->
      <div class="flex flex-col items-center">
        <div class="relative w-24 h-24">
          <svg class="w-24 h-24 -rotate-90" viewBox="0 0 80 80">
            <circle cx="40" cy="40" r="36" fill="none" stroke="currentColor"
              class="text-surface-light/20 dark:text-surface-dark/30" stroke-width="6" />
            <circle cx="40" cy="40" r="36" fill="none" :stroke="todayProgressPercent >= 100 ? 'url(#goalCompleteGradient)' : 'url(#goalGradient)'"
              stroke-width="6" stroke-linecap="round"
              :stroke-dasharray="circumference"
              :stroke-dashoffset="circumference - (todayProgressFraction * circumference)"
              class="transition-[stroke-dashoffset] duration-700 ease-out" />
            <defs>
              <linearGradient id="goalGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" class="text-primary" stop-color="currentColor" />
                <stop offset="100%" class="text-accent" stop-color="currentColor" />
              </linearGradient>
              <linearGradient id="goalCompleteGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" class="text-success" stop-color="currentColor" />
                <stop offset="100%" class="text-accent" stop-color="currentColor" />
              </linearGradient>
            </defs>
          </svg>
          <div class="absolute inset-0 flex flex-col items-center justify-center">
            <span class="text-lg font-bold text-text-light dark:text-text-dark leading-none">{{ Math.floor(todayMinutes / 60) }}h {{ todayMinutes % 60 }}m</span>
            <span class="text-[10px] text-text-light/50 dark:text-text-dark/50 mt-0.5">of {{ Math.floor(dailyGoalMinutes / 60) }}h goal</span>
          </div>
        </div>
        <span v-if="todayProgressPercent >= 100" class="mt-1.5 text-[11px] font-medium text-success flex items-center gap-1">
          <UIcon name="i-lucide-trophy" class="w-3 h-3" />
          Goal reached!
        </span>
        <span v-else-if="todayMinutes > 0" class="mt-1.5 text-[11px] font-medium text-primary">
          {{ todayProgressPercent }}% of daily goal
        </span>
      </div>

      <!-- Summary cards -->
      <div class="grid grid-cols-3 gap-2">
        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-xl p-3 text-center hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 transition-colors cursor-default">
          <UIcon name="i-lucide-sun" class="w-4 h-4 text-primary mx-auto mb-1" />
          <div class="text-base font-bold text-text-light dark:text-text-dark">{{ Math.floor(todayMinutes / 60) }}h {{ todayMinutes % 60 }}m</div>
          <div class="text-[10px] text-text-light/50 dark:text-text-dark/50 uppercase tracking-wider">Today</div>
        </div>
        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-xl p-3 text-center hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 transition-colors cursor-default">
          <UIcon name="i-lucide-calendar-days" class="w-4 h-4 text-secondary mx-auto mb-1" />
          <div class="text-base font-bold text-text-light dark:text-text-dark">{{ Math.floor(weeklyTotal / 60) }}h {{ weeklyTotal % 60 }}m</div>
          <div class="text-[10px] text-text-light/50 dark:text-text-dark/50 uppercase tracking-wider">Week</div>
        </div>
        <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-xl p-3 text-center hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 transition-colors cursor-default">
          <UIcon name="i-lucide-zap" class="w-4 h-4 text-accent mx-auto mb-1" />
          <div class="text-base font-bold text-text-light dark:text-text-dark">{{ focusStore.stats.total_sessions }}</div>
          <div class="text-[10px] text-text-light/50 dark:text-text-dark/50 uppercase tracking-wider">Sessions</div>
        </div>
      </div>

      <!-- Daily bar chart -->
      <div>
        <div class="flex items-center gap-1.5 mb-2.5">
          <UIcon name="i-lucide-chart-bar-big" class="w-3.5 h-3.5 text-text-light/40 dark:text-text-dark/40" />
          <h4 class="text-xs font-medium text-text-light/60 dark:text-text-dark/60">Daily Focus</h4>
        </div>
        <div class="flex items-end gap-1 h-20">
          <div v-for="day in focusStore.stats.daily_stats?.slice(-7)" :key="day.date"
            class="flex-1 flex flex-col items-center gap-1 h-full justify-end">
            <div class="w-full rounded-t-sm relative group cursor-default transition-all hover:opacity-80"
              :class="day.minutes > 0 ? 'bg-primary/60 dark:bg-primary/50' : 'bg-surface-light/20 dark:bg-surface-dark/20'"
              :style="{ height: `${Math.max((day.minutes / maxDailyMinutes) * 100, day.minutes > 0 ? 4 : 2)}%` }">
              <div class="absolute -top-7 left-1/2 -translate-x-1/2 bg-text-light dark:bg-text-dark text-surface-light dark:text-surface-dark text-[10px] font-medium px-1.5 py-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap shadow-lg">
                {{ day.minutes }}m
              </div>
            </div>
            <span class="text-[9px] text-text-light/40 dark:text-text-dark/40 font-medium">{{ new Date(day.date).toLocaleDateString('en-US', { weekday: 'narrow' }) }}</span>
          </div>
        </div>
      </div>

      <!-- Task breakdown -->
      <div v-if="focusStore.stats.task_breakdown?.length > 0">
        <div class="flex items-center gap-1.5 mb-2.5">
          <UIcon name="i-lucide-list-checks" class="w-3.5 h-3.5 text-text-light/40 dark:text-text-dark/40" />
          <h4 class="text-xs font-medium text-text-light/60 dark:text-text-dark/60">By Task</h4>
        </div>
        <div class="space-y-2">
          <div v-for="task in focusStore.stats.task_breakdown.slice(0, 5)" :key="task.task_id || 'no-task'"
            class="flex items-center gap-2.5">
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between text-xs mb-1">
                <span class="truncate text-text-light dark:text-text-dark font-medium">{{ task.task_title }}</span>
                <span class="text-text-light/50 dark:text-text-dark/50 shrink-0 ml-2 text-[11px]">{{ task.minutes }}m</span>
              </div>
              <div class="w-full h-1.5 bg-surface-light/20 dark:bg-surface-dark/20 rounded-full overflow-hidden">
                <div class="h-full bg-gradient-to-r from-primary to-accent rounded-full transition-all duration-500"
                  :style="{ width: `${(task.minutes / maxTaskMinutes) * 100}%` }" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No data -->
      <div v-if="focusStore.stats.total_minutes === 0" class="text-center py-6">
        <div class="w-12 h-12 rounded-full bg-text-light/5 dark:bg-text-dark/5 flex items-center justify-center mx-auto mb-3">
          <UIcon name="i-lucide-timer-off" class="w-6 h-6 text-text-light/20 dark:text-text-dark/20" />
        </div>
        <p class="text-xs text-text-light/40 dark:text-text-dark/40">No sessions yet</p>
        <p class="text-[11px] text-text-light/30 dark:text-text-dark/30 mt-1">Start a focus session to track your time</p>
      </div>
    </div>
  </div>

  <!-- Loading skeleton -->
  <div v-else class="space-y-4">
    <div class="animate-pulse flex flex-col items-center">
      <div class="w-24 h-24 bg-surface-light/20 dark:bg-surface-dark/20 rounded-full" />
    </div>
    <div class="grid grid-cols-3 gap-2">
      <div class="h-16 bg-surface-light/20 dark:bg-surface-dark/20 rounded-xl" />
      <div class="h-16 bg-surface-light/20 dark:bg-surface-dark/20 rounded-xl" />
      <div class="h-16 bg-surface-light/20 dark:bg-surface-dark/20 rounded-xl" />
    </div>
    <div class="h-20 bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg" />
    <div class="h-12 bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg" />
  </div>
</template>
