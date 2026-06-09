<script setup lang="ts">
import { ref, computed } from 'vue'

interface GoalTaskRecommendation {
  title: string
  description: string
  priority: number
  estimated_hours: number
  reasoning: string
}

const props = defineProps<{
  isOpen: boolean
  recommendations: GoalTaskRecommendation[]
  loading?: boolean
}>()

const emit = defineEmits<{
  close: []
  addTask: [recommendation: GoalTaskRecommendation]
}>()

const getPriorityColor = (priority: number) => {
  switch (priority) {
    case 1:
    case 2:
      return 'text-green-600 dark:text-green-400'
    case 3:
      return 'text-yellow-600 dark:text-yellow-400'
    case 4:
    case 5:
      return 'text-red-600 dark:text-red-400'
    default:
      return 'text-text-light/60 dark:text-text-dark/60'
  }
}

const handleAddTask = (recommendation: GoalTaskRecommendation) => {
  emit('addTask', recommendation)
}
</script>

<template>
  <Teleport to="#plan">
    <div v-if="isOpen" class="fixed inset-0 bg-black/50 dark:bg-black/70 flex items-center justify-center p-4 z-50"
      @click="emit('close')">
      <div
        class="bg-surface-light dark:bg-surface-dark rounded-lg w-full max-w-2xl max-h-[85vh] overflow-hidden shadow-xl border border-surface-light/20 dark:border-surface-dark/20"
        @click.stop>

        <!-- Modal Header -->
        <div class="flex items-center justify-between p-4 border-b border-surface-light/20 dark:border-surface-dark/20">
          <h2 class="text-lg font-semibold text-text-light dark:text-text-dark flex items-center gap-2">
            <svg class="w-5 h-5 text-primary" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd"
                d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                clip-rule="evenodd" />
            </svg>
            AI Task Recommendations
            <span v-if="recommendations.length" class="text-sm font-normal text-text-light/60 dark:text-text-dark/60">
              ({{ recommendations.length }})
            </span>
          </h2>
          <BaseButton @click="emit('close')" variant="default" size="sm" text="Close" />
        </div>

        <!-- Modal Body -->
        <div class="p-4 overflow-y-auto max-h-[60vh]">
          <!-- Loading State -->
          <div v-if="loading" class="flex items-center justify-center py-12">
            <div class="flex flex-col items-center gap-3">
              <svg class="w-8 h-8 animate-spin text-primary" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <p class="text-sm text-text-light/60 dark:text-text-dark/60">Analyzing goal and generating recommendations...</p>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else-if="recommendations.length === 0" class="text-center py-12">
            <p class="text-text-light/60 dark:text-text-dark/60">No recommendations available</p>
            <p class="text-sm text-text-light/40 dark:text-text-dark/40 mt-2">Try adding more details to your goal</p>
          </div>

          <!-- Recommendations List -->
          <div v-else class="space-y-3">
            <div v-for="(recommendation, index) in recommendations" :key="index"
              class="bg-primary/5 dark:bg-primary/10 rounded-lg p-4 border border-primary/20">
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="text-xs text-text-light/40 dark:text-text-dark/40">#{{ index + 1 }}</span>
                    <h5 class="text-sm font-medium text-text-light dark:text-text-dark">
                      {{ recommendation.title }}
                    </h5>
                  </div>
                  <p class="text-xs text-text-light/80 dark:text-text-dark/80 mb-3">
                    {{ recommendation.description }}
                  </p>
                  <div class="flex items-center gap-4 text-xs text-text-light/60 dark:text-text-dark/60">
                    <span :class="getPriorityColor(recommendation.priority)" class="flex items-center gap-1">
                      <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd"
                          d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
                          clip-rule="evenodd" />
                      </svg>
                      Priority {{ recommendation.priority }}
                    </span>
                    <span class="flex items-center gap-1">
                      <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                          clip-rule="evenodd" />
                      </svg>
                      ~{{ recommendation.estimated_hours }}h
                    </span>
                  </div>
                  <p class="text-xs text-text-light/70 dark:text-text-dark/70 mt-2 italic">
                    {{ recommendation.reasoning }}
                  </p>
                </div>
                <BaseButton @click="handleAddTask(recommendation)" variant="primary" size="sm" text="Add" class="shrink-0" />
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-4 border-t border-surface-light/20 dark:border-surface-dark/20 flex justify-end">
          <BaseButton @click="emit('close')" variant="default" size="md" text="Done" />
        </div>
      </div>
    </div>
  </Teleport>
</template>