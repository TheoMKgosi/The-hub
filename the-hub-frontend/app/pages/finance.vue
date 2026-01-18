<script setup lang="ts">
const selectedTab = useState('financeTab', () => 'Budget')
</script>
<template>
   <Tabs :tabs="['Budget', 'Transaction', 'Receipt']" v-model="selectedTab">
     <template #Budget>
       <Income></Income>
     </template>
     <template #Transaction>
       <Transaction></Transaction>
     </template>
     <template #Receipt>
       <Receipt></Receipt>
     </template>
   </Tabs>
</template>
