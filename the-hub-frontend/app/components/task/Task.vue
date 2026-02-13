<script setup lang="ts">
import EditIcon from '../ui/svg/EditIcon.vue';
import ThreeDotsIcon from '../ui/svg/ThreeDotsIcon.vue';
import UpArrowIcon from '../ui/svg/UpArrowIcon.vue';
import DownArrowIcon from '../ui/svg/DownArrowIcon.vue';
import DeleteIcon from '../ui/svg/DeleteIcon.vue';

const { fromNow } = useDate()

const taskStore = useTaskStore()

interface Props {
  task_id: string,
  status: string,
  title: string,
  description?: string,
  due_date?: Date | null,
  priority?: number,
  time_estimate_minutes?: number
}

const props = withDefaults(defineProps<Props>(), {
  status: 'pending',
  description: '',
  due_date: null,
  priority: 3,
})


const emit = defineEmits<{
  (e: 'completeTask', id: string): void;
  (e: 'moveTaskUp', id: string): void;
  (e: 'moveTaskDown', id: string): void;
  (e: 'edit', id: string): void;
}>()

const isMenuOpen = ref(false)
const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const startEdit = () => {
  isMenuOpen.value = false
  emit('edit', props.task_id)
}

const completeBtnClick = () => {
  const newStatus = props.status === 'pending' ? 'complete' : 'pending'

  taskStore.editTask({ task_id: props.task_id, status: newStatus })
}

const deleteBtnClick = () => {
  taskStore.deleteTask(props.task_id)
}

const moveUpBtnClick = () => {
  emit('moveTaskUp', props.task_id)
}

const moveDownBtnClick = () => {
  emit('moveTaskDown', props.task_id)
}

const handleDoubleClick = () => {
  emit('edit', props.task_id)
}
</script>

<template>
  <div
    class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
    :class="[status === 'complete' ? 'border-success' : 'border-warning',]" @dblclick="handleDoubleClick">
    <div class="flex flex-row justify-between">
      <div class="">
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
            {{ title }}
          </h3>
        </div>
        <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">
          {{ description }}
        </p>
        <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
          {{ due_date ? fromNow(due_date) : "" }}
        </p>
        <div class="flex items-center gap-2 mt-2">
          <input type="checkbox" @click="completeBtnClick" :checked="status === 'complete'"
            class="accent-success w-4 h-4" />
          <span class="text-sm font-medium text-text-light dark:text-text-dark capitalize">{{ status }}</span>
        </div>
        <div>
          <p class="text-sm text-text-light dark:text-text-dark/60 mt-1">
            Priority: {{ priority }}
          </p>
        </div>
      </div>

      <div class="flex flex-col justify-between">
        <!-- Three-dot menu button -->
        <div class="ml-auto">
          <BaseButton @click="toggleMenu" variant="clear" :iconOnly="true" :icon="ThreeDotsIcon"></BaseButton>

          <!-- Dropdown menu -->
          <div v-if="isMenuOpen"
            class="absolute right-4 mt-2 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-2xl border border-surface-light/20 dark:border-surface-dark/20 z-10">
            <div class="py-1">
              <BaseButton @click="startEdit" variant="clear" size="full" text="Edit" :icon="EditIcon"></BaseButton>
              <BaseButton @click="moveUpBtnClick" variant="clear" size="full" text="Move Up" :icon="UpArrowIcon">
              </BaseButton>
              <BaseButton @click="moveDownBtnClick" variant="clear" size="full" text="Move Down" :icon="DownArrowIcon">
              </BaseButton>
              <BaseButton @click="deleteBtnClick" variant="clear" size="full" text="Delete" :icon="DeleteIcon">
              </BaseButton>
            </div>
          </div>
        </div>
        <div v-if="time_estimate_minutes" class="flex items-center gap-1 mt-1">
          <span class="hidden sm:inline">⏱️</span>
          <span class="text-sm text-text-light dark:text-text-dark/60">
            Est: {{ Math.floor(time_estimate_minutes / 60) }}h {{ time_estimate_minutes % 60 }}m
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
