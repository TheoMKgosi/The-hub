<script setup lang="ts">
import type { Task } from '~/types/task'
interface Props {
  taskList: Task[]
}

const props = defineProps<Props>()

const tasks = computed(() => {
  return props.taskList.map((task: Task) => {
    const dateValue = task.due_date ? new Date(task.due_date) : null
    const isValid = dateValue && !isNaN(dateValue.getTime());
    const hour = isValid ? dateValue.getHours() : 99;
    return {
      title: task.title,
      start_time: hour,
      end_time: hour
    }
  })
})

// Add user settings integration
const userSettings = ref<any>({})

// Load user settings on mount
onMounted(async () => {
  try {
    const auth = useAuthStore()
    if (!auth.user?.user_id) return

    const { $api } = useNuxtApp()
    const response = await $api(`/users/${auth.user.user_id}/settings`)
    userSettings.value = response.settings || {}
  } catch (error) {
    console.warn('Failed to load user settings for tri-modal:', error)
  }
})

const tri_interface = computed(() => {
  return userSettings.value.task?.['tri_modal'] === true
})

const tri_modal = ['Planning', 'Execute', 'Analysis']
</script>
<template>
  <div id="plan" class="flex flex-col flex-1 min-h-screen">
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
        <div v-if="tri_interface">
          <p>Tri-Modal</p>
          <div class="flex">
            <SegmentedControl :texts="tri_modal" />
          </div>
        </div>
      </slot>
    </div>

    <div class="layout-content p-4 flex flex-1 flex-col md:flex-row">
      <div class="layout-tasks basis-1/3 grow">
        <slot name="tasks">
          <TaskList :tasks="taskList" />
        </slot>
      </div>

      <div v-if="tri_interface" class="layout-calendar-slot flex basis-2/3 ml-2 grow">
        <slot name="calendar-slot" class="flex w-full">
          <DateSlots class="grow basis-1/2" label="Today" :tasks="tasks" />
          <DateSlots class="grow basis-1/2" label="Tomorrow" :tasks="tasks" />
        </slot>
      </div>
    </div>
    <FormTask />
  </div>
</template>
