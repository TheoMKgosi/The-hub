<script setup lang="ts">
interface Props {
  time: string
}

const props = defineProps<Props>()
</script>
<template>
  <div class="py-5 border">
    <p>{{ time }}</p>
  </div>
</template>
