<script setup lang="ts">
interface Props {
  list: {
    text: string,
    icon: object,
    callBack: Function
  }[]
  position: 'above' | 'below'
}
const props = defineProps<Props>()

const emit = defineEmits(['close'])

const menuRef = ref<HTMLElement | null>(null)

const handleAction = (item: any) => {
  item.callBack();
  emit('close');
};

const handleClickOutside = (event: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    emit('close');
  }
};

onMounted(() => {
  // Use capture to ensure it runs before other click events if needed
  document.addEventListener('mousedown', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
});
</script>
<template>
  <div class="relative p-3 rounded-2xl">
    <div
      class="absolute right-0 w-48 z-50 shadow-2xl rounded-md overflow-hidden border bg-white dark:bg-gray-900 dark:border-gray-700"
      :class="props.position === 'above' ? 'bottom-full mb-2' : 'top-full mt-2'" 
      ref="menuRef">
      <BaseButton v-for="item in props.list" @click="handleAction(item)" variant="clear" size="full" :text="item.text"
        :icon="item.icon"></BaseButton>
    </div>
  </div>
</template>
