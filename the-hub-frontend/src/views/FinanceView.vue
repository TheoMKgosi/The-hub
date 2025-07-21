<script setup lang="ts">
import Tabs from '../components/Tabs.vue'
import Budget from '../components/finance/Budget.vue'
import Category from '../components/finance/Category.vue'
</script>

<template>
  <Tabs :tabs="['Overview', 'Budget', 'Category', 'Transaction']" :defaultTab="Overview">
    <template #Overview>
      Overview - coming soon
    </template>
    <template #Transaction>
      Transaction - coming soon
    </template>
    <template #Budget>
      <Budget></Budget>
    </template>
    <template #Category>
      <Category></Category>
    </template>
  </Tabs>
</template>
