<script setup>

const props = defineProps({
  tabs: Array,
  modelValue: String,
})

const emit = defineEmits(['update:modelValue'])

// const activeTab = ref(props.defaultTab || props.tabs[0])
const activeTab = computed({
  get: () => props.modelValue || props.tabs[0],
  set: (val) => emit('update:modelValue', val)
})
</script>

<template>
  <div>
    <div class="pt-4 flex justify-center">
      <div
          class="inline-flex justify-self-center border-b p-4 text-center text-2xl rounded-2xl bg-orange-200/20 backdrop-blur-md border-white/10 shadow-lg">
        <button v-for="tab in tabs" :key="tab" @click="activeTab = tab"
          :class="['px-4 py-2', activeTab === tab ? 'border-b-2 border-blue-500' : 'text-gray-500']">
          {{ tab }}
        </button>
      </div>
    </div>

    <div class="mt-4">
      <slot :name="activeTab"></slot>
    </div>
  </div>
</template>
