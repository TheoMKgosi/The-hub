<script setup lang="ts">
import ConfirmDialog from "@/components/ConfirmDialog.vue";
import ComboBox from "@/components/ui/ComboBox.vue";
import FormUI from "@/components/ui/FormUI.vue";

const incomeStore = useIncomeStore();
const categoryStore = useCategoryStore();
const budgetStore = useBudgetStore();

const activeIncomeId = ref<number | null>(null);
const showDialog = ref(false);
const showIncomeModal = ref(true);
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
  // Form reset is now handled automatically by FormUI component
};

onMounted(async () => {
  if (incomeStore.incomes.length === 0) {
    incomeStore.fetchIncomes();
  }
  if (categoryStore.categories.length === 0) {
    categoryStore.fetchCategory();
  }
  if (budgetStore.analytics.length === 0) {
    await budgetStore.fetchBudgetAnalytics("current");
  }
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

const handleBudgetFormSubmit = async (formData) => {
  // Handle category creation if needed
  if (formData.category_id && typeof formData.category_id === 'string') {
    // Check if this is a newly created category (not in the existing categories)
    const existingCategory = categoryStore.categories.find(
      cat => cat.budget_category_id === formData.category_id
    );
    if (!existingCategory) {
      // This was a newly created category, but FormUI doesn't handle creation directly
      // The category should have been created via the combobox @create event
      // We can assume it's already in the store
    }
  }

  // Call the existing submit logic
  await submitBudgetForm();
};
</script>

<template>
  <div class="space-y-6 p-4 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold text-text-light dark:text-text-dark">
      Income Management
    </h2>

    <!-- Filters + Search -->
    <div
      class="shadow-sm p-4 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg border border-surface-light/10 dark:border-surface-dark/10">
      <input v-model="searchQuery" placeholder="Search income sources..."
        class="w-full px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary" />
    </div>

    <!-- Floating Action Button -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="showIncomeModal" @click="showIncomeModal = false" class="fixed bottom-4 right-4 cursor-pointer z-40">
          <div
            class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
            <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
            </svg>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

    <!-- Income Modal -->
    <ClientOnly>
      <Teleport to="body">
        <div v-if="!showIncomeModal" @click="showIncomeModal = true"
          class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50">
          <div
            class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-md max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light dark:border-surface-dark"
            @click.stop>
            <!-- Modal Header -->
            <div class="flex items-center justify-between p-6 border-b border-surface-light dark:border-surface-dark">
              <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">
                Add New Income
              </h2>
              <UiButton @click="showIncomeModal = true" variant="default" size="sm" class="p-2">
                ×
              </UiButton>
            </div>

            <!-- Modal Body -->
            <div class="p-6">
              <form @submit.prevent="submitForm" class="space-y-4">
                <div>
                  <label for="source" class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Income
                    Source</label>
                  <input type="text" id="source" v-model="formData.source"
                    placeholder="e.g., Salary, Freelance, Business"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <div>
                  <label for="amount"
                    class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Amount</label>
                  <input type="number" id="amount" v-model="formData.amount" placeholder="0.00" step="0.01" min="0"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <div>
                  <label for="received"
                    class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Received Date</label>
                  <input type="date" id="received" v-model="formData.received_at"
                    class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary" />
                </div>

                <!-- Modal Footer -->
                <div
                  class="flex flex-col-reverse sm:flex-row gap-3 pt-6 border-t border-surface-light dark:border-surface-dark">
                  <UiButton type="button" @click="showIncomeModal = true" variant="default" size="md"
                    class="w-full sm:w-auto">
                    Cancel
                  </UiButton>
                  <UiButton type="submit" variant="primary" size="md" class="w-full sm:w-auto">
                    Create Income
                  </UiButton>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Teleport>
    </ClientOnly>

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

            <!-- Edit Form -->
            <div v-if="isEditingBudget && editingBudgetId === budget.budget_id"
              class="mt-3 p-3 bg-surface-light dark:bg-surface-dark rounded-md border border-surface-light dark:border-surface-dark">
              <form @submit.prevent="submitEditBudget" class="space-y-3">
                <div>
                  <label class="block text-xs font-medium text-text-light dark:text-text-dark mb-1">Category</label>
                  <ComboBox :model-value="editBudgetForm.category_id" :categories="categoryStore.categories"
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
                  <UiButton type="button" @click="cancelEditBudget" variant="default" size="xs">
                    Cancel
                  </UiButton>
                  <UiButton type="submit" variant="primary" size="xs" :disabled="budgetStore.updating">
                    {{ budgetStore.updating ? "Updating..." : "Update" }}
                  </UiButton>
                </div>
              </form>
            </div>

            <ConfirmDialog v-model:show="showDialog" :message="`Delete budget for ${budget.Category.name}?`"
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
          <UiButton @click="openForm(income.income_id)" variant="default" size="sm" class="w-full">
            Create Budget for This Income
          </UiButton>
        </div>

        <ClientOnly>
          <Teleport to="body">
            <Transition name="fade-scale">
              <div v-if="activeIncomeId === income.income_id"
                class="fixed inset-0 flex items-center justify-center bg-black/50 dark:bg-black/70 z-50 p-4">
                <FormUI
                  title="Create Budget"
                  :fields="budgetFormFields"
                  submit-label="Create Budget"
                  cancel-label="Cancel"
                  :initial-data="budgetForm"
                  size="md"
                  @submit="handleBudgetFormSubmit"
                  @cancel="closeForm"
                  @combobox-create="handleBudgetCategoryCreate"
                  class="w-full max-w-md"
                />
              </div>
            </Transition>
          </Teleport>
        </ClientOnly>
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
