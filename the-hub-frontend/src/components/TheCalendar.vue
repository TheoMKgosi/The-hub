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
  title: '',
  start: new Date(),
  end: new Date(),
})

const toDateTimeLocal = (date: Date) =>
  new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16)

const onCellClick = ({ cursor }) => {
  modalShow.value = true
  const clickedDate = new Date(cursor.date)

  // Round down to hour
  clickedDate.setMinutes(0)
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

  const dataToSend = {
    title: formData.title,
    start: new Date(formData.start),
    end: new Date(formData.end),
  }

  console.log(dataToSend)

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
        <label for="title">Title</label>
        <input type="text" id="title" v-model="formData.title" autocomplete="off"/>

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
