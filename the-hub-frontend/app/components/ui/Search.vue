<script setup lang="ts">
interface Props {
  placeholder: string 
}
</script>
<template>
  <div class="flex flex-1">
    <BaseInput placeholder="Search " />
    <BaseButton type="button" text="Search" />
  </div>
</template>
