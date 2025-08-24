<script setup lang="ts">
import draggable from 'vuedraggable'

const taskStore = useTaskStore()

callOnce(async () => {
  if (taskStore.tasks.length === 0) await taskStore.fetchTasks()
})

const filter = ref<'all' | 'complete' | 'pending' | 'linked' | 'priority-low' | 'priority-medium' | 'priority-high'>('all')
const priorityFilter = ref<'all' | '1' | '2' | '3' | '4' | '5'>('all')

const setFilter = (filterOption: string) => {
  filter.value = filterOption as any
}
const searchQuery = ref('')

const isFiltering = computed(() => (filter.value !== 'all' && filter.value !== 'linked' && !filter.value.startsWith('priority')) || priorityFilter.value !== 'all' || searchQuery.value !== '')
const matchFilter = (task) => {
  // Handle status filter
  if (filter.value === 'complete' || filter.value === 'pending') {
    if (task.status !== filter.value) return false
  }

  // Handle all filter - show only tasks not linked to goals
  if (filter.value === 'all') {
    if (task.goal_id) return false
  }

  // Handle linked filter - show only tasks linked to goals
  if (filter.value === 'linked') {
    if (!task.goal_id) return false
  }

  // Handle priority group filters
  if (filter.value === 'priority-low' && task.priority !== 1) return false
  if (filter.value === 'priority-medium' && (task.priority !== 2 && task.priority !== 3)) return false
  if (filter.value === 'priority-high' && (task.priority !== 4 && task.priority !== 5)) return false

  // Handle specific priority filter (from dropdown)
  if (priorityFilter.value !== 'all' && task.priority !== parseInt(priorityFilter.value)) return false

  // Handle search query
  if (searchQuery.value && !task.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    && !task.description.toLowerCase().includes(searchQuery.value.toLowerCase())) return false

  return true
}

const filteredTasks = computed(() => {
  return taskStore.tasks.filter(matchFilter)
})

// Editing a task
const editingTaskId = ref<number | null>(null)
const editFormData = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: 3,
  status: 'pending',
})

const startEdit = (task) => {
  editingTaskId.value = task.task_id
  Object.assign(editFormData, { ...task })
}

const cancelEdit = () => {
  editingTaskId.value = null
}

const saveEdit = async (id: number) => {
  const updated = { ...editFormData, task_id: id }
  if (updated.due_date) {
    const date = new Date(updated.due_date)
    updated.due_date = date.toISOString()
  }

  await taskStore.editTask(updated)
  editingTaskId.value = null
}


const completeTask = async (task) => {
  if (task.status == 'pending') {
    task.status = 'complete'
    taskStore.completeTask(task)
  } else {
    task.status = 'pending'
    taskStore.completeTask(task)
  }
}

const deleteTask = async (id: number) => {
  taskStore.deleteTask(id)
}

// Reordering a tasks
interface ReorderFormat {
  task_id: number
  order: number
}
// const taskReorderFormat = reactive({ task_id: 0, order: 0 })
// const taskReorder = ref([])
const reorderTasks = () => {
  const reorderedTasks = taskStore.tasks.map((task, index) => ({
    task_id: task.task_id,
    order: index + 1,
  }))
    taskStore.reorderTask(reorderedTasks)
}
</script>

