<script setup lang="ts">
import * as v from '@valibot/valibot';
const goalStore = useGoalStore()

const value = ref(3)

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
  } catch (err) {
    // Error is already handled in the store
  }
}
</script>

<template>
  <UModal>
    <UButton label="Add a goal"/>

    <template #body>
      <UForm>
        <UFormField label="Title" required>
          <UInput />
        </UFormField>
        <UFormField label="Description">
          <UTextarea />
        </UFormField>
        <UFormField label="Priority (Higher is more important)">
          <UInputNumber v-model="value" :min="0" :max="5" />
        </UFormField>
        <UFormField label="Category">
          <UInput />
        </UFormField>
        <UFormField label="Colour">
          <UInput />
        </UFormField>
      </UForm>
    </template>
  </UModal>
</template>
