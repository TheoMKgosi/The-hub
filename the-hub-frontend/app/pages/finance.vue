<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Budget')
</script>
<template>
  <Tabs :tabs="['Budget', 'Category', 'Transaction']" v-model="selectedTab">
    <template #Budget>
      <FinanceIncome></FinanceIncome>
    </template>
    <template #Category>
      <FinanceCategory></FinanceCategory>
    </template>
    <template #Transaction>
      <FinanceTransaction></FinanceTransaction>
    </template>
  </Tabs>
</template>
