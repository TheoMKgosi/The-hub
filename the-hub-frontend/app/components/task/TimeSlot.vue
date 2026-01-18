<script setup lang="ts">
interface Props {
  time: string
}

const props = defineProps<Props>()
</script>
<template>
  <div class="p-5 border dark:text-white">
    <p class="font-bold">{{ time }}</p>
  </div>
</template>
