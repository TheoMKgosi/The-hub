<script setup>
const props = defineProps({
  show: Boolean,
  message: {
    type: String,
    default: 'Are you sure?',
  },
})
const emit = defineEmits(['update:show', 'confirm'])

const cancel = () => emit('update:show', false)
const confirm = () => {
  emit('confirm')
  emit('update:show', false)
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex justify-center items-center z-50">
    <div class="bg-surface-light dark:bg-surface-dark p-5 rounded-lg shadow-lg w-72">
      <p class="mb-4 text-text-light dark:text-text-dark">{{ message }}</p>
      <div class="flex justify-end gap-3">
        <UiBaseButton @click="cancel" variant="default" size="sm">Cancel</UiBaseButton>
        <UiBaseButton @click="confirm" variant="danger" size="sm">Confirm</UiBaseButton>
      </div>
    </div>
  </div>
</template>


