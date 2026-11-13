<script setup lang="ts">
import { useValidation } from '@/composables/useValidation'

const activityStore = useActivityLogStore()
const { validateObject } = useValidation()

const formData = reactive({
  title: '',
  notes: '',
  category: '',
  tagsInput: '',
  activity_date: new Date().toISOString().slice(0, 10),
  duration_min: 0 as number | null
})

const validationErrors = ref<Record<string, string>>({})
const validationSchema = {
  title: { required: true, maxLength: 255 }
}

const isSubmitting = ref(false)

const parsedTags = computed(() =>
  formData.tagsInput.split(',').map(tag => tag.trim()).filter(Boolean)
)

const submitActivity = async () => {
  const validation = validateObject(formData, validationSchema)
  validationErrors.value = validation.errors
  if (!validation.isValid) return

  isSubmitting.value = true
  try {
    await activityStore.createEntry({
      title: formData.title.trim(),
      notes: formData.notes.trim(),
      category: formData.category.trim(),
      tags: parsedTags.value,
      activity_date: formData.activity_date,
      duration_min: formData.duration_min || 0
    })
    // Reset the free-text fields but keep category/tags/date for quick consecutive logging
    formData.title = ''
    formData.notes = ''
    formData.duration_min = 0
  } finally {
    isSubmitting.value = false
  }
}

const applySuggestedCategory = (category: string) => {
  formData.category = category
}
</script>

<template>
  <UForm :state="formData" class="space-y-3" @submit="submitActivity">
    <UFormField label="What did you do?" required :error="validationErrors.title">
      <UInput
        v-model="formData.title"
        class="w-full"
        placeholder="e.g. Worked on thesis chapter 3"
      />
    </UFormField>

    <div class="grid gap-3 sm:grid-cols-3">
      <UFormField label="Category">
        <UInput
          v-model="formData.category"
          class="w-full"
          placeholder="e.g. Work"
        />
        <div v-if="activityStore.categorySuggestions.length" class="mt-2 flex flex-wrap gap-1">
          <UBadge
            v-for="suggestion in activityStore.categorySuggestions.slice(0, 5)"
            :key="suggestion"
            variant="subtle"
            class="cursor-pointer"
            @click="applySuggestedCategory(suggestion)"
          >
            {{ suggestion }}
          </UBadge>
        </div>
      </UFormField>

      <UFormField label="Tags (comma separated)">
        <UInput
          v-model="formData.tagsInput"
          class="w-full"
          placeholder="focus, writing"
        />
        <div v-if="activityStore.tagSuggestions.length" class="mt-2 flex flex-wrap gap-1">
          <UBadge
            v-for="suggestion in activityStore.tagSuggestions.slice(0, 5)"
            :key="suggestion"
            variant="subtle"
            color="neutral"
            class="cursor-pointer"
            @click="formData.tagsInput = formData.tagsInput ? `${formData.tagsInput}, ${suggestion}` : suggestion"
          >
            {{ suggestion }}
          </UBadge>
        </div>
      </UFormField>

      <UFormField label="Date">
        <UInput
          v-model="formData.activity_date"
          type="date"
          class="w-full"
        />
      </UFormField>
    </div>

    <div class="grid gap-3 sm:grid-cols-3">
      <UFormField label="Duration (minutes)">
        <UInputNumber
          v-model="formData.duration_min"
          :min="0"
          class="w-full"
          placeholder="0"
        />
      </UFormField>

      <UFormField label="Notes (optional)" class="sm:col-span-2">
        <UTextarea
          v-model="formData.notes"
          class="w-full"
          :rows="1"
          placeholder="Any details worth remembering later…"
        />
      </UFormField>
    </div>

    <UButton
      type="submit"
      label="Log activity"
      icon="i-lucide-plus"
      :loading="isSubmitting"
    />
  </UForm>
</template>
