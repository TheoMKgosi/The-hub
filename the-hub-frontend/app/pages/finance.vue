<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Budget')
</script>
<template>
   <UiTabs :tabs="['Budget', 'Transaction', 'Receipt']" v-model="selectedTab">
     <template #Budget>
       <FinanceIncome></FinanceIncome>
     </template>
     <template #Transaction>
       <FinanceTransaction></FinanceTransaction>
     </template>
     <template #Receipt>
       <FinanceReceipt></FinanceReceipt>
     </template>
   </UiTabs>
</template>
