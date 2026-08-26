<script setup lang="ts">
const incomeStore = useIncomeStore();
const categoryStore = useCategoryStore();
const budgetStore = useBudgetStore();

const activeIncomeId = ref<number | null>(null);
const showDialog = ref(false);
const showIncomeModal = ref(false);
const searchQuery = ref("");

const budgetID = ref(0);
const incomeID = ref(0);

// Edit budget state
const isEditingBudget = ref(false);
const editingBudgetId = ref("");

const formData = reactive({
  source: "",
  amount: 0,
  received_at: null,
});

const budgetForm = reactive({
  income_id: 0,
  category_id: "",
  amount: 0,
  start_date: null,
  end_date: null,
});

const editBudgetForm = reactive({
  budget_id: "",
  category_id: "",
  amount: 0,
  start_date: "",
  end_date: "",
});

const budgetFormFields = computed(() => [
  {
    name: "amount",
    label: "Budget Amount",
    type: "number",
    placeholder: "0.00",
    min: 0.01,
    step: 0.01,
    required: true,
  },
  {
    name: "category_id",
    label: "Category",
    type: "combobox",
    categories: categoryStore.categories,
    placeholder: "Select or create category...",
    allowCreate: true,
    required: true,
  },
  {
    name: "start_date",
    label: "Start Date",
    type: "date",
    required: true,
  },
  {
    name: "end_date",
    label: "End Date",
    type: "date",
    required: true,
  },
]);

const filteredIncome = computed(() => {
  let result = incomeStore.incomes;

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter((income) =>
      income.source.toLowerCase().includes(query),
    );
  }

  return result;
});

const getBudgetAnalytics = (budgetId: string) => {
  return budgetStore.analytics.find(
    (analytic) => analytic.budget_id === budgetId,
  );
};

const deleteItem = (id: number, incomeID: number) => {
  budgetStore.deleteBudget(id, incomeID);
};

const submitForm = async () => {
  const dataToSend = { ...formData };
  incomeStore.submitForm(dataToSend);
  Object.assign(formData, {
    source: "",
    amount: 0,
    received_at: null,
  });
  showIncomeModal.value = true;
};

const submitBudgetForm = async () => {
  // Validation
  if (!budgetForm.category_id) {
    alert("Please select a category");
    return;
  }

  if (!budgetForm.amount || budgetForm.amount <= 0) {
    alert("Please enter a valid amount greater than 0");
    return;
  }

  if (!budgetForm.start_date) {
    alert("Please select a start date");
    return;
  }

  if (!budgetForm.end_date) {
    alert("Please select an end date");
    return;
  }

  const startDate = new Date(budgetForm.start_date);
  const endDate = new Date(budgetForm.end_date);

  if (endDate <= startDate) {
    alert("End date must be after start date");
    return;
  }

  const dataToSend = { ...budgetForm };
  budgetStore.submitForm(dataToSend);
};

onMounted(async () => {
  if (incomeStore.incomes.length === 0) {
    incomeStore.fetchIncomes();
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory();
  }
  // if (budgetStore.analytics.length === 0) {
  //   await budgetStore.fetchBudgetAnalytics("current");
  // }
});

const formatDate = (date) => new Date(date).toLocaleDateString();
const openForm = (id: number) => {
  activeIncomeId.value = id;
  budgetForm.income_id = id;
};

const closeForm = () => {
  activeIncomeId.value = null;
  Object.assign(budgetForm, {
    category_id: "",
    amount: 0,
    start_date: null,
    end_date: null,
  });
};

const remainingAmount = (amount, budgets) => {
  let remaining = amount;
  for (const budget of budgets) {
    remaining -= budget.amount;
  }
  return remaining;
};

const handleBudgetCategorySelect = (category) => {
  budgetForm.category_id = category.budget_category_id;
};

const handleBudgetCategoryCreate = async (fieldName, categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName });
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(
      (cat) => cat.name.toLowerCase() === categoryName.toLowerCase(),
    );
    if (newCategory) {
      budgetForm.category_id = newCategory.budget_category_id;
    }
  } catch (error) {
    console.error("Failed to create category:", error);
  }
};

