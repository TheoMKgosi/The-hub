<script setup lang="ts">
const goalStore = useFinancialGoalStore()
const categoryStore = useCategoryStore()
const toast = useToast()

const showForm = ref(false)
const editingGoalId = ref<string | null>(null)
const contributingGoalId = ref<string | null>(null)
const contributeAmount = ref(0)
const activeTab = ref<'savings' | 'checklist'>('savings')

const formData = reactive({
  title: '',
  description: '',
  type: 'savings' as 'savings' | 'checklist',
  target_amount: null as number | null,
  current_amount: 0,
  priority: null as number | null,
  target_date: null as string | null,
  category_id: null as string | null,
  color: '#3B82F6',
})

const editData = reactive({
  goal_id: '',
  title: '',
  description: '',
  type: 'savings' as 'savings' | 'checklist',
  target_amount: null as number | null,
  current_amount: 0,
  priority: null as number | null,
  target_date: null as string | null,
  category_id: null as string | null,
  color: '#3B82F6',
})

const savingsGoals = computed(() =>
  goalStore.goals.filter(g => g.type === 'savings')
)
const checklistItems = computed(() =>
  goalStore.goals.filter(g => g.type === 'checklist')
)
const activeCategoryGoals = computed(() => {
  const activeOnly = [...savingsGoals.value, ...checklistItems.value].filter(g => g.status === 'active')
  const completed = [...savingsGoals.value, ...checklistItems.value].filter(g => g.status === 'completed')
  return { active: activeOnly, completed }
})

function resetForm() {
  Object.assign(formData, {
    title: '',
    description: '',
    type: 'savings',
    target_amount: null,
    current_amount: 0,
    priority: null,
    target_date: null,
    category_id: null,
    color: '#3B82F6',
  })
}

async function submitForm() {
  if (!formData.title.trim()) {
    toast.add({ title: "Error", description: "Title is required", color: "error" })
    return
  }

  const payload: any = {
    title: formData.title,
    description: formData.description,
    type: formData.type,
    current_amount: formData.current_amount,
    priority: formData.priority,
    target_date: formData.target_date ? new Date(formData.target_date).toISOString() : null,
    category_id: formData.category_id,
    color: formData.color,
  }

  if (formData.type === 'savings') {
    payload.target_amount = formData.target_amount
  }

  await goalStore.submitForm(payload)
  showForm.value = false
  resetForm()
}

function startEdit(goal: any) {
  editingGoalId.value = goal.goal_id
  Object.assign(editData, {
    goal_id: goal.goal_id,
    title: goal.title,
    description: goal.description,
    type: goal.type,
    target_amount: goal.target_amount,
    current_amount: goal.current_amount,
    priority: goal.priority,
    target_date: goal.target_date ? new Date(goal.target_date).toISOString().slice(0, 10) : null,
    category_id: goal.category_id,
    color: goal.color || '#3B82F6',
  })
}

function cancelEdit() {
  editingGoalId.value = null
  Object.assign(editData, {
    goal_id: '',
    title: '',
    description: '',
    type: 'savings',
    target_amount: null,
    current_amount: 0,
    priority: null,
    target_date: null,
    category_id: null,
    color: '#3B82F6',
  })
}

async function submitEdit() {
  if (!editData.title.trim()) {
    toast.add({ title: "Error", description: "Title is required", color: "error" })
    return
  }

  const payload: any = {
    goal_id: editData.goal_id,
    title: editData.title,
    description: editData.description,
    type: editData.type,
    priority: editData.priority,
    target_date: editData.target_date ? new Date(editData.target_date).toISOString() : null,
    category_id: editData.category_id,
    color: editData.color,
    target_amount: editData.type === 'savings' ? editData.target_amount : null,
  }

  await goalStore.editGoal(payload)
  cancelEdit()
}

async function deleteGoalItem(goalId: string) {
  if (confirm('Delete this goal?')) {
    await goalStore.deleteGoal(goalId)
  }
}

