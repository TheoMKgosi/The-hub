<script setup lang="ts">
const props = defineProps<{
  show: boolean
  preSelectedTaskId?: string
}>()

const emit = defineEmits<{
  close: []
}>()

const focusStore = useFocusSessionStore()
const taskStore = useTaskStore()
const toast = useToast()

const mode = ref<'pomodoro' | 'timer'>('pomodoro')
const selectedTaskId = ref<string | undefined>(props.preSelectedTaskId)
const notes = ref('')
const timerStatus = ref<'idle' | 'running' | 'paused'>('idle')

const workDuration = ref(25)
const breakDuration = ref(5)
const cycleCount = ref(1)
const isBreak = ref(false)

const elapsedSeconds = ref(0)
const remainingSeconds = ref(0)
const totalPlannedSeconds = ref(0)
let timerInterval: ReturnType<typeof setInterval> | null = null
let activeSessionId: string | null = null
let cycleGroupId: string | null = null

const ringRadius = 90
const ringCircumference = 2 * Math.PI * ringRadius

const displayTime = computed(() => {
  const totalSecs = mode.value === 'pomodoro' ? remainingSeconds.value : elapsedSeconds.value
  const mins = Math.floor(totalSecs / 60)
  const secs = totalSecs % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
})

const displayLabel = computed(() => {
  if (timerStatus.value === 'idle') return mode.value === 'pomodoro' ? 'Ready' : 'Ready'
  if (isBreak.value) return 'Break'
  if (mode.value === 'pomodoro') return `Cycle ${cycleCount} / 4`
  return 'Focusing'
})

const ringProgress = computed(() => {
  if (timerStatus.value === 'idle') return ringCircumference
  if (mode.value === 'pomodoro') {
    const fraction = remainingSeconds.value / totalPlannedSeconds.value
    return ringCircumference * fraction
  }
  // Simple timer: show fill relative to a 30-min reference
  const ref = Math.max(elapsedSeconds.value, 1800)
  return ringCircumference * (1 - elapsedSeconds.value / ref)
})

const ringColor = computed(() => {
  if (timerStatus.value === 'paused') return '#f59e0b'
  if (timerStatus.value === 'running') return 'var(--color-primary)'
  if (isBreak.value) return '#f59e0b'
  return 'oklch(0.5 0.02 240)'
})

const ringTrackColor = computed(() => {
  if (isBreak.value) return 'oklch(0.5 0.1 85 / 0.15)'
  return 'oklch(0.5 0.02 240 / 0.15)'
})

const isActive = computed(() => timerStatus.value === 'running' || timerStatus.value === 'paused')
const canStart = computed(() => timerStatus.value === 'idle')

const tasks = computed(() => taskStore.tasks.filter(t => t.status !== 'completed'))

const todayMinutes = computed(() => focusStore.getTodayMinutes)

watch(() => props.show, (val) => {
  if (val) {
    resetTimer()
    if (props.preSelectedTaskId) {
      selectedTaskId.value = props.preSelectedTaskId
    }
    focusStore.fetchStats(7)
  }
})

onUnmounted(() => {
  clearTimerInterval()
})

function clearTimerInterval() {
  if (timerInterval) {
    clearInterval(timerInterval)
    timerInterval = null
  }
}

function resetTimer() {
  clearTimerInterval()
  timerStatus.value = 'idle'
  elapsedSeconds.value = 0
  remainingSeconds.value = workDuration.value * 60
  totalPlannedSeconds.value = 0
  isBreak.value = false
  cycleCount.value = 1
  activeSessionId = null
  cycleGroupId = null
}

function generateCycleGroup(): string {
  return crypto.randomUUID?.() || `${Date.now()}-${Math.random().toString(36).slice(2)}`
}

async function startTimer() {
  if (mode.value === 'pomodoro') {
    totalPlannedSeconds.value = workDuration.value * 60
    remainingSeconds.value = workDuration.value * 60
  }

  if (!cycleGroupId && mode.value === 'pomodoro') {
    cycleGroupId = generateCycleGroup()
  }

  timerStatus.value = 'running'

  const sessionType = mode.value === 'pomodoro' ? (isBreak.value ? 'break' : 'focus') : 'focus'

  const session = await focusStore.startSession({
    task_id: selectedTaskId.value || undefined,
    session_type: sessionType,
    planned_duration: mode.value === 'pomodoro'
      ? (isBreak.value ? breakDuration.value : workDuration.value)
      : undefined,
    notes: notes.value,
    cycle_group: cycleGroupId || undefined,
    cycle_order: isBreak.value ? cycleCount.value : cycleCount.value
  })

  if (session) {
    activeSessionId = session.focus_session_id
    startCountdown()
  }
}

