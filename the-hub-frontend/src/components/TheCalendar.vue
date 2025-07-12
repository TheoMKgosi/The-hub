<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useScheduleStore } from '@/stores/schedule'
import { useTaskStore } from '@/stores/tasks'
import { VueCal } from 'vue-cal'
import type { Event } from 'vue-cal'
import 'vue-cal/style'


const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const modalShow = ref(false)

async function fetchEvents() {
  if(scheduleStore.schedule === 0){
    await scheduleStore.fetchSchedule()
  }
  if(taskStore.tasks === 0) {
    await taskStore.fetchTasks()
  }
}

const formData = reactive({
  titleID: 0,
  title: '',
  start: null,
  end: null,
})

const toDateTimeLocal = (date: Date) =>
  new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)

const onCellClick = ({ cursor }) => {
  modalShow.value = true
  const clickedDate = new Date(cursor.date)
  formData.start = toDateTimeLocal(clickedDate)
  formData.end = toDateTimeLocal(new Date(clickedDate.getTime() + 60 * 60 * 1000))
}


function onViewChange(viewMeta: { start: Date; end: Date }) {
  fetchEvents(viewMeta.start, viewMeta.end)
}

async function onEventDropped(meta: { event: Event }) {
  console.log(meta.event)
  await useMyFetch(`calendar/${meta.event.task_id}`).patch({ start: meta.event.start }).json()
  fetchEvents()
}

async function onEventDelete(e) {
  scheduleStore.deleteSchedule(e.id)
}

function close() {
  modalShow.value = false
}

async function save() {
  const selectedTask = taskStore.tasks.find(t => t.task_id === formData.task_id)
  if (!selectedTask) {
    console.error('Task not found')
    return
  }

  const dataToSend = {
    task_id: formData.task_id,
    title: selectedTask.title,
    start: new Date(formData.start).toISOString(),
    end: new Date(formData.end).toISOString(),
  }

  await scheduleStore.submitForm(dataToSend)
  fetchEvents()
  close()
}

onMounted(() => {
  if (scheduleStore.schedule.length === 0) {
    scheduleStore.fetchSchedule()
  }
  if (taskStore.tasks.length === 0) {
    taskStore.fetchTasks()
  }
})
</script>

<template>
  <h2>Calendar</h2>
  <VueCal :events="scheduleStore.schedule" :editable-events="{ drag: false, resize: false, delete: true, create: true }" @view-change="onViewChange"
    @cell-click="onCellClick" @event-drop="onEventDropped" @event-delete="onEventDelete" />
  <div v-show="modalShow" class="fixed inset-0 bg-black/25 flex items-center justify-center">
    <div class="bg-white p-6 rounded shadow max-w-sm w-full">
      <form>
        <select v-model="formData.task_id" class="block">
          <option v-for="task in taskStore.tasks" :key="task.task_id" :value="task.task_id">
            {{ task.title }}
          </option>
        </select>

        <label for="start">Start Date:</label>
        <input type="datetime-local" id="start" v-model="formData.start" />

        <label for="end">End Date:</label>
        <input type="datetime-local" id="end" v-model="formData.end" />

      </form>
      <button @click="save" class="mt-4 p-2 bg-green-700 text-white">Save</button>
      <button @click="close" class="mt-4 ml-2 p-2  bg-red-700 text-white">Close</button>
    </div>
  </div>
</template>
