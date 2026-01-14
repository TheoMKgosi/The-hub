<script setup lang="ts">
import type { Task } from '~/types/task'
interface Props {
  taskList: Task[]
}

const { taskList = [] } = defineProps<Props>()
</script>
<template>
  <div class="flex flex-col flex-1 min-h-screen">
    <!--Tabs -->
    <div class="layout-tabs flex justify-center py-2">
        <Tabs :tabs="['Tasks', 'Goal']" />
    </div>
    <!-- Control -->
    <div class="layout-controls w-full">
      <slot name="control">
        <div class="flex">
          <BaseButton text="All" variant="primary" class="mr-2" />
          <BaseButton text="Linked" variant="primary" class="mr-2" />
          <BaseButton text="Pending" variant="primary" class="mr-2" />
        </div>
      </slot>
    </div>

    <div class="layout-content flex flex-1 md:flex-row">
      <!-- Tasks -->
      <div class="layout-tasks md:w-64">
        <slot name="tasks">
          <TaskList :tasks="taskList" />
        </slot>
      </div>

      <!-- CalendarSlots -->
      <div class="layout-calendar-slot flex flex-1 ml-2 shadow-2xs">
        <slot name="calendar-slot">
          <DateSlots class="grow" label="Today"/>
          <DateSlots class="grow" label="Tomorrow"/>
        </slot>
      </div>
    </div>

  </div>
</template>