async function contributeAmountToGoal(goalId: string) {
  if (contributeAmount.value <= 0) {
    toast.add({ title: "Error", description: "Enter a valid amount", color: "error" })
    return
  }

  await goalStore.addProgress(goalId, contributeAmount.value)
  contributeAmount.value = 0
  contributingGoalId.value = null
}

function getProgressPercent(goal: any) {
  if (!goal.target_amount || goal.target_amount <= 0) return 0
  return Math.min((goal.current_amount / goal.target_amount) * 100, 100)
}

function getProgressClass(goal: any) {
  const pct = getProgressPercent(goal)
  if (pct >= 100) return 'bg-green-500'
  if (pct >= 75) return 'bg-blue-500'
  if (pct >= 50) return 'bg-yellow-500'
  return 'bg-gray-400'
}

function getDaysLeft(dateStr: string | null) {
  if (!dateStr) return null
  const days = Math.ceil((new Date(dateStr).getTime() - Date.now()) / (1000 * 60 * 60 * 24))
  return days
}

function getDaysLeftClass(days: number | null) {
  if (days === null) return ''
  if (days < 0) return 'text-red-600 dark:text-red-400'
  if (days <= 7) return 'text-orange-600 dark:text-orange-400'
  return 'text-text-light/60 dark:text-text-dark/60'
}

function getDaysLeftText(days: number | null) {
  if (days === null) return ''
  if (days < 0) return `${Math.abs(days)}d overdue`
  if (days === 0) return 'Due today'
  return `${days}d left`
}

function getCategoryName(categoryId: string | null) {
  if (!categoryId) return 'None'
  return categoryStore.categories.find(c => c.budget_category_id === categoryId)?.name || 'Unknown'
}

const handleCategorySelect = (category: any) => {
  formData.category_id = category.budget_category_id
}

const handleEditCategorySelect = (category: any) => {
  editData.category_id = category.budget_category_id
}

const handleCategoryCreate = async (categoryName: string) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    const newCategory = categoryStore.categories.find(
      c => c.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) formData.category_id = newCategory.budget_category_id
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

const handleEditCategoryCreate = async (categoryName: string) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    const newCategory = categoryStore.categories.find(
      c => c.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) editData.category_id = newCategory.budget_category_id
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

onMounted(async () => {
  if (goalStore.goals.length === 0) goalStore.fetchGoals()
  if (categoryStore.categories.length === 0) categoryStore.fetchCategory()
})

const colorOptions = [
  { value: '#3B82F6', label: 'Blue' },
  { value: '#10B981', label: 'Green' },
  { value: '#F59E0B', label: 'Amber' },
  { value: '#EF4444', label: 'Red' },
  { value: '#8B5CF6', label: 'Purple' },
  { value: '#EC4899', label: 'Pink' },
]
</script>

