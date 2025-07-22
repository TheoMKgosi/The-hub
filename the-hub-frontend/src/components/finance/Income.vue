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
  <div>
    <!-- Filters -->
    <!-- Filters + Search -->
    <div class="shadow-sm p-3 bg-white/20 backdrop-blur-md rounded-lg mt-2">
      <div class="flex flex-wrap gap-2 items-center mb-2">
        <input v-model="searchQuery" placeholder="Search tasks..." class="flex-grow shadow-sm  bg-gradient-to-r from-gray-50 to-gray-100 px-3 py-2 rounded
          w-full sm:w-0" />
      </div>
    </div>

    <!-- Adding income for Category -->
    <div>
      <form @submit.prevent="submitForm" class="space-y-4 p-4 max-w-md mx-auto bg-white rounded-2xl shadow">
        <h2>Income Add Form</h2>
        <div>
          <label for="source" class="block text-sm font-medium text-gray-700">Income Source</label>
          <input type="text" id="source" v-model="formData.source"
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" />
        </div>

        <div>
          <label for="amount" class="block text-sm font-medium text-gray-700">Amount</label>
          <input type="number" id="amount" v-model="formData.amount"
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" />
        </div>

        <div>
          <label for="received" class="block text-sm font-medium text-gray-700">Received At</label>
          <input type="date" id="received" v-model="formData.received_at"
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500" />
        </div>

        <button type="submit" class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition">
          Create Income
        </button>
      </form>

    </div>

    <p>Double click budget to delete</p>

    <!-- Income cards -->
    <div class="mx-4 px-6 py-2 backdrop-blur-md bg-white/30 shadow-sm">
      <div v-if="incomeStore.incomes.length === 0">The are no incomes</div>
      <div v-for="income in incomeStore.incomes">
        <div class="flex justify-between">
          <div>
            <h2>{{ income.source }}</h2>
            <p>{{ formatDate(income.received_at) }}</p>
          </div>
          <p>P{{ income.amount }}</p>
        </div>

        <div>
          <p>Budgets Created</p>
          <div v-for="budget in income.budgets" class="flex justify-between hover:font-extrabold" @dblclick="showDialog
            = true">
            <p>{{ budget.Category.name }}: </p>
            <p>{{ budget.amount }} </p>
            <ConfirmDialog v-model:show="showDialog" :message="'Delete this budget?'"
              @confirm="deleteItem(budget.budget_id)" />
          </div>
          <hr class="p-2">

          <div>
            <div class="flex justify-between">
              <p>Remaining:</p>
              <p>{{ remainingAmount(income.amount, income.budgets) }}</p>
            </div>
            <div v-if="showForm">
              <button @click="openForm(income.income_id)" class="bg-white p-2 border rounded-lg">Create budget for this
                income</button>
            </div>
            <div v-if="!showForm">
              <form @submit.prevent="submitBudgetForm">
                <label for="category">Budget Category</label>
                <select id="category" class="border" v-model="budgetForm.category_id">
                  <option disabled value="">Please select category</option>
                  <option :value="category.budget_category_id" v-for="category in categoryStore.categories"
                    :key="category.id">
                    {{ category.name }}
                  </option>
                </select>

                <label for="amount">Amount</label>
                <input type="number" id="amount" class="border" v-model="budgetForm.amount">

                <label for="startDate">Start Date</label>
                <input type="date" id="startDate" class="border" v-model="budgetForm.start_date">

                <label for="endDate">End Date</label>
                <input type="date" id="endDate" class="border" v-model="budgetForm.end_date">

                <button type="submit" class="border">Create</button>
                <button @click="showForm = true" type="button" class="border">Cancel</button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
