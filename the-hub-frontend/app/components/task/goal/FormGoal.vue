<script setup lang="ts">
import type { FormSubmitEvent } from '@nuxt/ui';
import * as v from 'valibot';
const goalStore = useGoalStore()

const value = ref(3)

const schema = v.object({
  title: v.pipe(v.string()),
  description: v.pipe(v.string()),
  due_date: v.pipe(v.string()),
  priority: v.pipe(v.number()),
  category: v.pipe(v.string()),
  color: v.pipe(v.string())
})

type Schema = v.InferOutput<typeof schema>

const payload = reactive({
  title: '',
  description: '',
  due_date: null,
  priority: value,
  category: '',
  color: ''
})

const submitForm = async (event: FormSubmitEvent<Schema>) => {
  try {
    await goalStore.createGoal(event.data)
  } catch (err) {
    // Error is already handled in the store
  }
}
</script>

<template>
  <UModal title="Add Goal">
    <UButton label="Add a goal" />

    <template #body>
      <UForm @submit="submitForm" :state="payload" class="space-y-3">
        <UFormField label="Title" required>
          <UInput v-model="payload.title" class="w-full" />
        </UFormField>
        <UFormField label="Description" class="w-full">
          <UTextarea v-model="payload.description"/>
        </UFormField>
        <UFormField label="Priority (Higher is more important)">
          <UInputNumber v-model="payload.priority" :min="0" :max="5" />
        </UFormField>
        <UFormField label="Category" class="w-full">
          <UInput v-model="payload.category"/>
        </UFormField>
        <UFormField label="Colour" class="w-full">
          <UInput v-model="payload.color"/>
        </UFormField>
        <UButton label="Add Goal" type="submit" />
      </UForm>
    </template>
  </UModal>
</template>
