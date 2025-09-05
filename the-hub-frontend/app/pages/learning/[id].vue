<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const topicID = route.params.id as string

const topicStore = useTopicStore()
const taskLearningStore = useTaskLearningStore()
const studySessionStore = useStudySessionStore()
const resourceStore = useResourceStore()

const showModal = reactive({ value: false })
const showEditModal = reactive({ value: false })
const editingTask = ref(null)
const showResourceModal = reactive({ value: false })
const resourceFormData = reactive({
  title: '',
  link: '',
  type: 'article' as 'video' | 'article' | 'document' | 'book' | 'course',
  notes: '',
  task_id: null as string | null
})
const studyTimer = reactive({
  isRunning: false,
  startTime: null as Date | null,
  elapsedSeconds: 0,
  interval: null as NodeJS.Timeout | null
})

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
  resourceStore.fetchResources({ topic_id: topicID })
})

const formData = reactive({
  title: '',
})

const editFormData = reactive({
  title: '',
  notes: '',
  status: 'not_started'
})

const addTask = () => {
  if (formData.title.trim()) {
    taskLearningStore.createTask(topicID, { title: formData.title })
    formData.title = ''
    closeModal()
  }
}

const dateTransform = (date: Date) => new Date(date).toLocaleDateString()

const isComplete = (status) => status === "completed"

const completeTask = (task) => {
  if (task.status === 'not_started' || task.status === 'on_hold') {
    task.status = 'completed'
  } else if (task.status === 'completed') {
    task.status = 'not_started'
  }
  taskLearningStore.completeTask(topicID, task)
}

const editTask = (task) => {
  editingTask.value = task
  editFormData.title = task.title
  editFormData.notes = task.notes || ''
  editFormData.status = task.status
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingTask.value = null
  editFormData.title = ''
  editFormData.notes = ''
  editFormData.status = 'not_started'
}

const updateTask = async () => {
  if (!editingTask.value || !editFormData.title.trim()) return

  try {
    await taskLearningStore.updateTask(topicID, editingTask.value.task_learning_id, {
      title: editFormData.title,
      notes: editFormData.notes,
      status: editFormData.status
    })
    closeEditModal()
  } catch (error) {
    console.error('Failed to update task:', error)
  }
}

// Resource management functions
const openResourceModal = (taskId: string | null = null) => {
  resourceFormData.task_id = taskId
  resourceFormData.title = ''
  resourceFormData.link = ''
  resourceFormData.type = 'article'
  resourceFormData.notes = ''
  showResourceModal.value = true
}

const closeResourceModal = () => {
  showResourceModal.value = false
  resourceFormData.task_id = null
}

const addResource = async () => {
  if (!resourceFormData.title.trim()) return

  try {
    await resourceStore.createResource({
      topic_id: resourceFormData.task_id ? undefined : topicID,
      task_id: resourceFormData.task_id || undefined,
      title: resourceFormData.title,
      link: resourceFormData.link,
      type: resourceFormData.type,
      notes: resourceFormData.notes
    })
    closeResourceModal()
  } catch (error) {
    console.error('Failed to add resource:', error)
  }
}

const deleteResource = async (resourceId: string) => {
  if (confirm('Are you sure you want to delete this resource?')) {
    try {
      await resourceStore.deleteResource(resourceId)
    } catch (error) {
      console.error('Failed to delete resource:', error)
    }
  }
}

// Study timer functions
const startStudyTimer = () => {
  if (!studyTimer.isRunning) {
    studyTimer.isRunning = true
    studyTimer.startTime = new Date()
    studyTimer.interval = setInterval(() => {
      studyTimer.elapsedSeconds++
    }, 1000)
  }
}

const pauseStudyTimer = () => {
  if (studyTimer.isRunning && studyTimer.interval) {
    studyTimer.isRunning = false
    clearInterval(studyTimer.interval)
    studyTimer.interval = null
  }
}

const stopStudyTimer = async () => {
  if (studyTimer.interval) {
    clearInterval(studyTimer.interval)
    studyTimer.interval = null
  }

  if (studyTimer.elapsedSeconds > 0) {
    const durationMin = Math.ceil(studyTimer.elapsedSeconds / 60)
    try {
      await studySessionStore.createSession({
        topic_id: topicID,
        duration_min: durationMin
      })
    } catch (error) {
      console.error('Failed to save study session:', error)
    }
  }

  studyTimer.isRunning = false
  studyTimer.startTime = null
  studyTimer.elapsedSeconds = 0
}

