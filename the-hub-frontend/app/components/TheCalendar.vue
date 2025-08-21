<script setup lang="ts">
import { useScheduleStore } from '@/stores/schedule'
import { useTaskStore } from '@/stores/tasks'
// import { VueCal } from 'vue-cal'
// import 'vue-cal/style'


const scheduleStore = useScheduleStore()
const taskStore = useTaskStore()
const modalShow = ref(false)

async function fetchEvents() {
  if (scheduleStore.schedule === 0) {
    await scheduleStore.fetchSchedule()
  }
  if (taskStore.tasks === 0) {
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
  <!--
  <VueCal :events="scheduleStore.schedule" :editable-events="{ drag: false, resize: false, delete: true, create: true }"
    @view-change="onViewChange" @cell-click="onCellClick" @event-drop="onEventDropped" @event-delete="onEventDelete" />
  <div v-show="modalShow" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center z-50">
    <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-lg max-w-sm w-full mx-4">
      <form class="space-y-4">
        <div>
          <label for="title" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Title</label>
          <input type="text" id="title" v-model="formData.title" autocomplete="off"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div>
          <label for="start" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Start Date:</label>
          <input type="datetime-local" id="start" v-model="formData.start"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

        <div>
          <label for="end" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">End Date:</label>
          <input type="datetime-local" id="end" v-model="formData.end"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-text-light dark:text-text-dark focus:outline-none focus:ring-2 focus:ring-primary" />
        </div>

      </form>
      <div class="flex justify-end space-x-2 mt-6">
        <UiButton @click="close" variant="default" size="md">Cancel</UiButton>
        <UiButton @click="save" variant="primary" size="md">Save</UiButton>
      </div>
    </div>
  </div>
-->
</template>
