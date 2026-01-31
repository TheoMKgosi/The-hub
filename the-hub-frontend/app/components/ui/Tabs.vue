<script setup lang="ts">
/**
 * Tabs Component
 *
 * Usage:
 * <Tabs v-model="activeTab" :tabs="['Tab 1', 'Tab 2', 'Tab 3']">
 *   <template #Tab 1>
 *     Content for Tab 1
 *   </template>
 *   <template #Tab 2>
 *     Content for Tab 2
 *   </template>
 * </Tabs>
 */

interface TabsProps {
  tabs: string[]
  modelValue?: string
}

const props = withDefaults(defineProps<TabsProps>(), {
  modelValue: undefined
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const activeTab = computed({
  get: () => props.modelValue || props.tabs[0],
  set: (val) => emit('update:modelValue', val)
})

const handleTabClick = (tab: string) => {
  activeTab.value = tab
}
</script>

<template>
  <div>
    <div class="pt-4 flex justify-center">
      <div
        class="inline-flex justify-self-center border-b p-4 text-center text-2xl rounded-2xl bg-orange-200/20 dark:bg-orange-900/20 backdrop-blur-md border-white/10 dark:border-gray-700/50 shadow-lg">
        <button v-for="tab in tabs" :key="tab" :class="[
          'px-4 py-2 font-medium transition-all duration-300',
          activeTab === tab
            ? 'border-b-2 border-primary text-primary dark:text-primary'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
        ]" @click="handleTabClick(tab)">
          {{ tab }}
        </button>
      </div>
    </div>

    <div class="mt-4">
      <slot :name="activeTab" />
    </div>
  </div>
</template>