const formatTime = (seconds: number) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }
  return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// Cleanup timer on unmount
onUnmounted(() => {
  if (studyTimer.interval) {
    clearInterval(studyTimer.interval)
  }
})

const deleteTask = async (task) => {
  if (confirm('Are you sure you want to delete this task?')) {
    try {
      await taskLearningStore.deleteTask(topicID, task.task_learning_id)
    } catch (error) {
      console.error('Failed to delete task:', error)
    }
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
     <section class="bg-white dark:bg-surface-dark p-6 shadow rounded-lg space-y-4 border border-surface-light dark:border-surface-dark">
       <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Progress Overview</h2>
       <div class="flex items-center justify-between">
         <p class="text-sm text-text-light dark:text-text-dark/70">
           {{ taskLearningStore.tasks.filter(task => task.status === 'completed').length }} of {{ taskLearningStore.tasks.length }} tasks completed
         </p>
         <span class="text-sm font-medium px-3 py-1 rounded-full bg-primary/10 dark:bg-primary/20 text-primary">
           {{ Math.round((taskLearningStore.tasks.filter(task => task.status === 'completed').length / Math.max(taskLearningStore.tasks.length, 1)) * 100) }}%
         </span>
       </div>
       <div class="w-full bg-surface-light dark:bg-surface-dark h-3 rounded-full overflow-hidden">
         <div
           class="bg-primary h-full rounded-full transition-all duration-500 ease-out"
           :style="{ width: `${(taskLearningStore.tasks.filter(task => task.status === 'completed').length / Math.max(taskLearningStore.tasks.length, 1)) * 100}%` }"
         ></div>
       </div>
       <div class="flex flex-wrap gap-2 mt-3">
         <span
           v-for="task in taskLearningStore.tasks.filter(task => task.status === 'completed')"
           :key="task.task_learning_id"
           class="bg-success/10 dark:bg-success/20 text-success text-xs px-2 py-1 rounded-full"
         >
           {{ task.title.length > 20 ? task.title.substring(0, 20) + '...' : task.title }}
         </span>
       </div>
       <div v-if="topicStore.topic?.description" class="bg-surface-light dark:bg-surface-dark text-sm p-3 rounded-lg">
         <strong class="text-text-light dark:text-text-dark">Description:</strong>
         <span class="text-text-light dark:text-text-dark/80 ml-1">{{ topicStore.topic.description }}</span>
       </div>
     </section>

     <!-- Study Timer -->
     <section class="bg-white dark:bg-surface-dark p-6 shadow rounded-lg border border-surface-light dark:border-surface-dark">
       <h2 class="text-xl font-semibold text-text-light dark:text-text-dark mb-4">Study Session</h2>
       <div class="flex items-center justify-between">
         <div class="flex items-center gap-4">
           <div class="text-3xl font-mono font-bold text-primary">
             {{ formatTime(studyTimer.elapsedSeconds) }}
           </div>
           <div class="flex gap-2">
             <button
               v-if="!studyTimer.isRunning && studyTimer.elapsedSeconds === 0"
               @click="startStudyTimer"
               class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors"
             >
               Start Studying
             </button>
             <button
               v-else-if="studyTimer.isRunning"
               @click="pauseStudyTimer"
               class="px-4 py-2 bg-warning text-white rounded hover:bg-warning/90 transition-colors"
             >
               Pause
             </button>
             <button
               v-else-if="!studyTimer.isRunning && studyTimer.elapsedSeconds > 0"
               @click="startStudyTimer"
               class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors"
             >
               Resume
             </button>
             <button
               v-if="studyTimer.elapsedSeconds > 0"
               @click="stopStudyTimer"
               class="px-4 py-2 bg-success text-white rounded hover:bg-success/90 transition-colors"
             >
               Stop & Save
             </button>
           </div>
         </div>
         <div class="text-sm text-text-light dark:text-text-dark/70">
           {{ studyTimer.isRunning ? 'Studying...' : studyTimer.elapsedSeconds > 0 ? 'Paused' : 'Ready to study' }}
         </div>
       </div>
     </section>

     <!-- Tasks & Lessons -->
    <section class="space-y-4">
      <h2 class="text-xl font-semibold">Tasks & Lessons</h2>
      <p class="text-sm text-gray-500">Break down your learning into manageable tasks</p>

       <!-- Tasks List -->
       <div class="space-y-3">
         <template v-for="(task, index) in taskLearningStore.tasks" :key="task.task_learning_id">
            <div class="bg-white dark:bg-surface-dark p-4 rounded-lg shadow border border-surface-light dark:border-surface-dark hover:shadow-md transition-all duration-200 active:scale-[0.98] md:active:scale-100">
              <div class="flex items-start gap-3">
                <!-- Drag Handle (Desktop only) -->
                <div class="flex-shrink-0 mt-1 hidden md:block">
                  <svg class="w-4 h-4 text-text-light dark:text-text-dark/50 cursor-move" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
                  </svg>
                </div>

                <!-- Checkbox (larger on mobile) -->
                <input
                  type="checkbox"
                  :checked="isComplete(task.status)"
                  @click="completeTask(task)"
                  class="mt-1 flex-shrink-0 w-5 h-5 md:w-4 md:h-4"
                />

                <!-- Task Content -->
                <div class="flex flex-col flex-grow min-w-0">
                  <div class="flex items-start justify-between gap-2">
                    <span :class="[
                      'text-gray-900 dark:text-text-dark flex-grow text-base md:text-sm',
                      isComplete(task.status) && 'line-through text-text-light dark:text-text-dark/60'
                    ]">
                      {{ task.title }}
                    </span>

                    <!-- Status Badge -->
                    <span :class="[
                      'px-2 py-1 rounded-full text-xs font-medium flex-shrink-0',
                      task.status === 'completed' ? 'bg-success/10 dark:bg-success/20 text-success' :
                      task.status === 'in_progress' ? 'bg-primary/10 dark:bg-primary/20 text-primary' :
                      task.status === 'on_hold' ? 'bg-warning/10 dark:bg-warning/20 text-warning' :
                      'bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark/60'
                    ]">
                      {{ task.status.replace('_', ' ') }}
                    </span>
                  </div>

                 <!-- Task Notes -->
                 <div v-if="task.notes" class="mt-2 text-sm text-text-light dark:text-text-dark/70">
                   {{ task.notes }}
                 </div>

                  <!-- Task Actions -->
                  <div class="flex gap-2 mt-3">
                    <button
                      @click="editTask(task)"
                      class="flex-1 md:flex-none text-xs px-3 py-2 md:px-2 md:py-1 bg-surface-light dark:bg-surface-dark hover:bg-surface-light/80 dark:hover:bg-surface-dark/80 rounded transition-colors active:scale-95 md:active:scale-100"
                    >
                      Edit
                    </button>
                    <button
                      @click="deleteTask(task)"
                      class="flex-1 md:flex-none text-xs px-3 py-2 md:px-2 md:py-1 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/30 rounded transition-colors active:scale-95 md:active:scale-100"
                    >
                      Delete
                    </button>
                  </div>
               </div>
             </div>
           </div>
         </template>
       </div>
     </section>

     <!-- Resources Section -->
     <section class="bg-white dark:bg-surface-dark p-6 shadow rounded-lg border border-surface-light dark:border-surface-dark">
       <div class="flex justify-between items-center mb-4">
         <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Learning Resources</h2>
         <UiButton @click="openResourceModal()" variant="primary" size="sm">
           Add Resource
         </UiButton>
       </div>

       <div class="space-y-3">
         <template v-for="resource in resourceStore.resources" :key="resource.id">
           <div class="bg-surface-light dark:bg-surface-dark p-4 rounded-lg border border-surface-light dark:border-surface-dark hover:shadow-md transition-shadow">
             <div class="flex items-start gap-3">
               <div class="text-2xl">
                 {{ resourceStore.getResourceTypeIcon(resource.type) }}
               </div>
               <div class="flex-grow min-w-0">
                 <div class="flex items-start justify-between gap-2">
                   <div class="min-w-0 flex-grow">
                     <h3 class="font-medium text-text-light dark:text-text-dark truncate">
                       {{ resource.title }}
                     </h3>
                     <div class="flex items-center gap-2 mt-1">
                       <span :class="[
                         'px-2 py-1 rounded-full text-xs font-medium',
                         resourceStore.getResourceTypeColor(resource.type)
                       ]">
                         {{ resource.type }}
                       </span>
                       <template v-if="resource.link">
                         <a
                           :href="resource.link"
                           target="_blank"
                           rel="noopener noreferrer"
                           class="text-primary hover:text-primary/80 text-xs underline"
                         >
                           View Resource →
                         </a>
                       </template>
                     </div>
                     <template v-if="resource.notes">
                       <p class="text-sm text-text-light dark:text-text-dark/70 mt-2">
                         {{ resource.notes }}
                       </p>
                     </template>
                   </div>
                   <button
                     @click="deleteResource(resource.id)"
                     class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 p-1"
                     title="Delete resource"
                   >
                     <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                     </svg>
                   </button>
                 </div>
               </div>
             </div>
           </div>
         </template>

         <div v-if="resourceStore.resources.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
           <div class="text-4xl mb-2">📚</div>
           <p>No resources added yet</p>
           <p class="text-sm">Add articles, videos, documents, or other learning materials</p>
         </div>
       </div>
     </section>
   </div>

   <teleport to="body">
     <!-- Add Task Modal -->
     <div v-if="showModal.value" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
       <div class="bg-white dark:bg-surface-dark p-6 rounded shadow-md w-full max-w-md border border-surface-light dark:border-surface-dark">
         <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Add a New Task</h3>
         <form @submit.prevent="addTask">
           <div class="space-y-4">
             <input
               v-model="formData.title"
               placeholder="Task title"
               class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
               required
             />
           </div>
           <div class="flex justify-end gap-2 mt-6">
             <button type="button" @click="closeModal" class="px-4 py-2 bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded hover:bg-surface-light/80 dark:hover:bg-surface-dark/80 transition-colors">Cancel</button>
             <button type="submit" class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors">Save</button>
           </div>
         </form>
       </div>
     </div>

     <!-- Edit Task Modal -->
     <div v-if="showEditModal.value" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
       <div class="bg-white dark:bg-surface-dark p-6 rounded shadow-md w-full max-w-md border border-surface-light dark:border-surface-dark">
         <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Edit Task</h3>
         <form @submit.prevent="updateTask">
           <div class="space-y-4">
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Title</label>
               <input
                 v-model="editFormData.title"
                 placeholder="Task title"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
                 required
               />
             </div>
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Notes</label>
               <textarea
                 v-model="editFormData.notes"
                 placeholder="Task notes (optional)"
                 rows="3"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary resize-none"
               ></textarea>
             </div>
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Status</label>
               <select
                 v-model="editFormData.status"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
               >
                 <option value="not_started">Not Started</option>
                 <option value="in_progress">In Progress</option>
                 <option value="completed">Completed</option>
                 <option value="on_hold">On Hold</option>
               </select>
             </div>
           </div>
           <div class="flex justify-end gap-2 mt-6">
             <button type="button" @click="closeEditModal" class="px-4 py-2 bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded hover:bg-surface-light/80 dark:hover:bg-surface-dark/80 transition-colors">Cancel</button>
             <button type="submit" class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors">Update</button>
           </div>
         </form>
       </div>
     </div>

     <!-- Add Resource Modal -->
     <div v-if="showResourceModal.value" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
       <div class="bg-white dark:bg-surface-dark p-6 rounded shadow-md w-full max-w-md border border-surface-light dark:border-surface-dark">
         <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">
           Add Learning Resource
         </h3>
         <form @submit.prevent="addResource">
           <div class="space-y-4">
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Title *</label>
               <input
                 v-model="resourceFormData.title"
                 placeholder="Resource title"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
                 required
               />
             </div>
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Type</label>
               <select
                 v-model="resourceFormData.type"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
               >
                 <option value="article">📄 Article</option>
                 <option value="video">🎥 Video</option>
                 <option value="document">📋 Document</option>
                 <option value="book">📚 Book</option>
                 <option value="course">🎓 Course</option>
               </select>
             </div>
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Link</label>
               <input
                 v-model="resourceFormData.link"
                 type="url"
                 placeholder="https://example.com"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary"
               />
             </div>
             <div>
               <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Notes</label>
               <textarea
                 v-model="resourceFormData.notes"
                 placeholder="Additional notes (optional)"
                 rows="3"
                 class="w-full border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark p-2 rounded focus:outline-none focus:ring-2 focus:ring-primary resize-none"
               ></textarea>
             </div>
           </div>
           <div class="flex justify-end gap-2 mt-6">
             <button type="button" @click="closeResourceModal" class="px-4 py-2 bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded hover:bg-surface-light/80 dark:hover:bg-surface-dark/80 transition-colors">Cancel</button>
             <button type="submit" class="px-4 py-2 bg-primary text-white rounded hover:bg-primary/90 transition-colors">Add Resource</button>
           </div>
         </form>
       </div>
     </div>
   </teleport>

</template>