<template>
  <div class="px-6">
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg mt-2 border border-surface-light/10 dark:border-surface-dark/10">
      <div class="flex flex-col gap-3">
        <!-- Main Filters -->
        <div class="flex flex-wrap gap-2 items-center">
          <div class="flex gap-1">
            <button v-for="filterOption in ['all', 'linked', 'pending', 'complete']" :key="filterOption"
              :class="[
                'px-4 py-2 text-sm font-medium transition-all duration-300 rounded-md',
                filter === filterOption
                  ? 'border-b-2 border-primary text-primary dark:text-primary bg-primary/10'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800'
              ]"
              @click="setFilter(filterOption)">
              {{ filterOption === 'linked' ? 'Linked' : filterOption.charAt(0).toUpperCase() + filterOption.slice(1) }}
            </button>
          </div>
        </div>

        <!-- Priority Filters -->
        <div class="flex flex-wrap gap-2 items-center">
          <span class="text-sm text-text-light dark:text-text-dark font-medium">Priority:</span>
          <div class="flex gap-1">
            <button v-for="priorityOption in ['priority-low', 'priority-medium', 'priority-high']" :key="priorityOption"
              :class="[
                'px-4 py-2 text-sm font-medium transition-all duration-300 rounded-md',
                filter === priorityOption
                  ? 'border-b-2 border-primary text-primary dark:text-primary bg-primary/10'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800'
              ]"
              @click="setFilter(priorityOption)">
              {{ priorityOption === 'priority-low' ? 'Low' :
                 priorityOption === 'priority-medium' ? 'Medium' : 'High' }}
            </button>
          </div>
          <select v-model="priorityFilter" class="px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-1 focus:ring-primary">
            <option value="all">All</option>
            <option value="1">1 - Low</option>
            <option value="2">2 - Medium</option>
            <option value="3">3 - Medium</option>
            <option value="4">4 - High</option>
            <option value="5">5 - High</option>
          </select>
        </div>

        <!-- Search -->
        <input v-model="searchQuery" placeholder="Search tasks..."
          class="w-full shadow-sm bg-surface-light dark:bg-surface-dark px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
      </div>
    </div>

    <div class="px-3 py-5 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md shadow-sm mt-4 rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <TaskFormTask />
      <p v-if="taskStore.loading" class="text-text-light dark:text-text-dark">Loading...</p>

       <template v-else>
         <p v-if="taskStore.tasks.length === 0" class="text-text-light dark:text-text-dark">No tasks found</p>
         <draggable v-model="taskStore.tasks" item-key="task_id" @end="reorderTasks"
           :disabled="isFiltering" class="space-y-3" ghost-class="opacity-50" chosen-class="shadow-lg"
           drag-class="rotate-2">
             <template #item="{ element: task }">
               <div v-if="matchFilter(task)" :data-task-id="task.task_id">
              <div :data-task-id="task.task_id">
                <div class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 cursor-move hover:shadow-lg transition-shadow duration-200"
                  :class="task.status === 'complete' ? 'border-success' : 'border-warning'">
                  <!-- Normal view -->
                  <div v-if="editingTaskId !== task.task_id" class="flex justify-between items-start">
                    <div @dblclick="startEdit(task)" class="flex-1">
                      <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">{{ task.title }}</h3>
                      <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">{{ task.description }}</p>
                      <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
                        {{ task.due_date ? new Date(task.due_date).toLocaleString() : '' }}
                      </p>
                      <div class="flex items-center gap-2 mt-2">
                        <input type="checkbox" @click="completeTask(task)" :checked="task.status === 'complete'"
                          class="accent-success w-4 h-4" />
                        <span class="text-sm font-medium text-text-light dark:text-text-dark capitalize">{{ task.status }}</span>
                      </div>
                      <p class="text-sm text-text-light dark:text-text-dark/60 mt-1">Priority: {{ task.priority }}</p>
                    </div>
                    <UiButton @click="deleteTask(task.task_id)" variant="danger" size="sm" class="ml-2"
                      title="Delete task">
                      Delete
                    </UiButton>
                  </div>

                <!-- Edit mode -->
                <div v-else class="flex flex-col w-full space-y-3">
                  <input v-model="editFormData.title" placeholder="Task title"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                  <input v-model="editFormData.description" placeholder="Task description"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                  <input type="datetime-local" v-model="editFormData.due_date"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                  <input type="number" min="1" max="5" v-model="editFormData.priority" placeholder="Priority (1-5)"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                  <div class="flex gap-2">
                    <UiButton @click="saveEdit(task.task_id)" variant="primary" size="sm">
                      Save
                    </UiButton>
                    <UiButton @click="cancelEdit" variant="default" size="sm">
                      Cancel
                    </UiButton>
                  </div>
                </div>
               </div>
             </div>
             </div>
           </template>
         </draggable>

      </template>
    </div>
  </div>
</template>
