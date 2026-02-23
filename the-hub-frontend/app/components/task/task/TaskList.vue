<script setup lang="ts">
import type { Task } from '~/types/task';
interface Props {
  tasks: Task[]
}

const { tasks = [], } = defineProps<Props>()

const emit = defineEmits<{
  (e: 'edit', taskId: string): void
  (e: 'moveUp', taskId: string): void
  (e: 'moveDown', taskId: string): void
}>()

const handleEdit = (taskId: string) => {
  emit('edit', taskId)
}

const handleMoveUp = (taskId: string) => {
  emit('moveUp', taskId)
}

const handleMoveDown = (taskId: string) => {
  emit('moveDown', taskId)
}
</script>

<template>
  <div class="p-4 shadow rounded-2xl dark:inset-shadow-sm inset-shadow-gray-500/50">
    <div v-if="tasks.length === 0">No Tasks</div>
    <Task v-else v-for="task in tasks" :key="task.task_id" :task_id="task.task_id" :title="task.title"
      :description="task.description" :status="task.status" :due_date="task.due_date" :priority="task.priority"
      :order="task.order" :time_estimate_minutes="task.time_estimate_minutes" class="mt-3" @edit="handleEdit" @moveUpBtnClick="handleMoveUp" @moveDownBtnClick="handleMoveDown" />
  </div>
</template>
