<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useCategoryStore, useBudgetStore } from "@/stores/finance";

const formData = reactive({
  category_id: '',
  amount: 0,
  start_date: '',
  end_date: '',
})

const categoryStore = useCategoryStore()
const budgetStore = useBudgetStore()

const submitForm = async () => {
  const dataToSend = { ...formData }
  budgetStore.submitForm(dataToSend)
  Object.assign(formData, {
    category_id: '',
    amount: 0,
    start_date: '',
    end_date: '',
  })
}

onMounted(() => {
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory()
  }
  if (budgetStore.budgets.length === 0) {
    budgetStore.fetchBudget()
  }
})

const formatDate = (date) => new Date(date).toLocaleDateString()

</script>

<template>
  <div class="p-4">
    <form @submit.prevent="submitForm" class="space-y-2 mb-6 max-w-2xl">
      <select v-model="formData.category_id" class="border p-2 w-full">
        <option disabled value="">Select category</option>
        <option v-for="cat in categoryStore.categories" :key="cat.budget_category_id" :value="cat.budget_category_id">
          {{ cat.name }}
        </option>
      </select>

      <input v-model.number="formData.amount" type="number" placeholder="Amount" class="border p-2 w-full" />
      <input v-model="formData.start_date" type="date" class="border p-2 w-full" />
      <input v-model="formData.end_date" type="date" class="border p-2 w-full" />

      <button class="bg-blue-600 text-white px-4 py-2">Create Budget</button>
    </form>

    <div v-if="budgetStore.budgets.length === 0" class="text-gray-500">No budgets yet.</div>
    <div v-else class="space-y-6">
      <!--
      <div v-for="(items, month) in budgetStore.budgets" :key="month">
        <div class="flex justify-between items-center mb-2">
          <h2 class="text-xl font-bold">{{ month }}</h2>
          <span class="text-gray-700 font-medium">Total: ${{ getMonthTotal(items).toFixed(2) }}</span>
        </div>
-->

      <div class="grid gap-4">
        <div v-for="budget in budgetStore.budgets" :key="budget.budget_id" class="border p-4 rounded shadow">
          <h3 class="font-bold">{{ budget.Category.name }}</h3>
          <p>Amount: ${{ budget.amount.toFixed(2) }}</p>
          <p class="text-sm text-gray-600">
            {{ formatDate(budget.start_date) }} → {{ formatDate(budget.end_date) }}
          </p>
        </div>
      </div>
    </div>
  </div>

</template>
