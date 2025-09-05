<script setup lang="ts">
const studySessionStore = useStudySessionStore()
const topicStore = useTopicStore()
const selectedPeriod = ref<'7' | '30' | '90'>('30')

const fetchAnalytics = async () => {
  await studySessionStore.fetchStats(parseInt(selectedPeriod.value))
  await topicStore.fetchTopics()
}

const totalStudyTime = computed(() => studySessionStore.getTotalStudyTime(parseInt(selectedPeriod.value)))
const averageDailyTime = computed(() => studySessionStore.getAverageDailyTime(parseInt(selectedPeriod.value)))
const topicBreakdown = computed(() => studySessionStore.getTopicBreakdown())
const dailyStats = computed(() => studySessionStore.getDailyStats())

const formatTime = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}h ${mins}m`
}

const getTopicTitle = (topicId: string) => {
  const topic = topicStore.topics.find(t => t.topic_id === topicId)
  return topic?.title || 'Unknown Topic'
}

const getStudyStreak = () => {
  const stats = dailyStats.value
  if (stats.length === 0) return 0

  let streak = 0
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  for (let i = 0; i < stats.length; i++) {
    const statDate = new Date(stats[i].date)
    statDate.setHours(0, 0, 0, 0)

    if (statDate.getTime() === today.getTime() - (i * 24 * 60 * 60 * 1000) && stats[i].minutes > 0) {
      streak++
    } else {
      break
    }
  }

  return streak
}

const getMostProductiveDay = () => {
  if (dailyStats.value.length === 0) return null

  const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
  const dayStats = dailyStats.value.reduce((acc, stat) => {
    const day = new Date(stat.date).getDay()
    acc[day] = (acc[day] || 0) + stat.minutes
    return acc
  }, {} as Record<number, number>)

  const mostProductiveDay = Object.entries(dayStats).reduce((max, [day, minutes]) =>
    minutes > (dayStats[parseInt(max[0])] || 0) ? [day, minutes] : max
  , ['0', 0])

  return {
    day: dayNames[parseInt(mostProductiveDay[0])],
    minutes: mostProductiveDay[1] as number
  }
}

const getLearningEfficiency = () => {
  const completedTopics = topicStore.topics.filter(t => t.status === 'completed').length
  const totalTopics = topicStore.topics.length
  const completionRate = totalTopics > 0 ? (completedTopics / totalTopics) * 100 : 0

  if (totalStudyTime.value === 0) return 'Start studying to see efficiency!'

  const efficiency = completionRate / (totalStudyTime.value / 60) // completion % per hour
  if (efficiency > 2) return 'Excellent efficiency! 🎉'
  if (efficiency > 1) return 'Good efficiency! 👍'
  if (efficiency > 0.5) return 'Moderate efficiency 📚'
  return 'Keep studying to improve efficiency 📖'
}

onMounted(() => {
  fetchAnalytics()
})

watch(selectedPeriod, () => {
  fetchAnalytics()
})
</script>

<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-text-light dark:text-text-dark">Learning Analytics</h1>
        <p class="text-text-light dark:text-text-dark/70 mt-1">Insights into your learning patterns and progress</p>
      </div>

      <!-- Period Selector -->
      <div class="flex gap-2">
        <button
          v-for="period in [
            { value: '7', label: '7 Days' },
            { value: '30', label: '30 Days' },
            { value: '90', label: '90 Days' }
          ]"
          :key="period.value"
          @click="selectedPeriod = period.value as any"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
            selectedPeriod === period.value
              ? 'bg-primary text-white'
              : 'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark hover:bg-surface-light/80 dark:hover:bg-surface-dark/80'
          ]"
        >
          {{ period.label }}
        </button>
      </div>
    </div>

    <!-- Key Metrics -->
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4 mb-8">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <div class="flex items-center gap-3">
          <div class="p-3 bg-primary/10 dark:bg-primary/20 rounded-lg">
            <svg class="w-6 h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p class="text-2xl font-bold text-text-light dark:text-text-dark">{{ formatTime(totalStudyTime) }}</p>
            <p class="text-sm text-text-light dark:text-text-dark/70">Total Study Time</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <div class="flex items-center gap-3">
          <div class="p-3 bg-success/10 dark:bg-success/20 rounded-lg">
            <svg class="w-6 h-6 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <p class="text-2xl font-bold text-text-light dark:text-text-dark">{{ formatTime(averageDailyTime) }}</p>
            <p class="text-sm text-text-light dark:text-text-dark/70">Daily Average</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <div class="flex items-center gap-3">
          <div class="p-3 bg-warning/10 dark:bg-warning/20 rounded-lg">
            <svg class="w-6 h-6 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div>
            <p class="text-2xl font-bold text-text-light dark:text-text-dark">{{ getStudyStreak() }}</p>
            <p class="text-sm text-text-light dark:text-text-dark/70">Day Streak</p>
          </div>
        </div>
      </div>

      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <div class="flex items-center gap-3">
          <div class="p-3 bg-secondary/10 dark:bg-secondary/20 rounded-lg">
            <svg class="w-6 h-6 text-secondary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div>
            <p class="text-lg font-bold text-text-light dark:text-text-dark">
              {{ getMostProductiveDay()?.day || 'N/A' }}
            </p>
            <p class="text-sm text-text-light dark:text-text-dark/70">
              {{ getMostProductiveDay() ? formatTime(getMostProductiveDay()!.minutes) : 'Most Productive' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Learning Efficiency -->
    <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark mb-8">
      <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4">Learning Efficiency</h2>
      <div class="text-center">
        <p class="text-lg text-text-light dark:text-text-dark mb-2">{{ getLearningEfficiency() }}</p>
        <div class="flex items-center justify-center gap-4 text-sm text-text-light dark:text-text-dark/70">
          <span>{{ topicStore.topics.filter(t => t.status === 'completed').length }} completed</span>
          <span>•</span>
          <span>{{ topicStore.topics.length }} total topics</span>
          <span>•</span>
          <span>{{ formatTime(totalStudyTime) }} studied</span>
        </div>
      </div>
    </div>

    <!-- Topic Breakdown -->
    <div class="grid gap-6 md:grid-cols-2">
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4">Study Time by Topic</h2>
        <div v-if="topicBreakdown.length > 0" class="space-y-3">
          <div
            v-for="topic in topicBreakdown.slice(0, 5)"
            :key="topic.topic_id || topic.topic_title"
            class="flex items-center justify-between"
          >
            <span class="text-sm text-text-light dark:text-text-dark truncate flex-1 mr-4">
              {{ getTopicTitle(topic.topic_id || '') || topic.topic_title }}
            </span>
            <div class="flex items-center gap-2">
              <div class="w-20 bg-surface-light dark:bg-surface-dark h-2 rounded-full overflow-hidden">
                <div
                  class="bg-primary h-full rounded-full"
                  :style="{ width: `${(topic.minutes / Math.max(...topicBreakdown.map(t => t.minutes))) * 100}%` }"
                ></div>
              </div>
              <span class="text-sm font-medium text-text-light dark:text-text-dark min-w-[60px] text-right">
                {{ formatTime(topic.minutes) }}
              </span>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-8 text-text-light dark:text-text-dark/60">
          <p>No study sessions recorded yet</p>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-surface-light dark:bg-surface-dark rounded-lg p-6 border border-surface-light dark:border-surface-dark">
        <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4">Recent Study Sessions</h2>
        <div v-if="studySessionStore.sessions.length > 0" class="space-y-3">
          <div
            v-for="session in studySessionStore.sessions.slice(0, 5)"
            :key="session.id"
            class="flex items-center justify-between py-2"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-primary/10 dark:bg-primary/20 rounded-full flex items-center justify-center">
                <svg class="w-4 h-4 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium text-text-light dark:text-text-dark">
                  {{ getTopicTitle(session.topic_id || '') || 'General Study' }}
                </p>
                <p class="text-xs text-text-light dark:text-text-dark/60">
                  {{ new Date(session.started_at).toLocaleDateString() }}
                </p>
              </div>
            </div>
            <span class="text-sm font-medium text-text-light dark:text-text-dark">
              {{ formatTime(session.duration_min) }}
            </span>
          </div>
        </div>
        <div v-else class="text-center py-8 text-text-light dark:text-text-dark/60">
          <p>No recent study sessions</p>
        </div>
      </div>
    </div>
  </div>
</template>