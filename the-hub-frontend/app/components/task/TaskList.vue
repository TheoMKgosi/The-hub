<script setup lang="ts">
import type { Task } from '~/types/task';
interface Props {
  tasks: Task[]
}

const { tasks = [],} = defineProps<Props>()
</script>

<template>
  <div>
    <div v-if="tasks.length === 0">No Tasks</div>
    <Task v-else v-for="task in tasks" :key="task.task_id" :task_id="task.task_id" :title="task.title"
      :description="task.description" :status="task.status" :due_date="task.due_date" :priority="task.priority"
      :time_estimate_minutes="task.time_estimate_minutes" class="mt-3"/>
  </div>
</template>
