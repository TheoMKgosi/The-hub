<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const topicID = parseInt(route.params.id as string)

const topicStore = useTopicStore()
const taskLearningStore = useTaskLearningStore()

const showModal = reactive({ value: false })

const openModal = () => {
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  taskLearningStore.fetchTasks(topicID)
  topicStore.fetchTopic(topicID)
})

const formData = reactive({
  topic_id: topicID,
  title: '',
})

const addTask = () => {
  taskLearningStore.createTask(formData)
}

const dateTransform = (date: Date) => new Date(date).toLocaleDateString()

const isComplete = (status) => status === "finished"

const completeTask = (task) => {
  if (task.status === 'pending') {
    task.status = 'finished'
    taskLearningStore.completeTask(task)
  } else {
    task.status = 'pending'
    taskLearningStore.completeTask(task)
  }
}
</script>
<template>
  <div to="#learning" class="max-w-3xl mx-auto p-6 space-y-6">

    <!-- Header -->
    <header class="space-y-2">
      <!-- Go Back Button -->
      <button @click="goBack" class="mb-4 flex items-center gap-2 text-gray-600 hover:text-gray-800 transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
        </svg>
        Go Back
      </button>

      <h1 class="text-3xl font-bold">{{ topicStore.topic?.title }}</h1>
      <div class="flex flex-wrap gap-2 text-sm text-gray-600">
        <span class="bg-gray-100 px-2 py-1 rounded">🎯 goal</span>
        <span class="bg-gray-100 px-2 py-1 rounded">📅 Due {{ dateTransform(topicStore.topic?.deadline) }}</span>
        <!-- <span class="bg-gray-100 px-2 py-1 rounded">⏱️ 20h estimated</span> -->
      </div>
      <button @click="openModal" class="mt-2 px-4 py-2 bg-black text-white rounded hover:bg-gray-800">
        + Add Task
      </button>
    </header>

    <!-- Progress Overview -->
    <!--
    <section class="bg-white p-4 shadow rounded space-y-3">
      <h2 class="text-xl font-semibold">Progress Overview</h2>
      <p class="text-sm text-gray-500">1 of 2 tasks completed</p>
      <div class="w-full bg-gray-200 h-2 rounded">
        <div class="bg-black h-2 rounded" style="width: 50%;"></div>
      </div>
      <div class="flex gap-2 mt-2">
        <span class="bg-gray-200 text-sm px-2 py-1 rounded">Databases</span>
        <span class="bg-gray-200 text-sm px-2 py-1 rounded">Backend</span>
      </div>
      <p class="bg-gray-100 text-sm p-2 rounded">
        <strong>Motivation:</strong> Need this for my new job
      </p>
    </section>
      -->

    <!-- Tasks & Lessons -->
    <section class="space-y-4">
      <h2 class="text-xl font-semibold">Tasks & Lessons</h2>
      <p class="text-sm text-gray-500">Break down your learning into manageable tasks</p>

      <!-- Task 1 - Completed -->
      <template v-for="task in taskLearningStore.tasks" :key="task.task_learning_id">
        <div class="bg-white p-4 rounded shadow space-y-2">
          <div class="flex items-start gap-2">
            <input type="checkbox" :checked="isComplete(task.status)" @click="completeTask(task)" class="mt-1" />
            <div class="flex flex-col flex-grow">
              <span class="text-gray-900">{{ task.title }}</span>
              <!--
              <div class="flex gap-2 text-xs mt-1">
                <span class="bg-red-500 text-white px-2 py-0.5 rounded">high</span>
                <span class="text-gray-500">30min</span>
              </div>
              -->
            </div>
          </div>
        </div>
      </template>
    </section>
  </div>

  <teleport to="body">
    <div v-if="showModal.value" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded shadow-md w-full max-w-md">
        <h3 class="text-lg font-semibold mb-4">Add a New Task</h3>
        <form @submit.prevent="addTask">
          <input v-model="formData.title" placeholder="Title" class="w-full border p-2 mb-2 rounded" />
          <!--
          <textarea v-model="formData.notes" placeholder="Notes" class="w-full border p-2 mb-2 rounded"></textarea>
          <select v-model="formData.status" class="w-full border p-2 mb-2 rounded">
            <option value="pending">Pending</option>
            <option value="in_progress">In Progress</option>
            <option value="done">Done</option>
          </select>
          -->
          <!-- Add other fields as needed -->
          <div class="flex justify-end gap-2">
            <button type="button" @click="closeModal" class="px-4 py-2 bg-gray-200 rounded">Cancel</button>
            <button type="submit" class="px-4 py-2 bg-black text-white rounded">Save</button>
          </div>
        </form>
      </div>
    </div>
  </teleport>

</template>
