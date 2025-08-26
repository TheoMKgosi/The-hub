<script setup lang="ts">
import ConfirmDialog from '@/components/ConfirmDialog.vue'

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

const submitForm = async () => {
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
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">Transaction Management</h2>

    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search transactions..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <!-- Floating Action Button -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showTransactionModal" @click="showTransactionModal = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
          <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
            <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
            </svg>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Transaction Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="!showTransactionModal" @click="showTransactionModal = true" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50">
          <div class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark" @click.stop>

            <!-- Modal Header -->
            <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
              <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Add New Transaction</h2>
              <UiButton @click="showTransactionModal = true" variant="default" size="sm" class="p-2">
                ×
              </UiButton>
            </div>

            <!-- Modal Body -->
            <div class="p-6">
              <form @submit.prevent="submitForm" class="space-y-4">

                <div>
                  <label for="description" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Description</label>
                  <input type="text" id="description" v-model="formData.description" placeholder="e.g., Coffee, Salary, Rent"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" required />
                </div>

                <div>
                  <label for="amount" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
                  <input type="number" id="amount" v-model="formData.amount" placeholder="0.00" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" required />
                </div>

                <div>
                  <label for="type" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Type</label>
                  <select id="type" v-model="formData.type"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" required>
                    <option value="expense">Expense</option>
                    <option value="income">Income</option>
                  </select>
                </div>

                <div>
                  <label for="date" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Date</label>
                  <input type="date" id="date" v-model="formData.date"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" required />
                </div>

                <div>
                  <label for="category" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Category (Optional)</label>
                  <select id="category" v-model="formData.category_id"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
                    <option value="">Select a category</option>
                    <option v-for="category in categoryStore.categories" :value="category.budget_category_id"
                      :key="category.budget_category_id">
                      {{ category.name }}
                    </option>
                  </select>
                </div>

                <!-- Modal Footer -->
                <div class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                  <UiButton type="button" @click="showTransactionModal = true" variant="default" size="md" class="w-full sm:w-auto">
                    Cancel
                  </UiButton>
                  <UiButton type="submit" variant="primary" size="md" class="w-full sm:w-auto">
                    Create Transaction
                  </UiButton>
                </div>

              </form>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <p class="text-sm text-text-light dark:text-text-dark/60 text-center">Double-click a transaction to delete it</p>

    <!-- Transaction Cards -->
    <div class="space-y-4">
      <div v-if="transactionStore.transactions.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No transactions added yet</p>
        <p class="text-sm">Create your first transaction above to get started</p>
      </div>

      <div v-for="transaction in filteredTransactions" :key="transaction.transaction_id"
        class="p-6 rounded-lg shadow-md bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
        <!-- Transaction Header -->
        <div class="flex justify-between items-start mb-4">
          <div class="flex-1">
            <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-1">{{ transaction.description }}</h3>
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
              <span v-if="transaction.category" class="px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400">
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

        <!-- Delete on double-click -->
        <div class="mt-4">
          <div class="text-center">
            <p class="text-sm text-text-light dark:text-text-dark/60">Double-click to delete</p>
          </div>
        </div>

        <ConfirmDialog v-model:show="showDialog" :message="`Delete transaction '${transaction.description}'?`"
          @confirm="deleteItem(transactionToDelete)" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.25s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.fade-scale-enter-to {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-from {
  opacity: 1;
  transform: scale(1);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>