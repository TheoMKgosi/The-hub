<script setup lang="ts">
interface Props {
  modelValue?: string | number;
  label?: string;
  type?: string;
  placeholder?: string;
  disabled?: boolean;
  inputId?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  inputId: () => `input-${Math.random().toString(36).substring(2, 9)}`,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
}>();

const updateValue = (event: Event) => {
  const target = event.target as HTMLInputElement;
  emit('update:modelValue', target.value);
};
</script>

<template>
  <div class="base-input flex flex-col gap-1">
    <label 
      v-if="label" 
      :for="inputId" 
      class="base-input__label text-sm font-medium text-gray-700"
    >
      {{ label }}
    </label>
    
    <input
      :id="inputId"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="updateValue" 
      class="base-input__field border rounded px-3 py-2 outline-none focus:ring-2 focus:ring-blue-500 transition-all flex-1"
      :class="{ 'bg-gray-100 cursor-not-allowed opacity-50': disabled }"
    />
  </div>
</template>
