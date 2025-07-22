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
  <div v-if="show" class="fixed inset-0 bg-black/50 flex justify-center items-center z-50">
    <div class="bg-white p-5 rounded shadow w-72">
      <p class="mb-4">{{ message }}</p>
      <div class="flex justify-end gap-3">
        <button @click="cancel" class="px-3 py-1 border rounded">Cancel</button>
        <button @click="confirm" class="bg-red-600 text-white px-3 py-1 rounded">Confirm</button>
      </div>
    </div>
  </div>
</template>


