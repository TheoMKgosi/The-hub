<script setup lang="ts">
const transactionStore = useTransactionStore()
const categoryStore = useCategoryStore()

const activeTransactionId = ref<string | null>(null)
const showDialog = ref(false)
const showTransactionModal = ref(true)
const searchQuery = ref('')
const transactionToDelete = ref<string>('')

const formData = reactive({
  description: '',
  amount: 0,
  type: 'expense' as 'income' | 'expense',
  date: '',
  category_id: ''
})


const items = ref(['expense', 'income'])
const value = ref('expense')

const categories = categoryStore.categories.map(category => category.name)

// Auto-categorization logic
const suggestedCategory = computed(() => {
  if (!formData.description || formData.type !== 'expense') return null

  const description = formData.description.toLowerCase()

  // Common expense keywords and their categories
  const categoryKeywords = {
    'food': ['food', 'restaurant', 'dining', 'lunch', 'dinner', 'coffee', 'cafe', 'meal', 'groceries', 'supermarket'],
    'transportation': ['gas', 'fuel', 'uber', 'taxi', 'bus', 'train', 'parking', 'toll', 'car', 'vehicle'],
    'entertainment': ['movie', 'cinema', 'theater', 'concert', 'game', 'gaming', 'netflix', 'spotify', 'subscription'],
    'shopping': ['amazon', 'shopping', 'clothes', 'clothing', 'store', 'mall', 'retail', 'purchase'],
    'utilities': ['electricity', 'water', 'gas bill', 'internet', 'phone', 'mobile', 'utility', 'power'],
    'healthcare': ['doctor', 'medical', 'pharmacy', 'hospital', 'health', 'insurance', 'dental', 'clinic'],
    'education': ['book', 'course', 'tuition', 'school', 'university', 'learning', 'study', 'class'],
    'housing': ['rent', 'mortgage', 'apartment', 'house', 'property', 'maintenance', 'repair']
  }

  for (const [categoryId, keywords] of Object.entries(categoryKeywords)) {
    if (keywords.some(keyword => description.includes(keyword))) {
      // Find the actual category from the store
      const category = categoryStore.categories.find(cat =>
        cat.name.toLowerCase().includes(categoryId) ||
        categoryId === 'food' && cat.name.toLowerCase().includes('food') ||
        categoryId === 'transportation' && cat.name.toLowerCase().includes('transport') ||
        categoryId === 'entertainment' && cat.name.toLowerCase().includes('entertain') ||
        categoryId === 'shopping' && cat.name.toLowerCase().includes('shopping') ||
        categoryId === 'utilities' && cat.name.toLowerCase().includes('utilit') ||
        categoryId === 'healthcare' && cat.name.toLowerCase().includes('health') ||
        categoryId === 'education' && cat.name.toLowerCase().includes('educat') ||
        categoryId === 'housing' && cat.name.toLowerCase().includes('hous')
      )
      if (category) return category
    }
  }

  return null
})

// Auto-select suggested category
watch(() => suggestedCategory.value, (newCategory) => {
  if (newCategory && !formData.category_id) {
    formData.category_id = newCategory.budget_category_id
  }
}, { immediate: true })

const filteredTransactions = computed(() => {
  let result = transactionStore.transactions

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(transaction =>
      transaction.description.toLowerCase().includes(query) ||
      transaction.category?.name.toLowerCase().includes(query)
    )
  }

  return result
})

const deleteItem = (id: string) => {
  transactionStore.deleteTransaction(id)
}

const handleTransactionAdd = async () => {
  console.log("Go here")
  const dataToSend = {
    ...formData,
    category_id: formData.category_id || undefined
  }
  transactionStore.submitForm(dataToSend)
  Object.assign(formData, {
    description: '',
    amount: 0,
    type: 'expense',
    date: '',
    category_id: ''
  })
  showTransactionModal.value = true
}

const handleCategorySelect = (category) => {
  formData.category_id = category.budget_category_id
}

const handleCategoryCreate = async (categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName })
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(cat =>
      cat.name.toLowerCase() === categoryName.toLowerCase()
    )
    if (newCategory) {
      formData.category_id = newCategory.budget_category_id
    }
  } catch (error) {
    console.error('Failed to create category:', error)
  }
}

const formatDate = (date: string) => new Date(date).toLocaleDateString()
const formatCurrency = (amount: number) => `$${amount.toFixed(2)}`

const openForm = (id: string) => {
  activeTransactionId.value = id
}

const closeForm = () => {
  activeTransactionId.value = null
}

