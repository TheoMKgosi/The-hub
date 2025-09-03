<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useTaskStore } from '@/stores/tasks'

const taskStore = useTaskStore()
const showManager = ref(false)

const recurrenceRules = ref([])
const newRule = reactive({
  name: '',
  description: '',
  frequency: 'weekly',
  interval: 1,
  by_day: '',
  title_template: '',
  description_template: '',
  priority: null as number | null,
  time_estimate_minutes: null as number | null,
  due_date_offset_days: null as number | null,
})

const frequencyOptions = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'yearly', label: 'Yearly' },
]

const dayOptions = [
  { value: '0', label: 'Sunday' },
  { value: '1', label: 'Monday' },
  { value: '2', label: 'Tuesday' },
  { value: '3', label: 'Wednesday' },
  { value: '4', label: 'Thursday' },
  { value: '5', label: 'Friday' },
  { value: '6', label: 'Saturday' },
]

const loadRecurrenceRules = async () => {
  recurrenceRules.value = await taskStore.getRecurrenceRules()
}

const createRule = async () => {
  if (!newRule.name.trim() || !newRule.title_template.trim()) return

  try {
    await taskStore.createRecurrenceRule({
      ...newRule,
      interval: newRule.interval || 1,
      priority: newRule.priority || undefined,
      time_estimate_minutes: newRule.time_estimate_minutes || undefined,
      due_date_offset_days: newRule.due_date_offset_days || undefined,
    })

    // Reset form
    Object.assign(newRule, {
      name: '',
      description: '',
      frequency: 'weekly',
      interval: 1,
      by_day: '',
      title_template: '',
      description_template: '',
      priority: null,
      time_estimate_minutes: null,
      due_date_offset_days: null,
    })

    await loadRecurrenceRules()
  } catch (error) {
    console.error('Failed to create recurrence rule:', error)
  }
}

const generateTasks = async (ruleId: string, count: number = 1) => {
  try {
    await taskStore.generateRecurringTasks(ruleId, count)
  } catch (error) {
    console.error('Failed to generate tasks:', error)
  }
}

const openManager = async () => {
  showManager.value = true
  await loadRecurrenceRules()
}

const closeManager = () => {
  showManager.value = false
}
</script>

<template>
  <div>
    <!-- Trigger Button -->
    <UiButton @click="openManager" variant="default" size="sm">
      <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
      </svg>
      Recurring Tasks
    </UiButton>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showManager" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50" @click="closeManager">
        <div class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-4xl max-h-[90vh] overflow-y-auto shadow-xl border border-surface-light/30 dark:border-surface-dark/30" @click.stop>

          <!-- Header -->
          <div class="flex items-center justify-between p-6 border-b border-surface-light/20 dark:border-surface-dark/20">
            <h2 class="text-xl font-semibold text-text-light dark:text-text-dark">Recurring Tasks Manager</h2>
            <UiButton @click="closeManager" variant="default" size="sm" class="p-2">
              ×
            </UiButton>
          </div>

          <div class="p-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

              <!-- Create New Rule -->
              <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4">
                <h3 class="text-lg font-medium mb-4 text-text-light dark:text-text-dark">Create Recurrence Rule</h3>

                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Rule Name</label>
                    <input
                      v-model="newRule.name"
                      type="text"
                      placeholder="e.g., Weekly Standup"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Frequency</label>
                    <select
                      v-model="newRule.frequency"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                    >
                      <option v-for="option in frequencyOptions" :key="option.value" :value="option.value">
                        {{ option.label }}
                      </option>
                    </select>
                  </div>

                  <div v-if="newRule.frequency === 'weekly'">
                    <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Days of Week</label>
                    <div class="flex flex-wrap gap-2">
                      <label v-for="day in dayOptions" :key="day.value" class="flex items-center">
                        <input
                          v-model="newRule.by_day"
                          type="checkbox"
                          :value="day.value"
                          class="mr-2"
                        />
                        <span class="text-sm text-text-light dark:text-text-dark">{{ day.label }}</span>
                      </label>
                    </div>
                  </div>

                  <div>
                    <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Task Title Template</label>
                    <input
                      v-model="newRule.title_template"
                      type="text"
                      placeholder="e.g., {date} Standup Meeting"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium mb-1 text-text-light dark:text-text-dark">Priority</label>
                    <select
                      v-model="newRule.priority"
                      class="w-full px-3 py-2 border border-surface-light/30 dark:border-surface-dark/30 bg-surface-light/20 dark:bg-surface-dark/20 text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                    >
                      <option :value="null">No priority</option>
                      <option :value="1">1 - Low</option>
                      <option :value="2">2</option>
                      <option :value="3">3 - Medium</option>
                      <option :value="4">4</option>
                      <option :value="5">5 - High</option>
                    </select>
                  </div>

                  <UiButton @click="createRule" variant="primary" size="sm" class="w-full" :disabled="!newRule.name.trim() || !newRule.title_template.trim()">
                    Create Rule
                  </UiButton>
                </div>
              </div>

              <!-- Existing Rules -->
              <div class="bg-surface-light/10 dark:bg-surface-dark/10 rounded-lg p-4">
                <h3 class="text-lg font-medium mb-4 text-text-light dark:text-text-dark">Existing Rules</h3>

                <div v-if="recurrenceRules.length === 0" class="text-center py-8 text-text-light/60 dark:text-text-dark/60">
                  No recurrence rules yet
                </div>

                <div v-else class="space-y-3">
                  <div v-for="rule in recurrenceRules" :key="rule.recurrence_rule_id"
                    class="bg-surface-light/20 dark:bg-surface-dark/20 rounded-lg p-3 border border-surface-light/30 dark:border-surface-dark/30">

                    <div class="flex items-center justify-between mb-2">
                      <h4 class="font-medium text-text-light dark:text-text-dark">{{ rule.name }}</h4>
                      <span class="text-xs px-2 py-1 rounded-full bg-primary/10 dark:bg-primary/20 text-primary">
                        {{ rule.frequency }}
                      </span>
                    </div>

                    <p v-if="rule.description" class="text-sm text-text-light/80 dark:text-text-dark/80 mb-2">
                      {{ rule.description }}
                    </p>

                    <div class="flex items-center gap-2">
                      <UiButton @click="generateTasks(rule.recurrence_rule_id, 1)" variant="default" size="xs">
                        Generate 1 Task
                      </UiButton>
                      <UiButton @click="generateTasks(rule.recurrence_rule_id, 5)" variant="default" size="xs">
                        Generate 5 Tasks
                      </UiButton>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>