// Edit budget methods
const startEditBudget = (budget) => {
  isEditingBudget.value = true;
  editingBudgetId.value = budget.budget_id;
  Object.assign(editBudgetForm, {
    budget_id: budget.budget_id,
    category_id: budget.category_id,
    amount: budget.amount,
    start_date: budget.start_date,
    end_date: budget.end_date,
  });
};

const cancelEditBudget = () => {
  isEditingBudget.value = false;
  editingBudgetId.value = "";
  Object.assign(editBudgetForm, {
    budget_id: "",
    category_id: "",
    amount: 0,
    start_date: "",
    end_date: "",
  });
};

const submitEditBudget = async () => {
  // Validation
  if (!editBudgetForm.category_id) {
    alert("Please select a category");
    return;
  }

  if (!editBudgetForm.amount || editBudgetForm.amount <= 0) {
    alert("Please enter a valid amount greater than 0");
    return;
  }

  if (!editBudgetForm.start_date) {
    alert("Please select a start date");
    return;
  }

  if (!editBudgetForm.end_date) {
    alert("Please select an end date");
    return;
  }

  const startDate = new Date(editBudgetForm.start_date);
  const endDate = new Date(editBudgetForm.end_date);

  if (endDate <= startDate) {
    alert("End date must be after start date");
    return;
  }

  try {
    await budgetStore.editBudget(editBudgetForm);
    cancelEditBudget();
  } catch (error) {
    alert("Failed to update budget");
  }
};

const handleEditBudgetCategorySelect = (category) => {
  editBudgetForm.category_id = category.budget_category_id;
};

