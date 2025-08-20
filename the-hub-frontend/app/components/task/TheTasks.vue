<script setup lang="ts">
import draggable from 'vuedraggable'

const taskStore = useTaskStore()

callOnce(async () => {
  if (taskStore.tasks.length === 0) await taskStore.fetchTasks()
})

const filter = ref<'all' | 'complete' | 'pending'>('all')
const searchQuery = ref('')

const isFiltering = computed(() => filter.value !== 'all' || searchQuery.value !== '')
const matchFilter = (task) => {
  if (filter.value !== 'all' && task.status !== filter.value) return false
  if (searchQuery.value && !task.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    && !task.description.toLowerCase().includes(searchQuery.value.toLowerCase())) return false
  return true
}

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
    <div class="shadow-sm p-3 bg-white/20 backdrop-blur-md rounded-lg mt-2">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <div class="flex gap-2">
          <button @click="filter = 'all'" :class="{ 'font-bold': filter === 'all' }">All</button>
          <button @click="filter = 'pending'" :class="{ 'font-bold': filter === 'pending' }">Pending</button>
          <button @click="filter = 'complete'" :class="{ 'font-bold': filter === 'complete' }">Complete</button>
        </div>
        <input v-model="searchQuery" placeholder="Search tasks..." class="flex-grow shadow-sm bg-gradient-to-r from-gray-50 to-gray-100 px-3 py-2 rounded
          w-full sm:w-0" />
      </div>
    </div>

    <div class="px-3 py-5 bg-white/20 backdrop-blur-md shadow-sm mt-4 rounded-lg">
      <TaskFormTask />
      <p v-if="taskStore.loading">Loading...</p>

      <template v-else>
        <p v-if="taskStore.tasks.length === 0">No tasks found</p>
        <draggable v-model="taskStore.tasks" item-key="task_id" @end="reorderTasks"
          :disabled="isFiltering" class="space-y-3" ghost-class="opacity-50" chosen-class="shadow-lg"
          drag-class="rotate-2">
          <template #item="{ element: task }">
            <div :data-task-id="task.task_id">
              <div class="bg-white shadow rounded-lg p-4 border-l-4 cursor-move"
                :class="task.status === 'complete' ? 'border-green-500' : 'border-yellow-500'">
                <!-- Normal view -->
                <div v-if="editingTaskId !== task.task_id" class="flex justify-between items-start">
                  <div @dblclick="startEdit(task)">
                    <h3 class="text-lg font-semibold">{{ task.title }}</h3>
                    <p class="text-sm text-gray-600 mb-1">{{ task.description }}</p>
                    <p class="text-sm text-gray-500">
                      {{ task.due_date ? new Date(task.due_date).toLocaleString() : '' }}
                    </p>
                    <div class="flex items-center gap-2 mt-2">
                      <input type="checkbox" @click="completeTask(task)" :checked="task.status === 'complete'"
                        class="accent-green-600" />
                      <span class="text-sm font-medium">{{ task.status }}</span>
                    </div>
                    <p class="text-sm text-gray-500 mt-1">Priority: {{ task.priority }}</p>
                  </div>
                  <button @click="deleteTask(task.task_id)" class="text-red-500 hover:text-red-700 transition"
                    title="Delete task">
                    🗑
                  </button>
                </div>

                <!-- Edit mode -->
                <div v-else class="flex flex-col w-full space-y-2">
                  <input v-model="editFormData.title" class="border p-1 rounded" />
                  <input v-model="editFormData.description" class="border p-1 rounded" />
                  <input type="datetime-local" v-model="editFormData.due_date" class="border p-1 rounded" />
                  <input type="number" min="1" max="5" v-model="editFormData.priority" class="border p-1 rounded" />
                  <div class="flex gap-2">
                    <button @click="saveEdit(task.task_id)" class="bg-green-500 text-white px-3 py-1 rounded">
                      Save
                    </button>
                    <button @click="cancelEdit" class="bg-gray-400 text-white px-3 py-1 rounded">
                      Cancel
                    </button>
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
