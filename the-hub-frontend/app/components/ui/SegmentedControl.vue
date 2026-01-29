<script setup lang="ts">
interface Props {
  texts: string[]
}

const props = defineProps<Props>()
</script>
<template>
  <div v-for="text in props.texts">
    <BaseButton :text="text" variant="default" />
  </div>
</template>
