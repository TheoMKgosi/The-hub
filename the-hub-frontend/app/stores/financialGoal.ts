interface FinancialGoal {
  goal_id: string
  title: string
  description: string
  target_amount: number | null
  current_amount: number
  type: 'savings' | 'checklist'
  status: 'active' | 'completed' | 'cancelled'
  priority: number | null
  target_date: string | null
  category_id: string | null
  color: string
  created_at: string
}

export interface FinancialGoalResponse {
  goals: FinancialGoal[]
}

export const useFinancialGoalStore = defineStore('financialGoal', () => {
  const goals = ref<FinancialGoal[]>([])
  const loading = ref(false)
  const creating = ref(false)
  const updating = ref(false)
  const deleting = ref(false)
  const toast  = useToast()

  async function fetchGoals() {
    const { $api } = useNuxtApp()
    loading.value = true
    try {
      const response = await $api<FinancialGoalResponse>('financial-goals')
      if (response) goals.value = response.goals
    } catch (err) {
      toast.add({title: "Error", description: "Failed to load financial goals", color: "error"})
      console.log(err)
    } finally {
      loading.value = false
    }
  }

  async function submitForm(payload: Partial<FinancialGoal>) {
    creating.value = true

    const optimisticGoal: FinancialGoal = {
      goal_id: `temp-${Date.now()}`,
      title: payload.title || '',
      description: payload.description || '',
      target_amount: payload.target_amount ?? null,
      current_amount: payload.current_amount || 0,
      type: payload.type || 'savings',
      status: 'active',
      priority: payload.priority ?? null,
      target_date: payload.target_date ?? null,
      category_id: payload.category_id ?? null,
      color: payload.color || '#3B82F6',
      created_at: new Date().toISOString(),
    }

    goals.value.unshift(optimisticGoal)

    try {
      const { $api } = useNuxtApp()
      const data = await $api<FinancialGoal>('financial-goals', {
        method: 'POST',
        body: JSON.stringify(payload),
      })

      const optimisticIndex = goals.value.findIndex(g => g.goal_id === optimisticGoal.goal_id)
      if (optimisticIndex !== -1) {
        goals.value[optimisticIndex] = data
      }

      toast.add({title: "Finance", description: "Goal created successfully", color: "success"})
    } catch (err) {
      goals.value = goals.value.filter(g => g.goal_id !== optimisticGoal.goal_id)
      toast.add({title: "Error", description: "Failed to create goal", color: "error"})
    } finally {
      creating.value = false
    }
  }

  async function editGoal(payload: Partial<FinancialGoal> & { goal_id: string }) {
    updating.value = true

    const originalIndex = goals.value.findIndex(g => g.goal_id === payload.goal_id)
    const original = originalIndex !== -1 ? { ...goals.value[originalIndex] } : null

    if (originalIndex !== -1) {
      goals.value[originalIndex] = { ...goals.value[originalIndex], ...payload }
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<FinancialGoal>(`financial-goals/${payload.goal_id}`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
      })

      if (originalIndex !== -1 && data) {
        goals.value[originalIndex] = data
      }

      toast.add({title: "Finance", description: "Goal updated successfully", color: "success"})
    } catch (err) {
      if (original && originalIndex !== -1) {
        goals.value[originalIndex] = original
      }
      toast.add({title: "Error", description: "Failed to update goal", color: "error"})
    } finally {
      updating.value = false
    }
  }

  async function deleteGoal(goalId: string) {
    deleting.value = true

    const goalToDelete = goals.value.find(g => g.goal_id === goalId)
    if (!goalToDelete) {
      toast.add({title: "Error", description: "Goal not found", color: "error"})
      deleting.value = false
      return
    }

    goals.value = goals.value.filter(g => g.goal_id !== goalId)

    try {
      const { $api } = useNuxtApp()
      await $api(`financial-goals/${goalId}`, { method: 'DELETE' })
      toast.add({title: "Finance", description: "Goal deleted successfully", color: "success"})
    } catch (err) {
      goals.value.push(goalToDelete)
      toast.add({title: "Error", description: "Failed to delete goal", color: "error"})
    } finally {
      deleting.value = false
    }
  }

  async function addProgress(goalId: string, amount: number) {
    const goal = goals.value.find(g => g.goal_id === goalId)
    if (!goal) return

    const previousAmount = goal.current_amount
    goal.current_amount += amount

    if (goal.target_amount && goal.current_amount >= goal.target_amount) {
      goal.status = 'completed'
    }

    try {
      const { $api } = useNuxtApp()
      const data = await $api<FinancialGoal>(`financial-goals/${goalId}/progress`, {
        method: 'PATCH',
        body: JSON.stringify({ amount }),
      })

      const index = goals.value.findIndex(g => g.goal_id === goalId)
      if (index !== -1) {
        goals.value[index] = data
      }

      toast.add({title: "Finance", description: `Added $${amount.toFixed(2)} to goal`, color: "success"})
    } catch (err) {
      goal.current_amount = previousAmount
      goal.status = previousAmount >= (goal.target_amount || Infinity) ? 'completed' : 'active'
      toast.add({title: "Error", description: "Failed to update progress", color: "error"})
    }
  }

  async function toggleChecklist(goalId: string) {
    const goal = goals.value.find(g => g.goal_id === goalId)
    if (!goal || goal.type !== 'checklist') return

    const newStatus = goal.status === 'completed' ? 'active' : 'completed'
    const previousStatus = goal.status
    goal.status = newStatus

    try {
      const { $api } = useNuxtApp()
      await $api(`financial-goals/${goalId}`, {
        method: 'PATCH',
        body: JSON.stringify({ status: newStatus }),
      })
    } catch (err) {
      goal.status = previousStatus
      toast.add({title: "Error", description: "Failed to update goal", color: "error"})
    }
  }

  function reset() {
    goals.value = []
  }

  return {
    goals,
    loading,
    creating,
    updating,
    deleting,
    fetchGoals,
    submitForm,
    editGoal,
    deleteGoal,
    addProgress,
    toggleChecklist,
    reset,
  }
})