<template>
  <div class="space-y-6 p-4 max-w-4xl mx-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">Financial Goals</h2>
      <UButton :label="showForm ? 'Cancel' : 'New Goal'" :icon="showForm ? '' : 'i-lucide-plus'"
        :color="showForm ? 'error' : 'primary'" @click="showForm = !showForm" />
    </div>

    <!-- Create Form -->
    <UCollapsible v-model:open="showForm">
      <template #content>
        <div
          class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark">
          <h3 class="text-lg font-semibold mb-3">Create Goal</h3>
          <div class="space-y-3">
            <UFormField label="Title" required>
              <UInput v-model="formData.title" placeholder="e.g., WiFi Extender, Course Exam" class="w-full"
                size="lg" required/>
            </UFormField>
            <UFormField label="Description">
              <UInput v-model="formData.description" placeholder="Note about this goal..." class="w-full"
                size="lg" />
            </UFormField>
            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Type</label>
              <div class="flex gap-3">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input v-model="formData.type" type="radio" value="savings" class="text-primary focus:ring-primary" />
                  <span class="text-text-light dark:text-text-dark">Savings Goal</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input v-model="formData.type" type="radio" value="checklist"
                    class="text-primary focus:ring-primary" />
                  <span class="text-text-light dark:text-text-dark">Checklist</span>
                </label>
              </div>
            </div>
            <template v-if="formData.type === 'savings'">
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Target
                    Amount</label>
                  <input v-model.number="formData.target_amount" type="number" placeholder="0.00" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Already
                    Saved</label>
                  <input v-model.number="formData.current_amount" type="number" placeholder="0.00" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Target Date</label>
                  <input v-model="formData.target_date" type="date"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Priority</label>
                  <select v-model.number="formData.priority"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                    <option :value="null">None</option>
                    <option :value="1">1 - Highest</option>
                    <option :value="2">2 - High</option>
                    <option :value="3">3 - Medium</option>
                    <option :value="4">4 - Low</option>
                    <option :value="5">5 - Lowest</option>
                  </select>
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Link to Budget
                  Category</label>
                <ComboBox :model-value="formData.category_id" :categories="categoryStore.categories"
                  placeholder="Optional - link to budget category" @select="handleCategorySelect"
                  @create="handleCategoryCreate" />
              </div>
            </template>
            <div>
              <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Color</label>
              <div class="flex gap-2">
                <button v-for="c in colorOptions" :key="c.value" type="button" @click="formData.color = c.value"
                  :class="['w-8 h-8 rounded-full border-2 transition-all', formData.color === c.value ? 'border-text-light dark:border-text-dark scale-110' : 'border-transparent']"
                  :style="{ backgroundColor: c.value }" :title="c.label" />
              </div>
            </div>
            <div class="flex justify-end">
              <UButton :label="goalStore.creating ? 'Creating...' : 'Create Goal'" @click="submitForm" />
            </div>
          </div>
        </div>
      </template>
    </UCollapsible>

    <!-- Edit Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="editingGoalId"
          class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50"
          @click.self="cancelEdit">
          <div
            class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark p-6"
            @click.stop>
            <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Edit Goal</h3>
            <div class="space-y-3">
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Title</label>
                <input v-model="editData.title" type="text"
                  class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
              </div>
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Description</label>
                <input v-model="editData.description" type="text"
                  class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
              </div>
              <template v-if="editData.type === 'savings'">
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Target
                    Amount</label>
                  <input v-model.number="editData.target_amount" type="number" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Target Date</label>
                  <input v-model="editData.target_date" type="date"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Priority</label>
                  <select v-model.number="editData.priority"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                    <option :value="null">None</option>
                    <option :value="1">1 - Highest</option>
                    <option :value="2">2 - High</option>
                    <option :value="3">3 - Medium</option>
                    <option :value="4">4 - Low</option>
                    <option :value="5">5 - Lowest</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Link to Budget
                    Category</label>
                  <ComboBox :model-value="editData.category_id" :categories="categoryStore.categories"
                    placeholder="Optional - link to budget category" @select="handleEditCategorySelect"
                    @create="handleEditCategoryCreate" />
                </div>
              </template>
              <div>
                <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Color</label>
                <div class="flex gap-2">
                  <button v-for="c in colorOptions" :key="c.value" type="button" @click="editData.color = c.value"
                    :class="['w-8 h-8 rounded-full border-2 transition-all', editData.color === c.value ? 'border-text-light dark:border-text-dark scale-110' : 'border-transparent']"
                    :style="{ backgroundColor: c.value }" :title="c.label" />
                </div>
              </div>
              <div class="flex justify-end gap-2 pt-2">
                <UiButton variant="default" size="sm" @click="cancelEdit">Cancel</UiButton>
                <UiButton variant="primary" size="sm" :disabled="goalStore.updating" @click="submitEdit">
                  {{ goalStore.updating ? 'Saving...' : 'Save' }}
                </UiButton>
              </div>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Tab Navigation for Active Goals -->
    <div class="flex border-b border-surface-light dark:border-surface-dark">
      <button @click="activeTab = 'savings'"
        :class="['px-4 py-2 text-sm font-medium border-b-2 transition-colors', activeTab === 'savings' ? 'border-primary text-primary' : 'border-transparent text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark']">
        Savings Goals ({{savingsGoals.filter(g => g.status === 'active').length}})
      </button>
      <button @click="activeTab = 'checklist'"
        :class="['px-4 py-2 text-sm font-medium border-b-2 transition-colors', activeTab === 'checklist' ? 'border-primary text-primary' : 'border-transparent text-text-light/60 dark:text-text-dark/60 hover:text-text-light dark:hover:text-text-dark']">
        Checklist ({{checklistItems.filter(g => g.status === 'active').length}})
      </button>
    </div>

    <!-- Savings Goals -->
    <div v-if="activeTab === 'savings'">
      <div v-if="savingsGoals.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No savings goals yet</p>
        <p class="text-sm">Create your first goal above</p>
      </div>

      <!-- Active Savings Goals -->
      <div v-if="savingsGoals.filter(g => g.status === 'active').length > 0" class="space-y-3">
        <div v-for="goal in savingsGoals" :key="goal.goal_id" v-show="goal.status === 'active'"
          class="p-4 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
          <div class="flex items-start justify-between mb-2">
            <div class="flex-1">
              <div class="flex items-center gap-2">
                <span class="w-3 h-3 rounded-full flex-shrink-0" :style="{ backgroundColor: goal.color }"></span>
                <h3 class="font-semibold text-text-light dark:text-text-dark">{{ goal.title }}</h3>
              </div>
              <p v-if="goal.description" class="text-sm text-text-light/60 dark:text-text-dark/60 mt-1">{{
                goal.description
                }}</p>
              <div class="flex items-center gap-3 mt-1 text-xs text-text-light/60 dark:text-text-dark/60">
                <span v-if="goal.priority" class="px-1.5 py-0.5 rounded bg-surface-light/50 dark:bg-surface-dark/50">
                  P{{ goal.priority }}
                </span>
                <span :class="getDaysLeftClass(getDaysLeft(goal.target_date))">
                  {{ getDaysLeftText(getDaysLeft(goal.target_date)) }}
                </span>
                <span v-if="goal.category_id">Linked: {{ getCategoryName(goal.category_id) }}</span>
              </div>
            </div>
            <div class="flex items-center gap-1">
              <button @click="contributeAmountToGoal(goal.goal_id); contributingGoalId = goal.goal_id"
                class="p-1.5 bg-primary/10 hover:bg-primary/20 rounded-lg text-primary transition-colors"
                title="Add to savings">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
              </button>
              <button @click="startEdit(goal)" class="p-1.5 text-blue-600 hover:text-blue-800 rounded-lg" title="Edit">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button @click="deleteGoalItem(goal.goal_id)" class="p-1.5 text-red-600 hover:text-red-800 rounded-lg"
                title="Delete">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Contribute Input -->
          <div v-if="contributingGoalId === goal.goal_id" class="mt-2 flex items-center gap-2">
            <input v-model.number="contributeAmount" type="number" placeholder="Amount" step="0.01" min="0.01"
              class="flex-1 px-3 py-1.5 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-1 focus:ring-primary" />
            <UiButton variant="primary" size="xs" @click="contributeAmountToGoal(goal.goal_id)">Add</UiButton>
            <UiButton variant="default" size="xs" @click="contributingGoalId = null; contributeAmount = 0">Cancel
            </UiButton>
          </div>

          <!-- Progress -->
          <div v-if="goal.target_amount" class="mt-2">
            <div class="flex justify-between text-xs mb-1">
              <span class="text-text-light dark:text-text-dark">${{ goal.current_amount.toFixed(0) }}</span>
              <span class="text-text-light/60 dark:text-text-dark/60">{{ getProgressPercent(goal).toFixed(0) }}%</span>
              <span class="text-text-light dark:text-text-dark">${{ goal.target_amount.toFixed(0) }}</span>
            </div>
            <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div :class="['h-2 rounded-full transition-all duration-300', getProgressClass(goal)]"
                :style="{ width: getProgressPercent(goal) + '%' }"></div>
            </div>
            <p v-if="goal.target_amount" class="text-xs text-text-light/60 dark:text-text-dark/60 mt-1 text-right">
              ${{ (goal.target_amount - goal.current_amount).toFixed(2) }} remaining
            </p>
          </div>
        </div>
      </div>

      <!-- Completed Goals -->
      <div v-if="savingsGoals.filter(g => g.status === 'completed').length > 0" class="mt-6">
        <h3 class="text-lg font-semibold mb-3 text-text-light dark:text-text-dark">Completed</h3>
        <div v-for="goal in savingsGoals" :key="'done-' + goal.goal_id" v-show="goal.status === 'completed'"
          class="p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 opacity-75 mb-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full flex-shrink-0" :style="{ backgroundColor: goal.color }"></span>
              <div>
                <h4 class="font-medium text-text-light dark:text-text-dark line-through">{{ goal.title }}</h4>
                <p v-if="goal.target_amount" class="text-xs text-green-600 dark:text-green-400">
                  ${{ goal.current_amount.toFixed(2) }} / ${{ goal.target_amount.toFixed(2) }}
                </p>
              </div>
            </div>
            <div class="flex gap-1">
              <button @click="goalStore.editGoal({ goal_id: goal.goal_id, status: 'active' })"
                class="p-1 text-blue-600 hover:text-blue-800 text-xs">Reactivate</button>
              <button @click="deleteGoalItem(goal.goal_id)" class="p-1 text-red-600 hover:text-red-800">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Checklist -->
    <div v-if="activeTab === 'checklist'">
      <div v-if="checklistItems.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No checklist items yet</p>
        <p class="text-sm">Add things you plan to buy or do eventually</p>
      </div>

      <div v-if="checklistItems.filter(g => g.status === 'active').length > 0" class="space-y-2">
        <div v-for="item in checklistItems" :key="item.goal_id" v-show="item.status === 'active'"
          class="p-3 rounded-lg bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark flex items-center gap-3 hover:shadow transition-shadow">
          <button @click="goalStore.toggleChecklist(item.goal_id)" :class="['w-5 h-5 rounded border-2 flex items-center justify-center flex-shrink-0 transition-colors',
            item.status === 'completed' ? 'bg-primary border-primary' : 'border-gray-300 dark:border-gray-600']">
            <svg v-if="item.status === 'completed'" class="w-3 h-3 text-white" fill="none" stroke="currentColor"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
            </svg>
          </button>
          <div class="flex-1 min-w-0">
            <p class="font-medium text-text-light dark:text-text-dark truncate">{{ item.title }}</p>
            <p v-if="item.description" class="text-xs text-text-light/60 dark:text-text-dark/60 truncate">{{
              item.description }}</p>
          </div>
          <div class="flex items-center gap-2 flex-shrink-0">
            <span v-if="item.target_amount" class="text-sm font-medium text-text-light dark:text-text-dark">
              ~${{ item.target_amount.toFixed(2) }}
            </span>
            <span v-if="item.target_date" :class="['text-xs', getDaysLeftClass(getDaysLeft(item.target_date))]">
              {{ getDaysLeftText(getDaysLeft(item.target_date)) }}
            </span>
            <span v-if="item.category_id"
              class="text-xs px-1.5 py-0.5 rounded bg-surface-light/50 dark:bg-surface-dark/50 text-text-light/60">
              {{ getCategoryName(item.category_id) }}
            </span>
            <button @click="startEdit(item)" class="p-1 text-blue-600 hover:text-blue-800">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
            <button @click="deleteGoalItem(item.goal_id)" class="p-1 text-red-600 hover:text-red-800">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Completed Checklist -->
      <div v-if="checklistItems.filter(g => g.status === 'completed').length > 0" class="mt-6">
        <h3 class="text-lg font-semibold mb-3 text-text-light dark:text-text-dark">Completed</h3>
        <div v-for="item in checklistItems" :key="'done-' + item.goal_id" v-show="item.status === 'completed'"
          class="p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 opacity-75 mb-2 flex items-center gap-3">
          <svg class="w-5 h-5 text-green-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <div class="flex-1">
            <p class="font-medium text-text-light dark:text-text-dark line-through">{{ item.title }}</p>
          </div>
          <div class="flex gap-1">
            <UButton label="Undo" @click="goalStore.toggleChecklist(item.goal_id)" />
            <UButton icon="i-lucide-trash-2" color="error" @click="deleteGoalItem(item.goal_id)" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
