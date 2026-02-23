<script setup lang="ts">
const times: number[] = []
for (let index = 0; index < 24; index++) {
  times.push(index)
}

interface Props {
  label: string
  date: Date
  tasks: { task_id: string, start_time: Date, end_time: Date, title: string }[]
}

interface Slot {
  time: number
  task: string | null
  task_id: string | null
  position: 'end' | 'start' | 'both' | ''
  complete?: boolean | null
}


const props = defineProps<Props>()

const sloting = computed<Slot[]>(() => {
  return times.map(t => {
    const matchingTasks = props.tasks.find(task => {
      const startHour = task.start_time?.getHours()
      const endHour = task.end_time?.getHours()
      
      // Handle tasks that span midnight (e.g., 23:00 to 01:00)
      if (startHour !== undefined && endHour !== undefined) {
        if (startHour < endHour) {
          // Normal case: task within same day (e.g., 10:00 to 12:00)
          return t >= startHour && t < endHour
        } else {
          // Task spans midnight (e.g., 23:00 to 01:00)
          return t >= startHour || t < endHour
        }
      }
      return false
    })

    if (!matchingTasks) {
      return { time: t, task: null, task_id: null, position: '', complete: null }
    }

    // 1. Check for 'both' first
    const isStart = matchingTasks.start_time.getHours() === t
    const isEnd = matchingTasks.end_time.getHours() === t

    let position: 'end' | 'start' | 'both' | '' = ''
    if (isStart && isEnd) position = 'both'
    else if (isStart) position = 'start'
    else if (isEnd) position = 'end'

    return {
      time: t,
      task: matchingTasks.title,
      task_id: matchingTasks.task_id,
      position: position,
      complete: matchingTasks.title == '' ? null : false
    }
  })
})

</script>
<template>
  <div class="p-4 shadow rounded-2xl dark:inset-shadow-sm inset-shadow-gray-500/50 mr-2">
    <h2 class="font-bold p-2">{{ label }}</h2>
    <TimeSlot v-for="time in sloting" :key="time.time" :time="time.time" :date="props.date" :title="time.task" :task-id="time.task_id" :position="time.position" :complete="time.complete"/>
  </div>
</template>