const handleEditBudgetCategoryCreate = async (categoryName) => {
  try {
    await categoryStore.submitForm({ name: categoryName });
    // The new category should now be available in the store
    // Find it and set it as selected
    const newCategory = categoryStore.categories.find(
      (cat) => cat.name.toLowerCase() === categoryName.toLowerCase(),
    );
    if (newCategory) {
      editBudgetForm.category_id = newCategory.budget_category_id;
    }
  } catch (error) {
    console.error("Failed to create category:", error);
  }
};
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">
      Income Management
    </h2>

    <!-- Action Buttons -->
    <UModal title="Add Income" v-model:open="showIncomeModal">
      <UButton label="Add Income" variant="outline" />

      <template #body>
        <UForm @submit="submitForm" class="space-y-4">
          <UFormField label="Income Source">
            <UInput v-model="formData.source" placeholder="e.g., Salary, Freelance, Business" class="w-full" />
          </UFormField>

          <UFormField label="Amount">
            <UInputNumber v-model="formData.amount" :format-options="{ style: 'currency', currency: 'BWP' }" />
          </UFormField>

          <div>
            <label for="received" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Received
              Date</label>
            <input type="date" id="received" v-model="formData.received_at"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
          </div>

          <!-- Modal Footer -->
          <div
            class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
            <UButton label="Cancel" color="neutral" variant="outline" @click="showIncomeModal = false" />
            <UButton label="Create Income" color="neutral" variant="soft" type="submit" />
          </div>
        </Uform>
      </template>
    </UModal>

    <!-- Filters + Search -->
    <div
      class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search income sources..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <p class="text-sm text-text-light dark:text-text-dark/60 text-center">
      Double-click a budget to delete it
    </p>

    <!-- Income Cards -->
    <div class="space-y-4">
      <div v-if="incomeStore.incomes.length === 0" class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p class="text-lg mb-2">No income sources added yet</p>
        <p class="text-sm">
          Create your first income source above to get started
        </p>
      </div>

      <div v-for="income in filteredIncome" :key="income.income_id"
        class="p-6 rounded-lg shadow-md bg-surface-light dark:bg-surface-dark border border-surface-light dark:border-surface-dark hover:shadow-lg transition-shadow duration-200">
        <!-- Income Header -->
        <div class="flex justify-between items-start mb-4">
          <div class="flex-1">
            <h3 class="text-xl font-semibold text-text-light dark:text-text-dark mb-1">
              {{ income.source }}
            </h3>
            <p class="text-sm text-text-light dark:text-text-dark/60">
              {{ formatDate(income.received_at) }}
            </p>
          </div>
          <p class="text-xl font-bold text-success">
            ${{ income.amount.toFixed(2) }}
          </p>
        </div>

        <!-- Budgets -->
        <div class="space-y-3">
          <h4 class="font-medium text-text-light dark:text-text-dark">
            Allocated Budgets
          </h4>
          <div v-if="income.budgets.length === 0" class="text-sm text-text-light dark:text-text-dark/60 italic">
            No budgets created yet
          </div>
          <div v-else v-for="budget in income.budgets" :key="budget.budget_id"
            class="p-3 rounded-md bg-surface-light/50 dark:bg-surface-dark/50 border border-surface-light dark:border-surface-dark hover:bg-primary/10 dark:hover:bg-primary/20 hover:border-primary/50 dark:hover:border-primary/30 transition-colors">
            <!-- Budget Header with Actions -->
            <div class="flex justify-between items-start mb-2">
              <div class="flex-1">
                <p class="font-medium text-text-light dark:text-text-dark">
                  {{ budget.Category.name }}
                </p>
                <p class="text-xs text-text-light dark:text-text-dark/60">
                  {{ formatDate(budget.start_date) }} -
                  {{ formatDate(budget.end_date) }}
                </p>
              </div>
              <div class="flex items-center gap-2">
                <p class="font-semibold text-text-light dark:text-text-dark">
                  P{{ budget.amount.toFixed(2) }}
                </p>
                <div class="flex gap-1">
                  <button @click="startEditBudget(budget)"
                    class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 p-1 rounded"
                    title="Edit budget">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z">
                      </path>
                    </svg>
                  </button>
                  <button @click="
                    showDialog = true;
                  budgetID = budget.budget_id;
                  incomeID = income.income_id;
                  " class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 p-1 rounded"
                    title="Delete budget">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                      </path>
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            <!-- Budget Performance -->
            <!--
            <div v-if="getBudgetAnalytics(budget.budget_id)" class="space-y-2">
              <div class="flex justify-between text-xs text-text-light dark:text-text-dark/60">
                <span>Spent: ${{
                  getBudgetAnalytics(budget.budget_id).spent_amount.toFixed(2)
                }}</span>
                <span>Remaining: ${{
                  getBudgetAnalytics(
                    budget.budget_id,
                  ).remaining_amount.toFixed(2)
                }}</span>
              </div>
              <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div :class="[
                  'h-2 rounded-full transition-all duration-300',
                  getBudgetAnalytics(budget.budget_id).status ===
                    'over_budget'
                    ? 'bg-red-500'
                    : getBudgetAnalytics(budget.budget_id).status ===
                      'warning'
                      ? 'bg-yellow-500'
                      : getBudgetAnalytics(budget.budget_id).status ===
                        'caution'
                        ? 'bg-orange-500'
                        : 'bg-green-500',
                ]" :style="{
                    width:
                      Math.min(
                        getBudgetAnalytics(budget.budget_id).utilization_rate,
                        100,
                      ) + '%',
                  }"></div>
              </div>
              <div class="flex justify-between items-center text-xs">
                <span class="text-text-light dark:text-text-dark/60">
                  {{
                    getBudgetAnalytics(
                      budget.budget_id,
                    ).utilization_rate.toFixed(1)
                  }}% used
                </span>
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  getBudgetAnalytics(budget.budget_id).status === 'on_track'
                    ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100'
                    : getBudgetAnalytics(budget.budget_id).status ===
                      'caution'
                      ? 'bg-orange-100 text-orange-800 dark:bg-orange-800 dark:text-orange-100'
                      : getBudgetAnalytics(budget.budget_id).status ===
                        'warning'
                        ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100'
                        : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100',
                ]">
                  {{
                    getBudgetAnalytics(budget.budget_id)
                      .status.replace("_", " ")
                      .toUpperCase()
                  }}
                </span>
              </div>
            </div>
            -->

            <!-- Edit Form -->
            <div v-if="isEditingBudget && editingBudgetId === budget.budget_id"
              class="mt-3 p-3 bg-surface-light dark:bg-surface-dark rounded-md border border-surface-light dark:border-surface-dark">
              <form @submit.prevent="submitEditBudget" class="space-y-3">
                <div>
                  <label class="block text-xs font-medium text-text-light dark:text-text-dark mb-1">Category</label>
                  <UiBaseComboBox :model-value="editBudgetForm.category_id" :categories="categoryStore.categories"
                    placeholder="Select category..." @select="handleEditBudgetCategorySelect"
                    @create="handleEditBudgetCategoryCreate" />
                </div>

                <div>
                  <label class="block text-xs font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
                  <input v-model.number="editBudgetForm.amount" type="number" placeholder="0.00" step="0.01" min="0.01"
                    required
                    class="w-full px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded focus:outline-none focus:ring-1 focus:ring-primary" />
                </div>

                <div class="grid grid-cols-2 gap-2">
                  <div>
                    <label class="block text-xs font-medium text-text-light dark:text-text-dark mb-1">Start Date</label>
                    <input v-model="editBudgetForm.start_date" type="date" required
                      class="w-full px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded focus:outline-none focus:ring-1 focus:ring-primary" />
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-text-light dark:text-text-dark mb-1">End Date</label>
                    <input v-model="editBudgetForm.end_date" type="date" required
                      class="w-full px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded focus:outline-none focus:ring-1 focus:ring-primary" />
                  </div>
                </div>

                <div class="flex justify-end gap-2 pt-2">
                  <UButton label="Cancel" type="button" @click="cancelEditBudget" variant="outline" color="neutral"
                    size="xs" />
                  <UButton :label="budgetStore.updating ? ' Updating...' : 'Update'" type="submit" variant="soft"
                    color="neutral" size="xs" :disabled="budgetStore.updating" />
                </div>
              </form>
            </div>

            <UiConfirmDialog v-model:show="showDialog" :message="`Delete budget for ${budget.Category.name}?`"
              @confirm="deleteItem(budgetID, incomeID)" />
          </div>
        </div>

        <hr class="my-4 border-surface-light dark:border-surface-dark" />

        <!-- Remaining -->
        <div class="flex justify-between items-center font-medium">
          <p class="text-text-light dark:text-text-dark">Remaining:</p>
          <p class="text-lg" :class="remainingAmount(income.amount, income.budgets) >= 0
            ? 'text-success'
            : 'text-red-500 dark:text-red-400'
            ">
            ${{ remainingAmount(income.amount, income.budgets).toFixed(2) }}
          </p>
        </div>

        <!-- Budget Form Toggle -->
        <div v-if="activeIncomeId !== income.income_id" class="mt-4">
          <UButton label="Create Budget for This Income" @click="openForm(income.income_id)" variant="outline" size="sm"
            class="w-full" />
        </div>

        <ClientOnly>
          <Teleport to="body">
            <Transition name="fade-scale">
              <div v-if="activeIncomeId === income.income_id"
                class="fixed inset-0 flex items-center justify-center bg-black/50 dark:bg-black/70 z-50 p-4">
                <UForm @submit="submitBudgetForm" class="space-y-2">
                  <UFormField label="Amount">
                    <UInput v-model="budgetForm.amount"/>
                  </UFormField>
                  <UFormField>
                    <USelectMenu :items="categoryStore.categories"/>
                  </UFormField>
                  <input v-model="budgetForm.start_date" type="date"/>
                  <input v-model="budgetForm.end_date" type="date"/>
                  <UButton label="Create Budget" type="submit"/>
                </UForm>
              </div>
            </Transition>
          </Teleport>
        </ClientOnly>
      </div>
    </div>
  </div>
</template>
