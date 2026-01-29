<script setup lang="ts">
import type { Task } from '~/types/task'
interface Props {
  taskList: Task[]
}

const props = defineProps<Props>()

const tasks = computed(() => {
  return props.taskList.map((task) => ({
    title: task.title,
    start_time: task.due_date?.getHours() ?? 99,
    end_time: task.due_date?.getHours() ?? 99
  }))

})

const tri_modal = ['Planning', 'Execute', 'Analysis']
</script>
<template>
  <div class="flex flex-col flex-1 min-h-screen">
    <!--Tabs -->
    <div class="layout-tabs flex justify-center py-2">
      <Tabs :tabs="['Tasks', 'Goal']" />
    </div>
    <!-- Control -->
    <div class="layout-controls w-full flex">
      <slot name="control">
        <div>
          <p>Filter</p>
          <div class="flex">
            <BaseButton text="All" variant="primary" class="mr-2" />
            <BaseButton text="Linked" variant="primary" class="mr-2" />
            <BaseButton text="Pending" variant="primary" class="mr-2" />
          </div>
        </div>
        <div>
          <p>Tri-Modal</p>
          <div class="flex">
            <SegmentedControl :texts="tri_modal" />
          </div>
        </div>
      </slot>
    </div>

    <div class="layout-content p-4 flex flex-1 md:flex-row">
      <!-- Tasks -->
      <div class="layout-tasks md:w-64">
        <slot name="tasks">
          <TaskList :tasks="taskList" />
        </slot>
      </div>

      <!-- CalendarSlots -->
      <div class="layout-calendar-slot  flex flex-1 ml-2">
        <slot name="calendar-slot">
          <DateSlots class="grow" label="Today" :tasks="tasks" />
          <DateSlots class="grow" label="Tomorrow" :tasks="tasks" />
        </slot>
      </div>
    </div>

  </div>
</template>