function startCountdown() {
  clearTimerInterval()
  timerInterval = setInterval(() => {
    if (timerStatus.value !== 'running') return

    if (mode.value === 'timer') {
      elapsedSeconds.value++
    } else {
      remainingSeconds.value--
      if (remainingSeconds.value <= 0) {
        completePomodoroCycle()
      }
    }
  }, 1000)
}

async function completePomodoroCycle() {
  if (activeSessionId) {
    await focusStore.stopSession(activeSessionId, 'completed')
    activeSessionId = null
  }

  if (!isBreak.value) {
    isBreak.value = true
    remainingSeconds.value = breakDuration.value * 60
    toast.add({title: "Task", description: `Work cycle ${cycleCount.value} complete! Starting break.`})
    await startTimer()
  } else {
    if (cycleCount.value >= 4) {
      toast.add({title: "Task", description: "Pomodoro session complete! Great focus!", color: "success"})
      resetTimer()
      return
    }
    cycleCount.value++
    isBreak.value = false
    remainingSeconds.value = workDuration.value * 60
    toast.add({title: "Task", description: `Break over! Starting cycle ${cycleCount.value}.`})
    await startTimer()
  }
}

async function pauseTimer() {
  if (timerStatus.value === 'running') {
    timerStatus.value = 'paused'
  } else if (timerStatus.value === 'paused') {
    timerStatus.value = 'running'
  }
}

async function stopTimer() {
  if (activeSessionId) {
    await focusStore.stopSession(activeSessionId, 'interrupted')
    activeSessionId = null
  }
  resetTimer()
  focusStore.fetchStats(7)
}

async function finishTimer() {
  if (activeSessionId) {
    await focusStore.stopSession(activeSessionId, 'completed')
    activeSessionId = null
  }

  const durationStr = mode.value === 'pomodoro'
    ? `${Math.round(totalPlannedSeconds.value / 60)} min`
    : `${Math.floor(elapsedSeconds.value / 60)} min ${elapsedSeconds.value % 60} sec`
  toast.add({title: "Task", description: `Session complete! Duration: ${durationStr}`, color: "success"})

  resetTimer()
  focusStore.fetchSessions({ limit: 5 })
  focusStore.fetchStats(7)
}

