<script setup lang="ts">
const goalStore = useGoalStore()
const { schemas } = useValidation()

const showForm = ref(false) 

const fields = [
  {
    name: 'title',
    label: 'Title',
    type: 'text' as const,
    placeholder: 'Goal title',
    required: true
  },
  {
    name: 'description',
    label: 'Description',
    type: 'textarea' as const,
    placeholder: 'Optional description',
    rows: 3
  },
  {
    name: 'due_date',
    label: 'Due Date',
    type: 'date' as const
  },
  {
    name: 'priority',
    label: 'Priority',
    type: 'select' as const,
    options: [
      { value: null, label: 'No priority' },
      { value: 1, label: '1 - Low' },
      { value: 2, label: '2' },
      { value: 3, label: '3 - Medium' },
      { value: 4, label: '4' },
      { value: 5, label: '5 - High' }
    ]
  },
  {
    name: 'category',
    label: 'Category',
    type: 'text' as const,
    placeholder: 'e.g., Work, Personal, Health'
  },
  {
    name: 'color',
    label: 'Color',
    type: 'color' as const
  }
]

const initialData = {
  title: '',
  description: '',
  due_date: '',
  priority: null,
  category: '',
  color: '#3B82F6',
}

const submitForm = async (data: Record<string, any>) => {
  const payload = {
    title: data.title.trim(),
    description: data.description.trim(),
    due_date: data.due_date ? new Date(data.due_date).toISOString() : undefined,
    priority: data.priority || undefined,
    category: data.category.trim() || undefined,
    color: data.color,
  }

  try {
    await goalStore.createGoal(payload)
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
    title="Create a Goal"
    :fields="fields"
    :initial-data="initialData"
    :show-form="showForm"
    :validation-schema="schemas.goal.create"
    submit-label="Create Goal"
    teleport-target="body"
    @submit="submitForm"
    @cancel="cancelForm"
    @close="closeModal"
  />
</template>
