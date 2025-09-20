<script setup lang="ts">
import FormGoal from '@/components/FormGoal.vue'

const goalStore = useGoalStore()
const { validateObject, schemas } = useValidation()

const searchQuery = ref('')
const filter = ref('all')

// Editing state
const editingGoalId = ref<string | null>(null)
const editFormData = reactive({
  title: '',
  description: '',
  due_date: '',
  priority: null as number | null,
  category: '',
  color: '#3B82F6',
  status: 'active'
})
const editValidationErrors = ref<Record<string, string>>({})

onMounted(() => {
  goalStore.fetchGoals()
})

// Computed properties for filtering
const isFiltering = computed(() => filter.value !== 'all' || searchQuery.value !== '')

const filteredGoals = computed(() => {
  return goalStore.goals.filter(goal => {
    // Search filter
    if (searchQuery.value &&
        !goal.title.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
        !goal.description.toLowerCase().includes(searchQuery.value.toLowerCase())) {
      return false
    }
    return true
  })
})

// Edit functionality
const startEdit = (goal: any) => {
  editingGoalId.value = goal.goal_id
  editValidationErrors.value = {}
  Object.assign(editFormData, {
    title: goal.title,
    description: goal.description,
    due_date: goal.due_date || '',
    priority: goal.priority || null,
    category: goal.category || '',
    color: goal.color || '#3B82F6',
    status: goal.status || 'active'
  })
}

const cancelEdit = () => {
  editingGoalId.value = null
  editValidationErrors.value = {}
  Object.assign(editFormData, {
    title: '',
    description: '',
    due_date: '',
    priority: null,
    category: '',
    color: '#3B82F6',
    status: 'active'
  })
}

const saveEdit = async (id: string) => {
  editValidationErrors.value = {}

  const payload = {
    title: editFormData.title.trim(),
    description: editFormData.description.trim(),
    due_date: editFormData.due_date ? new Date(editFormData.due_date).toISOString() : undefined,
    priority: editFormData.priority || undefined,
    category: editFormData.category.trim() || undefined,
    color: editFormData.color,
    status: editFormData.status
  }

  const validation = validateObject(payload, schemas.goal.update)

  if (!validation.isValid) {
    editValidationErrors.value = validation.errors
    return
  }

  const updatedGoal = {
    goal_id: id,
    ...payload,
    tasks: null
  }

  try {
    await goalStore.updateGoal(updatedGoal)
    editingGoalId.value = null
    editValidationErrors.value = {}
  } catch (err) {
    // Error is already handled in the store
  }
}

// Delete functionality
const deleteGoal = async (id: string) => {
  if (confirm('Are you sure you want to delete this goal?')) {
    await goalStore.deleteGoal(id)
  }
}
</script>

