<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Plan')
</script>
<template>
   <Tabs :tabs="['Plan', 'Income', 'Transactions', 'Receipts', 'Goals']" v-model="selectedTab">
     <template #Plan>
       <MonthPlan></MonthPlan>
     </template>
     <template #Income>
       <Income></Income>
     </template>
     <template #Transactions>
       <Transaction></Transaction>
     </template>
     <template #Receipts>
       <Receipt></Receipt>
     </template>
     <template #Goals>
       <FinancialGoal></FinancialGoal>
     </template>
   </Tabs>
</template>