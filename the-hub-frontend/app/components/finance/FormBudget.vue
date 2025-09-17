<script setup lang="ts">
import { ref, computed } from "vue";
import { useBudgetStore, useCategoryStore } from '@/stores/finance'
import { useValidation } from '@/composables/useValidation'
import FormUI from '@/components/ui/FormUI.vue'

const budgetStore = useBudgetStore()
const categoryStore = useCategoryStore()
const { validateObject, schemas } = useValidation()

const showForm = ref(false)

// Fetch categories if not already loaded
if (categoryStore.categories.length === 0) {
  categoryStore.fetchCategory()
}

const categoryOptions = computed(() => {
  return categoryStore.categories.map(cat => ({
    value: cat.budget_category_id,
    label: cat.name
  }))
})

const fields = [
  {
    name: 'category_id',
    label: 'Category',
    type: 'select' as const,
    options: categoryOptions,
    placeholder: 'Select a category',
    required: true
  },
  {
    name: 'amount',
    label: 'Budget Amount',
    type: 'number' as const,
    placeholder: 'Enter budget amount',
    required: true,
    min: 0.01,
    step: 0.01
  },
  {
    name: 'start_date',
    label: 'Start Date',
    type: 'date' as const,
    required: true
  },
  {
    name: 'end_date',
    label: 'End Date',
    type: 'date' as const,
    required: true
  }
]

const initialData = {
  category_id: '',
  amount: '',
  start_date: '',
  end_date: ''
}

const submitForm = async (data: Record<string, any>) => {
  const payload = {
    category_id: data.category_id,
    amount: parseFloat(data.amount),
    start_date: data.start_date,
    end_date: data.end_date
  }

  const validation = validateObject(payload, schemas.budget.create)
  if (!validation.isValid) {
    // Validation errors will be shown by FormUI component
    return
  }

  try {
    await budgetStore.submitForm(payload)
    showForm.value = false // Close modal
  } catch (err) {
    // Error is already handled in the store
  }
}

const cancelForm = () => {
  showForm.value = false
}

const closeModal = () => {
  showForm.value = false
}

const openModal = () => {
  showForm.value = true
}
</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div v-if="!showForm" @click="openModal" class="fixed bottom-4 right-4 cursor-pointer z-40">
        <div class="bg-primary shadow-lg rounded-full p-4 hover:bg-primary/90 transition-all duration-200 hover:scale-105">
          <svg fill="currentColor" height="24px" width="24px" class="text-white" viewBox="0 0 24 24">
            <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
          </svg>
        </div>
      </div>
    </Teleport>
  </ClientOnly>

  <FormUI
    title="Create a Budget"
    :fields="fields"
    :initial-data="initialData"
    :show-form="showForm"
    submit-label="Create Budget"
    teleport-target="body"
    @submit="submitForm"
    @cancel="cancelForm"
    @close="closeModal"
  />
</template>