<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Income') 
</script>
<template>
  <Tabs :tabs="['Income', 'Category', 'Transaction']" v-model="selectedTab">
    <template #Income>
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
