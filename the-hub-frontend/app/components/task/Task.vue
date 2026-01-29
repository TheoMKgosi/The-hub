<script setup lang="ts">
import dayjs from 'dayjs';
import EditIcon from '../ui/svg/EditIcon.vue';
import ThreeDotsIcon from '../ui/svg/ThreeDotsIcon.vue';
import UpArrowIcon from '../ui/svg/UpArrowIcon.vue';
import DownArrowIcon from '../ui/svg/DownArrowIcon.vue';
import DeleteIcon from '../ui/svg/DeleteIcon.vue';

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
  (e: 'deleteTask', id: string): void;
  (e: 'moveTaskUp', id: string): void;
  (e: 'moveTaskDown', id: string): void;
  (e: 'edit', id: string, updates: any): void;
}>()

const draft = reactive({
  title: props.title,
  description: props.description,
  priority: props.priority,
  due_date: props.due_date
})

const isMenuOpen = ref(false)
const isEditing = ref(false)
const menuPosition = ref<'above' | 'below'>('below');
const triggerRef = ref<any>(null); // Ref for the BaseButton

const toggleMenu = () => {
  if (!isMenuOpen.value) {
    const el = triggerRef.value?.$el || triggerRef.value;
    const rect = el.getBoundingClientRect();

    const spaceBelow = window.innerHeight - rect.bottom;

    menuPosition.value = spaceBelow < 200 ? 'above' : 'below';
  }
  isMenuOpen.value = !isMenuOpen.value
}

const startEdit = () => {
  isEditing.value = true
}

const completeBtnClick = () => {
  const newStatus = props.status === 'pending' ? 'complete' : 'pending'
  useTaskStore().editTask({ task_id: props.task_id, status: newStatus })
}

const titleOnly = computed(() => {
  return true
})

const deleteBtnClick = () => {
  emit('deleteTask', props.task_id)
}

const moveUpBtnClick = () => {
  emit('moveTaskUp', props.task_id)
}

const moveDownBtnClick = () => {
  emit('moveTaskDown', props.task_id)
}

const cancelEdit = () => {
  // Reset draft to original prop values
  draft.title = props.title
  draft.description = props.description
  draft.priority = props.priority
  draft.due_date = props.due_date
  isEditing.value = false
}

const saveEdit = () => {
  emit('edit', props.task_id, { ...draft })
  isEditing.value = false
}

const menuItems = [{ text: 'Edit', icon: EditIcon, callBack: startEdit },
{ text: 'Move Up', icon: UpArrowIcon, callBack: moveUpBtnClick },
{ text: 'Move Down', icon: DownArrowIcon, callBack: moveDownBtnClick },
{ text: 'Delete', icon: DeleteIcon, callBack: deleteBtnClick },
]
</script>

<template>
  <div
    class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
    :class="[status === 'complete' ? 'border-success' : 'border-warning',]">
    <div class="flex flex-row justify-between">
      <div v-if="isEditing" class="space-y-3">
        <input v-model="draft.title" class="w-full p-2 border rounded dark:bg-gray-800 dark:text-white"
          placeholder="Task Title" />
        <textarea v-model="draft.description" class="w-full p-2 border rounded dark:bg-gray-800 dark:text-white text-sm"
          placeholder="Description"></textarea>
        <div class="flex gap-2">
          <button @click="saveEdit" class="px-3 py-1 bg-success text-white rounded text-sm font-bold">Save</button>
          <button @click="cancelEdit" class="px-3 py-1 bg-gray-500 text-white rounded text-sm">Cancel</button>
        </div>
      </div>
      <div class="" v-else>
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
            {{ title }}
          </h3>
        </div>
        <span v-if="titleOnly">
          <p class="text-sm text-text-light dark:text-text-dark/80 mb-2">
            {{ description }}
          </p>
          <p class="text-sm text-text-light dark:text-text-dark/60 mb-2">
            {{ due_date ? dayjs(due_date).fromNow() : "" }}
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
        </span>
      </div>

      <div class="flex flex-col justify-between" v-if="titleOnly">
        <!-- Three-dot menu button -->
        <div class="ml-auto">
          <BaseButton ref="triggerRef" @click="toggleMenu" variant="clear" :iconOnly="true" :icon="ThreeDotsIcon" />

          <!-- Dropdown menu -->
          <div v-if="isMenuOpen">
            <Menu :list="menuItems" :position="menuPosition" @close="isMenuOpen = false" />
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
