<script setup lang="ts">
const taskStore = useTaskStore()
const focusStore = useFocusSessionStore()

const showFocusModal = ref(false)
const showFocusPanel = ref(false)
const focusTaskId = ref<string | undefined>(undefined)

const todayMinutes = computed(() => focusStore.getTodayMinutes)

onMounted(() => {
  if (taskStore.tasks.length === 0) taskStore.fetchTasks()
  focusStore.fetchStats(7)
  focusStore.fetchSessions({ limit: 5, session_type: 'focus' })
})

function openFocus() {
  focusTaskId.value = undefined
  showFocusModal.value = true
}

function openFocusForTask(taskId: string) {
  focusTaskId.value = taskId
  showFocusModal.value = true
}

function toggleFocusPanel() {
  showFocusPanel.value = !showFocusPanel.value
}
</script>

<template>
  <div class="relative">
    <!-- Page header -->
    <div class="px-4 md:px-6 pt-4 md:pt-6 pb-2">
      <div class="flex items-start justify-between gap-4">
        <div class="min-w-0">
          <h1 class="text-xl md:text-2xl font-bold ">Plan Your Day</h1>
          <p class="text-sm text-text-light/50 dark:text-text-dark/50 mt-0.5">Organize tasks and stay focused</p>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <UButton
            :label="'Focus'" icon="i-lucide-timer"
            color="primary" @click="openFocus"
            class="shadow-sm shadow-primary/20" />

          <UButton
            :icon="showFocusPanel ? 'i-lucide-chevron-right' : 'i-lucide-chevron-left'"
            variant="ghost"
            color="neutral"
            @click="toggleFocusPanel"
            :label="showFocusPanel ? '' : undefined"
            class="hidden md:inline-flex"
            :ui="{ trailingIcon: 'hidden' }">
            <template v-if="todayMinutes > 0" #trailing>
              <span class="ml-1.5 text-xs text-primary font-semibold">
                {{ Math.floor(todayMinutes / 60) }}h {{ todayMinutes % 60 }}m
              </span>
            </template>
          </UButton>
        </div>
      </div>
    </div>

    <!-- Main content -->
    <div class="px-2 md:px-4 pb-4 md:pb-6">
      <TaskTemplate :taskList="taskStore.tasks" />
    </div>

    <!-- Focus Panel (collapsible) -->
    <FocusPanel :show="showFocusPanel" @close="showFocusPanel = false" @startFocus="openFocus" />

    <!-- Focus Session Modal -->
    <FocusSessionModal :show="showFocusModal" :preSelectedTaskId="focusTaskId" @close="showFocusModal = false" />
  </div>
</template>
