<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Income') 
</script>
<template>
  <Tabs :tabs="['Income', 'Category']" v-model="selectedTab">
    <!--
    <template #Transaction>
      Transaction - coming soon
    </template>
    <template #Budget>
      <FinanceBudget></FinanceBudget>
    </template>
    -->
    <template #Income>
      <FinanceIncome></FinanceIncome>
    </template>
    <template #Category>
      <FinanceCategory></FinanceCategory>
    </template>
  </Tabs>
</template>
