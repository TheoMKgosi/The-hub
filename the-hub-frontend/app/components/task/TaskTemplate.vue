<script setup lang="ts">
import type { Task } from '~/types/task'
interface Props {
  taskList: Task[]
}

const props = defineProps<Props>()

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

// Computed property for tri-modal visibility
const showTriModal = computed(() => {
  return userSettings.value.task?.['tri-modal'] === true
})

const tasks = computed(() => {
  return props.taskList.map((task) => ({
    title: task.title,
    start_time: task.due_date?.getHours() ?? 99,
    end_time: task.due_date?.getHours() ?? 99
  }))

})

const showAddModal = ref(false)

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
        <div v-if="showTriModal">
          <p>Tri-Modal</p>
          <div class="flex">
            <SegmentedControl :texts="tri_modal" />
          </div>
        </div>
      </slot>
    </div>

    <div class="layout-content p-4 flex flex-1 md:flex-row">
      <!-- Tasks -->
      <div class="layout-tasks grow">
        <TaskList :tasks="taskList" />
      </div>

      <!-- CalendarSlots -->
      <div v-show="showTriModal" class="layout-calendar-slot  flex flex-1 ml-2 grow">
        <DateSlots class="grow-2" label="Today" :tasks="tasks" />
        <DateSlots class="grow-2" label="Tomorrow" :tasks="tasks" />
      </div>
    </div>

  </div>
  <FormTask />
</template>
