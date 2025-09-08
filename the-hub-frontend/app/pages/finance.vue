<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Dashboard')
</script>
<template>
  <Tabs :tabs="['Dashboard', 'Budget', 'Category', 'Transaction']" v-model="selectedTab">
    <template #Dashboard>
      <BudgetDashboard></BudgetDashboard>
    </template>
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
