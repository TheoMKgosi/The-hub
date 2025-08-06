<script setup lang="ts">
import Tabs from '../components/Tabs.vue'
import Budget from '../components/finance/Budget.vue'
import Category from '../components/finance/Category.vue'
import Income from '../components/finance/Income.vue'

</script>

<template>
  <Tabs :tabs="['Transaction','Income', 'Category']">
    <template #Transaction>
      Transaction - coming soon
    </template>
    <template #Budget>
      <Budget></Budget>
    </template>
    <template #Income>
      <Income></Income>
    </template>
    <template #Category>
      <Category></Category>
    </template>
  </Tabs>
</template>
