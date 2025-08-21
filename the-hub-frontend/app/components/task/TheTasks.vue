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
    <div class="shadow-sm p-3 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg mt-2 border border-surface-light/10 dark:border-surface-dark/10">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <div class="flex gap-2">
          <UiNavLink v-for="filterOption in ['all', 'pending', 'complete']" :key="filterOption"
            :active="filter === filterOption" variant="tab" @click="filter = filterOption">
            {{ filterOption.charAt(0).toUpperCase() + filterOption.slice(1) }}
          </UiNavLink>
        </div>
        <input v-model="searchQuery" placeholder="Search tasks..."
          class="flex-grow shadow-sm bg-surface-light dark:bg-surface-dark px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary w-full sm:w-0" />
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
          </template>
        </draggable>

      </template>
    </div>
  </div>
</template>
