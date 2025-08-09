<script setup lang="ts">
import { ref, reactive, computed, onMounted } from "vue";
import { useIncomeStore } from "@/stores/income";
import { useCategoryStore, useBudgetStore } from "@/stores/finance";
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const incomeStore = useIncomeStore()
const categoryStore = useCategoryStore()
const budgetStore = useBudgetStore()

const showForm = ref(true)
const showDialog = ref(false)
const searchQuery = ref('')

const formData = reactive({
  source: '',
  amount: 0,
  received_at: null
})

const budgetForm = reactive({
  income_id: 0,
  category_id: 0,
  amount: 0,
  start_date: null,
  end_date: null
})

const filteredTasks = computed(() => {
  let result = incomeStore.incomes

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(income =>
      income.source.toLowerCase().includes(query)
    )
  }

  return result
})

const deleteItem = (id) => {
  budgetStore.deleteBudget(id)
}

const submitForm = async () => {
  const dataToSend = { ...formData }
  incomeStore.submitForm(dataToSend)
  Object.assign(formData, {
    source: '',
    amount: '',
    received_at: null
  })
}

const submitBudgetForm = async () => {
  const dataToSend = { ...budgetForm }
  budgetStore.submitForm(dataToSend)
  Object.assign(budgetForm, {
    income_id: 0,
    category_id: 0,
    amount: 0,
    start_date: null,
    end_date: null
  })
}

onMounted(() => {
  if (incomeStore.incomes.length === 0) {
    incomeStore.fetchIncomes()
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
})

const formatDate = (date) => new Date(date).toLocaleDateString()
const openForm = (id) => {
  showForm.value = false
  budgetForm.income_id = id
}

const remainingAmount = (amount, budgets) => {
  let remaining = amount
  for (const budget of budgets) {
    remaining -= budget.amount
  }
  return remaining
}
</script>

<template>
  <div class="space-y-6 p-4">
    <!-- Filters + Search -->
    <div class="shadow-sm p-4 bg-white/30 backdrop-blur-md rounded-xl">
      <input
        v-model="searchQuery"
        placeholder="Search tasks..."
        class="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-400"
      />
    </div>

    <!-- Income Form -->
    <form
      @submit.prevent="submitForm"
      class="space-y-4 p-6 max-w-lg mx-auto bg-white rounded-xl shadow-lg"
    >
      <h2 class="text-lg font-semibold text-gray-800">Add Income</h2>

      <div>
        <label for="source" class="block text-sm font-medium text-gray-700">Income Source</label>
        <input
          type="text"
          id="source"
          v-model="formData.source"
          class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
        />
      </div>

      <div>
        <label for="amount" class="block text-sm font-medium text-gray-700">Amount</label>
        <input
          type="number"
          id="amount"
          v-model="formData.amount"
          class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
        />
      </div>

      <div>
        <label for="received" class="block text-sm font-medium text-gray-700">Received At</label>
        <input
          type="date"
          id="received"
          v-model="formData.received_at"
          class="mt-1 w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
        />
      </div>

      <button
        type="submit"
        class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition"
      >
        Create Income
      </button>
    </form>

    <p class="text-sm text-gray-600 text-center">Double-click a budget to delete</p>

    <!-- Income Cards -->
    <div class="space-y-4">
      <div v-if="incomeStore.incomes.length === 0" class="text-center text-gray-500">
        There are no incomes
      </div>

      <div
        v-for="income in incomeStore.incomes"
        :key="income.income_id"
        class="p-4 rounded-lg shadow-sm bg-white/40 backdrop-blur-md"
      >
        <!-- Income Header -->
        <div class="flex justify-between items-center mb-2">
          <div>
            <h3 class="text-lg font-semibold">{{ income.source }}</h3>
            <p class="text-sm text-gray-500">{{ formatDate(income.received_at) }}</p>
          </div>
          <p class="text-lg font-bold text-green-600">P{{ income.amount }}</p>
        </div>

        <!-- Budgets -->
        <div class="space-y-2">
          <p class="font-medium">Budgets Created</p>
          <div
            v-for="budget in income.budgets"
            :key="budget.budget_id"
            class="flex justify-between p-2 rounded-lg hover:bg-red-100 hover:cursor-pointer transition"
            @dblclick="showDialog = true"
          >
            <p>{{ budget.Category.name }}</p>
            <p class="font-semibold">{{ budget.amount }}</p>
            <ConfirmDialog
              v-model:show="showDialog"
              :message="'Delete this budget?'"
              @confirm="deleteItem(budget.budget_id)"
            />
          </div>
        </div>

        <hr class="my-3" />

        <!-- Remaining -->
        <div class="flex justify-between font-medium">
          <p>Remaining:</p>
          <p>{{ remainingAmount(income.amount, income.budgets) }}</p>
        </div>

        <!-- Budget Form Toggle -->
        <div class="mt-3">
          <button
            v-if="showForm"
            @click="openForm(income.income_id)"
            class="bg-white p-2 border rounded-lg hover:bg-gray-100 transition"
          >
            Create budget for this income
          </button>

          <form
            v-else
            @submit.prevent="submitBudgetForm"
            class="mt-3 space-y-3 bg-gray-50 p-4 rounded-lg border"
          >
            <div>
              <label for="category" class="block text-sm font-medium">Budget Category</label>
              <select
                id="category"
                class="border w-full px-2 py-1 rounded-lg"
                v-model="budgetForm.category_id"
              >
                <option disabled value="">Please select category</option>
                <option
                  v-for="category in categoryStore.categories"
                  :value="category.budget_category_id"
                  :key="category.id"
                >
                  {{ category.name }}
                </option>
              </select>
            </div>

            <div>
              <label for="amount" class="block text-sm font-medium">Amount</label>
              <input
                type="number"
                id="amount"
                class="border w-full px-2 py-1 rounded-lg"
                v-model="budgetForm.amount"
              />
            </div>

            <div>
              <label for="startDate" class="block text-sm font-medium">Start Date</label>
              <input
                type="date"
                id="startDate"
                class="border w-full px-2 py-1 rounded-lg"
                v-model="budgetForm.start_date"
              />
            </div>

            <div>
              <label for="endDate" class="block text-sm font-medium">End Date</label>
              <input
                type="date"
                id="endDate"
                class="border w-full px-2 py-1 rounded-lg"
                v-model="budgetForm.end_date"
              />
            </div>

            <div class="flex gap-2">
              <button type="submit" class="border px-3 py-1 rounded-lg hover:bg-green-100 transition">
                Create
              </button>
              <button
                @click="showForm = true"
                type="button"
                class="border px-3 py-1 rounded-lg hover:bg-gray-100 transition"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