function handleClose() {
  if (isActive.value) {
    toast.add({title: "Task", description: "Session stopped - modal closed", color: "warning"})
    stopTimer()
  }
  resetTimer()
  emit('close')
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 bg-black/60 dark:bg-black/80 flex items-center justify-center z-50 p-4"
    @click.self="handleClose">
    <div
      class="bg-surface-light dark:bg-surface-dark rounded-2xl shadow-2xl max-w-md w-full max-h-[90vh] overflow-y-auto custom-scrollbar">

      <!-- Header -->
      <div class="p-5 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-xl bg-primary/10 dark:bg-primary/20 flex items-center justify-center">
            <UIcon name="i-lucide-timer" class="w-5 h-5 text-primary" />
          </div>
          <div>
            <h2 class="text-base font-semibold text-text-light dark:text-text-dark leading-none">Focus Session</h2>
            <p class="text-[11px] text-text-light/50 dark:text-text-dark/50 mt-0.5">Today: {{ Math.floor(todayMinutes /
              60) }}h {{ todayMinutes % 60 }}m</p>
          </div>
        </div>
        <button @click="handleClose"
          class="p-2 rounded-lg hover:bg-surface-light/10 dark:hover:bg-surface-dark/10 text-text-light/40 dark:text-text-dark/40 hover:text-text-light dark:hover:text-text-dark transition-colors">
          <UIcon name="i-lucide-x" class="w-4 h-4" />
        </button>
      </div>

      <!-- Mode toggle -->
      <div class="px-5 pb-4">
        <div class="flex bg-surface-light/10 dark:bg-surface-dark/10 rounded-xl p-1">
          <button :class="[
            'flex-1 flex items-center justify-center gap-1.5 py-2.5 text-sm font-medium rounded-lg transition-all duration-200',
            mode === 'pomodoro'
              ? 'bg-primary text-white shadow-sm shadow-primary/20'
              : 'text-text-light/50 dark:text-text-dark/50 hover:text-text-light dark:hover:text-text-dark'
          ]" @click="mode = 'pomodoro'">
            <UIcon name="i-lucide-clock" class="w-4 h-4" />
            Pomodoro
          </button>
          <button :class="[
            'flex-1 flex items-center justify-center gap-1.5 py-2.5 text-sm font-medium rounded-lg transition-all duration-200',
            mode === 'timer'
              ? 'bg-primary text-white shadow-sm shadow-primary/20'
              : 'text-text-light/50 dark:text-text-dark/50 hover:text-text-light dark:hover:text-text-dark'
          ]" @click="mode = 'timer'">
            <UIcon name="i-lucide-hourglass" class="w-4 h-4" />
            Timer
          </button>
        </div>
      </div>

      <!-- Circular progress ring -->
      <div class="px-5 pb-4 flex justify-center">
        <div class="relative w-56 h-56">
          <svg class="w-56 h-56 -rotate-90" viewBox="0 0 200 200">
            <circle cx="100" cy="100" :r="ringRadius" fill="none" :stroke="ringTrackColor" stroke-width="8" />
            <circle cx="100" cy="100" :r="ringRadius" fill="none" :stroke="ringColor" stroke-width="8"
              stroke-linecap="round" :stroke-dasharray="ringCircumference" :stroke-dashoffset="ringProgress"
              class="transition-[stroke-dashoffset] duration-700 ease-out"
              :class="{ 'drop-shadow-[0_0_8px_var(--color-primary)]': timerStatus === 'running' }" />
          </svg>
          <div class="absolute inset-0 flex flex-col items-center justify-center">
            <span class="text-4xl font-mono font-bold text-text-light dark:text-text-dark tracking-wider tabular-nums"
              :class="{ 'text-primary': timerStatus === 'running', 'text-warning': timerStatus === 'paused' }">
              {{ displayTime }}
            </span>
            <span class="text-xs font-medium mt-1.5 px-2.5 py-0.5 rounded-full" :class="{
              'bg-primary/10 dark:bg-primary/20 text-primary': timerStatus === 'running',
              'bg-amber-500/10 dark:bg-amber-500/20 text-amber-500': timerStatus === 'paused',
              'bg-surface-light/10 dark:bg-surface-dark/10 text-text-light/40 dark:text-text-dark/40': timerStatus === 'idle',
            }">
              {{ timerStatus === 'paused' ? 'PAUSED' : displayLabel }}
            </span>
          </div>
        </div>
      </div>

      <!-- Pomodoro settings (only when idle) -->
      <div v-if="mode === 'pomodoro' && !isActive" class="px-5 pb-4">
        <div class="flex gap-3 justify-center">
          <div class="flex flex-col items-center gap-1.5">
            <label class="text-[11px] text-text-light/50 dark:text-text-dark/50 font-medium">Work</label>
            <div class="flex items-center gap-1">
              <button @click="workDuration = Math.max(5, workDuration - 5)"
                class="w-7 h-7 rounded-lg bg-surface-light/10 dark:bg-surface-dark/10 flex items-center justify-center text-text-light/60 dark:text-text-dark/60 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 hover:text-text-light dark:hover:text-text-dark transition-colors text-sm">
                -
              </button>
              <span class="w-8 text-center text-sm font-semibold text-text-light dark:text-text-dark tabular-nums">{{
                workDuration }}</span>
              <button @click="workDuration = Math.min(120, workDuration + 5)"
                class="w-7 h-7 rounded-lg bg-surface-light/10 dark:bg-surface-dark/10 flex items-center justify-center text-text-light/60 dark:text-text-dark/60 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 hover:text-text-light dark:hover:text-text-dark transition-colors text-sm">
                +
              </button>
            </div>
            <span class="text-[10px] text-text-light/40 dark:text-text-dark/40">min</span>
          </div>
          <div class="w-px bg-surface-light/20 dark:bg-surface-dark/20" />
          <div class="flex flex-col items-center gap-1.5">
            <label class="text-[11px] text-text-light/50 dark:text-text-dark/50 font-medium">Break</label>
            <div class="flex items-center gap-1">
              <button @click="breakDuration = Math.max(1, breakDuration - 1)"
                class="w-7 h-7 rounded-lg bg-surface-light/10 dark:bg-surface-dark/10 flex items-center justify-center text-text-light/60 dark:text-text-dark/60 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 hover:text-text-light dark:hover:text-text-dark transition-colors text-sm">
                -
              </button>
              <span class="w-8 text-center text-sm font-semibold text-text-light dark:text-text-dark tabular-nums">{{
                breakDuration }}</span>
              <button @click="breakDuration = Math.min(30, breakDuration + 1)"
                class="w-7 h-7 rounded-lg bg-surface-light/10 dark:bg-surface-dark/10 flex items-center justify-center text-text-light/60 dark:text-text-dark/60 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 hover:text-text-light dark:hover:text-text-dark transition-colors text-sm">
                +
              </button>
            </div>
            <span class="text-[10px] text-text-light/40 dark:text-text-dark/40">min</span>
          </div>
        </div>
      </div>

      <!-- Task selector -->
      <div class="px-5 pb-4">
        <label class="block text-[11px] font-medium text-text-light/50 dark:text-text-dark/50 mb-1.5">Task
          (optional)</label>
        <div class="relative">
          <select v-model="selectedTaskId"
            class="w-full px-3 py-2.5 pr-8 bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 text-text-light dark:text-text-dark rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 appearance-none cursor-pointer transition-colors">
            <option :value="undefined">No task — just focus</option>
            <option v-for="task in tasks" :key="task.task_id" :value="task.task_id">
              {{ task.title }}
            </option>
          </select>
          <UIcon name="i-lucide-chevron-down"
            class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-light/40 dark:text-text-dark/40 pointer-events-none" />
        </div>
      </div>

      <!-- Notes -->
      <div class="px-5 pb-4">
        <label class="block text-[11px] font-medium text-text-light/50 dark:text-text-dark/50 mb-1.5">Notes</label>
        <textarea v-model="notes" placeholder="What are you focusing on?"
          class="w-full px-3 py-2.5 bg-surface-light/10 dark:bg-surface-dark/10 border border-surface-light/20 dark:border-surface-dark/20 text-text-light dark:text-text-dark placeholder:text-text-light/30 dark:placeholder:text-text-dark/30 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 resize-none transition-colors"
          rows="2"></textarea>
      </div>

      <!-- Controls -->
      <div class="px-5 pb-5 flex justify-center gap-2.5">
        <UButton v-if="canStart" label="Start" icon="i-lucide-play" color="primary" size="lg" @click="startTimer" />

        <UButton v-if="timerStatus === 'running'" label="Pause" icon="i-lucide-pause" color="warning" size="lg"
          @click="pauseTimer" />

        <UButton v-if="timerStatus === 'paused'" label="Resume" icon="i-lucide-play" color="primary" size="lg"
          @click="pauseTimer" />

        <UButton v-if="isActive" label="Complete" icon="i-lucide-check" color="success" size="lg"
          @click="finishTimer" />

        <UButton v-if="isActive" label="Stop" icon="i-lucide-stop-circle" color="error" variant="outline" size="lg"
          @click="stopTimer" />
      </div>

      <!-- Recent sessions -->
      <div v-if="focusStore.sessions.length > 0" class="px-5 pb-5">
        <div class="flex items-center gap-1.5 mb-2.5">
          <UIcon name="i-lucide-history" class="w-3.5 h-3.5 text-text-light/40 dark:text-text-dark/40" />
          <h4 class="text-xs font-medium text-text-light/50 dark:text-text-dark/50 uppercase tracking-wider">Recent
            Sessions</h4>
        </div>
        <div class="space-y-1 max-h-32 overflow-y-auto custom-scrollbar">
          <div v-for="session in focusStore.sessions.slice(0, 5)" :key="session.focus_session_id"
            class="flex items-center justify-between text-xs py-2 px-2.5 rounded-lg bg-surface-light/10 dark:bg-surface-dark/10 hover:bg-surface-light/20 dark:hover:bg-surface-dark/20 transition-colors">
            <div class="flex items-center gap-2 min-w-0">
              <div :class="[
                'w-1.5 h-1.5 rounded-full shrink-0',
                session.session_type === 'break' ? 'bg-amber-400' : 'bg-primary'
              ]" />
              <span class="truncate text-text-light dark:text-text-dark">{{ session.notes || (session.session_type ===
                'break' ? 'Break' : 'Focus') }}</span>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <span class="text-text-light/50 dark:text-text-dark/50">{{ session.duration_min }}m</span>
              <span :class="[
                'px-1.5 py-0.5 rounded-full text-[10px] font-medium',
                session.status === 'completed' ? 'bg-success/10 text-success' : 'bg-red-500/10 text-red-400'
              ]">{{ session.status }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