onMounted(() => {
  if (transactionStore.transactions.length === 0) {
    transactionStore.fetchTransactions()
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})

// Mobile features
const isMobile = ref(false)
const swipeStartX = ref(0)
const swipeStartY = ref(0)
const swipeThreshold = 50

// Initialize mobile features
onMounted(() => {
  isMobile.value = window.innerWidth < 768
})

// Swipe gesture handlers
const handleTouchStart = (event) => {
  if (!isMobile.value) return
  swipeStartX.value = event.touches[0].clientX
  swipeStartY.value = event.touches[0].clientY
}

const handleTouchEnd = (event, transactionId) => {
  if (!isMobile.value) return

  const deltaX = event.changedTouches[0].clientX - swipeStartX.value
  const deltaY = event.changedTouches[0].clientY - swipeStartY.value

  // Check if it's a horizontal swipe (not vertical)
  if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > swipeThreshold) {
    if (deltaX > 0) {
      // Swipe right - edit
      openForm(transactionId)
    } else {
      // Swipe left - delete
      transactionToDelete.value = transactionId
      showDialog.value = true
    }
  }
}
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold">Transaction Management</h2>

    <!-- Button Actions -->
    <UModal title="Add Transaction">
      <UButton label="Add Transaction" variant="outline" />

      <template #body>
        <UForm @submit="handleTransactionAdd" :state="formData" class="space-y-4" v-model:open="showTransactionModal">
          <UFormField label="Description" required>
            <UInput v-model="formData.description" placeholder="e.g., Coffee, Salary, Rent" class="w-full" />
          </UFormField>

          <UFormField label="Amount" required>
            <UInputNumber v-model="formData.amount" :format-options="{ style: 'currency', currency: 'BWP' }" />
          </UFormField>

          <UFormField label="Type">
            <URadioGroup v-model="formData.type" :items="items" orientation="horizontal" />
          </UFormField>

          <div>
            <label for="date" class="block text-sm font-medium mb-1">Date</label>
            <input type="date" id="date" v-model="formData.date"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
              required />
          </div>

          <UFormField label="Category">
            <USelectMenu :items="categories" @select="handleCategorySelect" placeholder="Select Category" />
            <p v-if="formData.type === 'expense' && !formData.category_id"
              class="text-xs text-text-light dark:text-text-dark/60 mt-1">
              💡 Categories help track expenses against budgets
            </p>
          </UFormField>

          <!-- Modal Footer -->
          <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 ">
            <UButton label="Cancel" color="neutral" variant="outline" @click="showTransactionModal = false" />
            <UButton label="Create Transaction" variant="soft" color="neutral" type="submit" />
          </div>
        </UForm>
      </template>
    </UModal>

    <!-- Filters + Search -->
    <div
      class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search transactions..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <p class="text-sm text-text-light dark:text-text-dark/60 text-center">
      <span v-if="isMobile">Swipe left to delete, swipe right to edit</span>
      <span v-else>Double-click a transaction to delete it</span>
    </p>

    <!-- Transaction Cards -->
    <div class="space-y-4">
      <div v-if="transactionStore.transactions.length === 0"
        class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No transactions added yet</p>
        <p class="text-sm">Create your first transaction above to get started</p>
      </div>

      <div v-for="transaction in filteredTransactions" :key="transaction.transaction_id"
        class="p-6 rounded-lg shadow-md bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200"
        :class="{ 'cursor-pointer': isMobile }" @touchstart="handleTouchStart"
        @touchend="(event) => handleTouchEnd(event, transaction.transaction_id)"
        @dblclick="isMobile ? null : (transactionToDelete = transaction.transaction_id, showDialog = true)">
        <!-- Transaction Header -->
        <div class="flex justify-between items-start mb-4">
          <div class="flex-1">
            <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-1">{{ transaction.description }}
            </h3>
            <p class="text-sm text-text-light dark:text-text-dark/60">{{ formatDate(transaction.date) }}</p>
            <div class="flex items-center gap-2 mt-2">
              <span :class="[
                'px-2 py-1 rounded-full text-xs font-medium',
                transaction.type === 'income'
                  ? 'bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400'
                  : 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-400'
              ]">
                {{ transaction.type === 'income' ? 'Income' : 'Expense' }}
              </span>
              <span v-if="transaction.category"
                class="px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400">
                {{ transaction.category.name }}
              </span>
            </div>
          </div>
          <p :class="[
            'text-xl font-bold',
            transaction.type === 'income' ? 'text-success' : 'text-red-500 dark:text-red-400'
          ]">
            {{ transaction.type === 'income' ? '+' : '-' }}{{ formatCurrency(transaction.amount) }}
          </p>
        </div>
        <ConfirmDialog v-model:show="showDialog" :message="`Delete transaction '${transaction.description}'?`"
          @confirm="deleteItem(transactionToDelete)" />
      </div>
    </div>
  </div>
</template>
