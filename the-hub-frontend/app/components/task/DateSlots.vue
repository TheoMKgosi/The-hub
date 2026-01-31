<script setup lang="ts">
const times: number[] = []
for (let index = 0; index < 24; index++) {
  times.push(index)
}

interface Props {
  label: string
  tasks: { start_time: number, end_time: number, title: string }[]
}

interface Slot {
  time: number
  task: string | null
  position: 'end' | 'start' | 'both' | ''
  complete?: boolean | null
}


const props = defineProps<Props>()

const sloting = computed<Slot[]>(() => {
  return times.map(t => {
    const matchingTasks = props.tasks.find(task => {
      return (t >= task.start_time) && (t <= task.end_time)
    })

    if (!matchingTasks) {
      return { time: t, task: null, position: '' }
    }

    // 1. Check for 'both' first
    const isStart = matchingTasks.start_time === t
    const isEnd = matchingTasks.end_time === t

    let position: 'end' | 'start' | 'both' | '' = ''
    if (isStart && isEnd) position = 'both'
    else if (isStart) position = 'start'
    else if (isEnd) position = 'end'

    return {
      time: t,
      task: matchingTasks.title,
      position: position,
      complete: matchingTasks.title == '' ? null : false
    }
  })
})

</script>
<template>
  <div class="p-4 shadow rounded-2xl dark:inset-shadow-sm inset-shadow-gray-500/50 mr-2">
    <h2 class="font-bold p-2">{{ label }}</h2>
    <TimeSlot v-for="time in sloting" :time="time.time" :title="time.task" :position="time.position" :complete="time.complete"/>
  </div>
</template>