<template>
  <div class="px-6">
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg mt-2 border border-surface-light/10 dark:border-surface-dark/10">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <div class="flex gap-2">
          <UiNavLink v-for="filterOption in ['all']" :key="filterOption"
            :active="filter === filterOption" variant="tab" @click="filter = filterOption">
            {{ filterOption.charAt(0).toUpperCase() + filterOption.slice(1) }}
          </UiNavLink>
        </div>
        <input v-model="searchQuery" placeholder="Search goals..."
          class="flex-grow shadow-sm bg-surface-light dark:bg-surface-dark px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary w-full sm:w-0" />
      </div>
    </div>

    <div class="px-3 py-5 bg-surface-light/10 dark:bg-surface-dark/10 backdrop-blur-md shadow-sm mt-4 rounded-lg border border-surface-light/20 dark:border-surface-dark/20">
      <!-- Create Goal Form -->
      <FormGoal />

      <p v-if="goalStore.loading" class="text-text-light dark:text-text-dark">Loading...</p>

      <template v-else>
        <p v-if="goalStore.goals.length === 0" class="text-text-light dark:text-text-dark/60">No goals added yet</p>
        <div v-else-if="filteredGoals.length === 0" class="text-text-light dark:text-text-dark/60">No goals match your search</div>

        <div v-if="filteredGoals.length > 0" class="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 mt-4 sm:mt-6">
          <div v-for="goal in filteredGoals" :key="goal.goal_id"
            class="bg-surface-light/20 dark:bg-surface-dark/20 border border-surface-light/30 dark:border-surface-dark/30 rounded-lg p-4 sm:p-6 shadow-sm hover:shadow-md transition-all duration-200 hover:border-surface-light/40 dark:hover:border-surface-dark/40">

            <!-- Normal view -->
             <div v-if="editingGoalId !== goal.goal_id" class="flex flex-col h-full">
                <!-- Goal Header -->
                <div @dblclick="startEdit(goal)" class="flex-1 cursor-pointer mb-4">
                  <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">{{ goal.title }}</h3>
                  <p class="text-text-light dark:text-text-dark/80 text-sm leading-relaxed">{{ goal.description }}</p>
                </div>

                <!-- Goal Progress Section -->
                <div class="mb-4">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-sm font-medium text-text-light dark:text-text-dark">Progress</span>
                    <span class="text-sm text-text-light dark:text-text-dark/70">
                      {{ goal.completed_tasks }}/{{ goal.total_tasks }} tasks
                    </span>
                  </div>
                  <div class="w-full bg-surface-light/30 dark:bg-surface-dark/30 rounded-full h-2">
                    <div
                      class="h-2 rounded-full transition-all duration-300"
                      :class="goal.progress === 100 ? 'bg-success' : 'bg-primary'"
                      :style="{ width: goal.progress + '%' }"
                    ></div>
                  </div>
                  <div class="text-right mt-1">
                    <span class="text-xs text-text-light dark:text-text-dark/60">
                      {{ Math.round(goal.progress) }}% complete
                    </span>
                  </div>
                </div>

                <!-- Goal Metadata -->
                <div class="mb-4 flex flex-wrap gap-2">
                  <span v-if="goal.due_date"
                    class="px-2 py-1 text-xs font-medium rounded-full"
                    :class="new Date(goal.due_date) < new Date() && goal.status !== 'completed'
                      ? 'bg-error/10 dark:bg-error/20 text-error dark:text-error'
                      : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
                    Due: {{ new Date(goal.due_date).toLocaleDateString() }}
                  </span>
                  <span v-if="goal.priority"
                    class="px-2 py-1 text-xs font-medium bg-secondary/10 dark:bg-secondary/20 text-secondary dark:text-secondary rounded-full">
                    Priority {{ goal.priority }}
                  </span>
                  <span v-if="goal.category"
                    class="px-2 py-1 text-xs font-medium bg-primary/10 dark:bg-primary/20 text-primary dark:text-primary rounded-full">
                    {{ goal.category }}
                  </span>
                  <span class="px-2 py-1 text-xs font-medium rounded-full"
                    :class="goal.status === 'completed'
                      ? 'bg-success/10 dark:bg-success/20 text-success dark:text-success'
                      : goal.status === 'active'
                      ? 'bg-primary/10 dark:bg-primary/20 text-primary dark:text-primary'
                      : 'bg-warning/10 dark:bg-warning/20 text-warning dark:text-warning'">
                    {{ goal.status }}
                  </span>
                </div>

                <!-- Goal Tasks Section -->
                <div class="mb-4">
                  <GoalTasks :goal-id="goal.goal_id" />
                </div>

                <!-- Add Task to Goal -->
                <div class="mb-4">
                  <AddTaskToGoal
                    :goal-id="goal.goal_id"
                    :goal-title="goal.title"
                  />
                </div>

                <!-- Goal Actions -->
                <div class="flex items-center justify-between mt-auto pt-4 border-t border-surface-light/20 dark:border-surface-dark/20">
                  <span class="text-xs text-text-light dark:text-text-dark/60">
                    Double-click to edit
                  </span>
                  <div class="flex gap-2">
                    <UiButton @click="startEdit(goal)" variant="default" size="sm">
                      Edit
                    </UiButton>
                    <UiButton @click="deleteGoal(goal.goal_id)" variant="danger" size="sm">
                      Delete
                    </UiButton>
                  </div>
                </div>
              </div>

              <!-- Edit mode -->
              <div v-else class="flex flex-col w-full space-y-4">
                <div class="space-y-3">
                  <div class="flex flex-col">
                    <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Title</label>
                    <input v-model="editFormData.title" placeholder="Goal title"
                      class="w-full border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors"
                      :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.title }" />
                    <p v-if="editValidationErrors.title" class="mt-1 text-sm text-red-500 dark:text-red-400">
                      {{ editValidationErrors.title }}
                    </p>
                  </div>

                  <div class="flex flex-col">
                    <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Description</label>
                    <textarea v-model="editFormData.description" placeholder="Goal description" rows="3"
                      class="w-full border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary resize-none transition-colors"
                      :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.description }"></textarea>
                    <p v-if="editValidationErrors.description" class="mt-1 text-sm text-red-500 dark:text-red-400">
                      {{ editValidationErrors.description }}
                    </p>
                  </div>

                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Due Date</label>
                      <input type="date" v-model="editFormData.due_date"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors"
                        :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.due_date }" />
                      <p v-if="editValidationErrors.due_date" class="mt-1 text-sm text-red-500 dark:text-red-400">
                        {{ editValidationErrors.due_date }}
                      </p>
                    </div>

                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Priority</label>
                      <select v-model="editFormData.priority"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors"
                        :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.priority }">
                        <option :value="null">No priority</option>
                        <option :value="1">1 - Low</option>
                        <option :value="2">2</option>
                        <option :value="3">3 - Medium</option>
                        <option :value="4">4</option>
                        <option :value="5">5 - High</option>
                      </select>
                      <p v-if="editValidationErrors.priority" class="mt-1 text-sm text-red-500 dark:text-red-400">
                        {{ editValidationErrors.priority }}
                      </p>
                    </div>
                  </div>

                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Category</label>
                      <input type="text" v-model="editFormData.category" placeholder="e.g., Work, Personal, Health"
                        class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors placeholder:text-text-light/50 dark:placeholder:text-text-dark/50"
                        :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.category }" />
                      <p v-if="editValidationErrors.category" class="mt-1 text-sm text-red-500 dark:text-red-400">
                        {{ editValidationErrors.category }}
                      </p>
                    </div>

                    <div class="flex flex-col">
                      <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Color</label>
                      <input type="color" v-model="editFormData.color"
                        class="w-full h-10 px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors cursor-pointer"
                        :class="{ 'border-red-500 focus:ring-red-500': editValidationErrors.color }" />
                      <p v-if="editValidationErrors.color" class="mt-1 text-sm text-red-500 dark:text-red-400">
                        {{ editValidationErrors.color }}
                      </p>
                    </div>
                  </div>

                  <div class="flex flex-col">
                    <label class="mb-2 font-medium text-sm text-text-light dark:text-text-dark">Status</label>
                    <select v-model="editFormData.status"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary transition-colors">
                      <option value="active">Active</option>
                      <option value="completed">Completed</option>
                      <option value="paused">Paused</option>
                    </select>
                  </div>
                </div>

                <div class="flex gap-2 pt-2 border-t border-surface-light/20 dark:border-surface-dark/20">
                  <UiButton @click="saveEdit(goal.goal_id)" variant="primary" size="sm" :disabled="!editFormData.title.trim()">
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
    </div>
  </div>
</template>
