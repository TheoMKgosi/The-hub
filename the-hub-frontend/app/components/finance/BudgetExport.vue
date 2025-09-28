<script setup>
const budgetStore = useBudgetStore()
const transactionStore = useTransactionStore()
const { addToast } = useToast()

// Export settings
const exportSettings = reactive({
  format: 'csv',
  includeAnalytics: true,
  includeTransactions: true,
  dateRange: 'current',
  categories: []
})

// Available formats
const formats = [
  { value: 'csv', label: 'CSV' },
  { value: 'json', label: 'JSON' },
  { value: 'pdf', label: 'PDF Report' }
]

// Date ranges
const dateRanges = [
  { value: 'current', label: 'Current Period' },
  { value: 'last_month', label: 'Last Month' },
  { value: 'last_3_months', label: 'Last 3 Months' },
  { value: 'last_6_months', label: 'Last 6 Months' },
  { value: 'year_to_date', label: 'Year to Date' },
  { value: 'all', label: 'All Time' }
]

// Export budget data
const exportBudgets = async () => {
  try {
    const data = {
      budgets: budgetStore.budgets,
      analytics: exportSettings.includeAnalytics ? budgetStore.analytics : [],
      transactions: exportSettings.includeTransactions ? transactionStore.transactions : [],
      exportDate: new Date().toISOString(),
      settings: exportSettings
    }

    let content = ''
    let filename = ''
    let mimeType = ''

    switch (exportSettings.format) {
      case 'csv':
        content = convertToCSV(data)
        filename = `budget-export-${new Date().toISOString().split('T')[0]}.csv`
        mimeType = 'text/csv'
        break
      case 'json':
        content = JSON.stringify(data, null, 2)
        filename = `budget-export-${new Date().toISOString().split('T')[0]}.json`
        mimeType = 'application/json'
        break
      case 'pdf':
        await generatePDF(data)
        return
    }

    // Create and download file
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    addToast('Export completed successfully', 'success')
  } catch (error) {
    addToast('Export failed', 'error')
  }
}

// Convert data to CSV format
const convertToCSV = (data) => {
  const rows = []

  // Budgets header
  rows.push(['Budgets'])
  rows.push(['ID', 'Category', 'Amount', 'Start Date', 'End Date'])
  data.budgets.forEach(budget => {
    rows.push([
      budget.budget_id,
      budget.Category.name,
      budget.amount,
      budget.start_date,
      budget.end_date
    ])
  })

  // Analytics header
  if (data.analytics.length > 0) {
    rows.push([])
    rows.push(['Budget Analytics'])
    rows.push(['Category', 'Budget Amount', 'Spent Amount', 'Remaining', 'Utilization %', 'Daily Rate', 'Status'])
    data.analytics.forEach(analytic => {
      rows.push([
        analytic.category_name,
        analytic.budget_amount,
        analytic.spent_amount,
        analytic.remaining_amount,
        analytic.utilization_rate,
        analytic.daily_spending_rate,
        analytic.status
      ])
    })
  }

  // Transactions header
  if (data.transactions.length > 0) {
    rows.push([])
    rows.push(['Transactions'])
    rows.push(['Description', 'Amount', 'Type', 'Date', 'Category'])
    data.transactions.forEach(transaction => {
      rows.push([
        transaction.description,
        transaction.amount,
        transaction.type,
        transaction.date,
        transaction.category?.name || ''
      ])
    })
  }

  return rows.map(row => row.map(cell => `"${cell}"`).join(',')).join('\n')
}

// Generate PDF report
const generatePDF = async (data) => {
  // In a real implementation, you would use a PDF library like jsPDF
  addToast('PDF export feature coming soon', 'info')
}

// Tax preparation export
const exportForTaxes = async () => {
  const taxData = {
    year: new Date().getFullYear(),
    income: transactionStore.transactions
      .filter(t => t.type === 'income')
      .reduce((sum, t) => sum + t.amount, 0),
    expenses: transactionStore.transactions
      .filter(t => t.type === 'expense')
      .reduce((sum, t) => sum + t.amount, 0),
    categorizedExpenses: transactionStore.transactions
      .filter(t => t.type === 'expense')
      .reduce((acc, t) => {
        const category = t.category?.name || 'Uncategorized'
        acc[category] = (acc[category] || 0) + t.amount
        return acc
      }, {}),
    exportDate: new Date().toISOString()
  }

  const content = JSON.stringify(taxData, null, 2)
  const blob = new Blob([content], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `tax-prep-${new Date().getFullYear()}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)

  addToast('Tax preparation export completed', 'success')
}

onMounted(async () => {
  if (budgetStore.budgets.length === 0) {
    await budgetStore.fetchBudget()
  }
  if (transactionStore.transactions.length === 0) {
    await transactionStore.fetchTransactions()
  }
})
</script>

<template>
  <div class="p-4 max-w-4xl mx-auto">
    <h2 class="text-2xl font-bold mb-6 text-text-light dark:text-text-dark">Export Budget Data</h2>

    <!-- Export Settings -->
    <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark mb-6">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Export Settings</h3>
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Format</label>
            <select v-model="exportSettings.format"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
              <option v-for="format in formats" :key="format.value" :value="format.value">
                {{ format.label }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-text-light dark:text-text-dark mb-1">Date Range</label>
            <select v-model="exportSettings.dateRange"
              class="w-full px-3 py-2 border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary">
              <option v-for="range in dateRanges" :key="range.value" :value="range.value">
                {{ range.label }}
              </option>
            </select>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <label class="flex items-center">
            <input v-model="exportSettings.includeAnalytics" type="checkbox"
              class="rounded border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-primary focus:ring-primary" />
            <span class="ml-2 text-sm text-text-light dark:text-text-dark">Include budget analytics</span>
          </label>

          <label class="flex items-center">
            <input v-model="exportSettings.includeTransactions" type="checkbox"
              class="rounded border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-primary focus:ring-primary" />
            <span class="ml-2 text-sm text-text-light dark:text-text-dark">Include transactions</span>
          </label>
        </div>

        <div class="flex gap-4">
          <UiButton @click="exportBudgets" variant="primary" size="md">
            Export Budget Data
          </UiButton>
          <UiButton @click="exportForTaxes" variant="secondary" size="md">
            Export for Tax Preparation
          </UiButton>
        </div>
      </div>
    </div>

    <!-- Export History -->
    <div class="bg-surface-light dark:bg-surface-dark p-6 rounded-lg shadow-md border border-surface-light dark:border-surface-dark">
      <h3 class="text-lg font-semibold mb-4 text-text-light dark:text-text-dark">Recent Exports</h3>
      <div class="text-center py-8 text-text-light dark:text-text-dark/60">
        <p>No recent exports. Your exported files will appear here.</p>
      </div>
    </div>
  </div>
</template>
