<script setup lang="ts">
const goalStore = useGoalStore()

const searchQuery = ref('')
const filter = ref('all')

// Editing state
const editingGoalId = ref<string | null>(null)
const editFormData = reactive({
  title: '',
  description: '',
})

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
  Object.assign(editFormData, {
    title: goal.title,
    description: goal.description
  })
}

const cancelEdit = () => {
  editingGoalId.value = null
  Object.assign(editFormData, {
    title: '',
    description: ''
  })
}

const saveEdit = async (id: string) => {
  if (!editFormData.title.trim()) return

  const updatedGoal = {
    goal_id: id,
    title: editFormData.title.trim(),
    description: editFormData.description.trim(),
    tasks: null
  }

  await goalStore.updateGoal(updatedGoal)
  editingGoalId.value = null
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

    <div class="px-3 py-5 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md shadow-sm mt-4 rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <!-- Create Goal Form -->
      <FormGoal />

      <p v-if="goalStore.loading" class="text-text-light dark:text-text-dark">Loading...</p>

      <template v-else>
        <p v-if="goalStore.goals.length === 0" class="text-text-light dark:text-text-dark/60">No goals added yet</p>
        <div v-else-if="filteredGoals.length === 0" class="text-text-light dark:text-text-dark/60">No goals match your search</div>

        <div v-if="filteredGoals.length > 0" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3 mt-6">
          <div v-for="goal in filteredGoals" :key="goal.goal_id"
            class="bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark rounded-lg p-4 shadow-md hover:shadow-lg transition-shadow duration-200">

            <!-- Normal view -->
             <div v-if="editingGoalId !== goal.goal_id" class="flex flex-col h-full">
               <div @dblclick="startEdit(goal)" class="flex-1 cursor-pointer">
                 <h3 class="text-lg font-semibold text-text-light dark:text-text-dark mb-2">{{ goal.title }}</h3>
                 <p class="text-text-light dark:text-text-dark/80 mb-4">{{ goal.description }}</p>
               </div>

               <!-- Goal Tasks Section -->
               <GoalTasks :goal-id="goal.goal_id" class="mb-4" />

               <!-- Add Task to Goal -->
               <AddTaskToGoal
                 :goal-id="goal.goal_id"
                 :goal-title="goal.title"
                 class="mb-4"
               />

               <div class="flex items-center justify-between mt-auto">
                 <span class="text-sm text-text-light dark:text-text-dark/60">
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
            <div v-else class="flex flex-col w-full space-y-3">
              <input v-model="editFormData.title" placeholder="Goal title"
                class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />

              <textarea v-model="editFormData.description" placeholder="Goal description" rows="3"
                class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary resize-none"></textarea>

              <div class="flex gap-2">
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
