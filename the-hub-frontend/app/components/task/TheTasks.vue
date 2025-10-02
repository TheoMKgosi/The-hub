<script setup lang="ts">
import draggable from "vuedraggable";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime.js";

const taskStore = useTaskStore();

callOnce(async () => {
  if (taskStore.tasks.length === 0) await taskStore.fetchTasks();
});

const filter = ref<
  | "all"
  | "complete"
  | "pending"
  | "linked"
  | "priority-low"
  | "priority-medium"
  | "priority-high"
>("all");
const priorityFilter = ref<"all" | "1" | "2" | "3" | "4" | "5">("all");

dayjs.extend(relativeTime);

// Mobile reordering state
const isReorderMode = ref(false);
const reorderTaskId = ref<string | null>(null);
const longPressTimer = ref<NodeJS.Timeout | null>(null);
const longPressDelay = 500; // ms

const setFilter = (filterOption: string) => {
  filter.value = filterOption as any;
};
const searchQuery = ref("");

const isFiltering = computed(
  () =>
    (filter.value !== "all" &&
      filter.value !== "linked" &&
      !filter.value.startsWith("priority")) ||
    priorityFilter.value !== "all" ||
    searchQuery.value !== "",
);

// Detect mobile device
const isMobile = computed(() => {
  if (process.client) {
    return window.innerWidth < 768 || "ontouchstart" in window;
  }
  return false;
});
const matchFilter = (task) => {
  // Handle status filter
  if (filter.value === "complete" || filter.value === "pending") {
    if (task.status !== filter.value) return false;
  }

  // Handle all filter - show only tasks not linked to goals
  if (filter.value === "all") {
    if (task.goal_id) return false;
  }

  // Handle linked filter - show only tasks linked to goals
  if (filter.value === "linked") {
    if (!task.goal_id) return false;
  }

  // Handle priority group filters
  if (filter.value === "priority-low" && task.priority !== 1) return false;
  if (
    filter.value === "priority-medium" &&
    task.priority !== 2 &&
    task.priority !== 3
  )
    return false;
  if (
    filter.value === "priority-high" &&
    task.priority !== 4 &&
    task.priority !== 5
  )
    return false;

  // Handle specific priority filter (from dropdown)
  if (
    priorityFilter.value !== "all" &&
    task.priority !== parseInt(priorityFilter.value)
  )
    return false;

  // Handle search query
  if (
    searchQuery.value &&
    !task.title.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
    !task.description.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
    return false;

  return true;
};

const filteredTasks = computed(() => {
  return taskStore.tasks.filter(matchFilter);
});

// Editing a task
const editingTaskId = ref<number | null>(null);
const editFormData = reactive({
  title: "",
  description: "",
  due_date: null,
  priority: 3,
  status: "pending",
});

// Task menu state
const openMenuTaskId = ref<number | null>(null);

const startEdit = (task) => {
  editingTaskId.value = task.task_id;
  Object.assign(editFormData, { ...task });
  openMenuTaskId.value = null; // Close menu when editing
};

const cancelEdit = () => {
  editingTaskId.value = null;
};

const saveEdit = async (id: number) => {
  const updated = { ...editFormData, task_id: id };
  if (updated.due_date) {
    const date = new Date(updated.due_date);
    updated.due_date = date.toISOString();
  }

  await taskStore.editTask(updated);
  editingTaskId.value = null;
};

const completeTask = async (task) => {
  if (task.status == "pending") {
    task.status = "complete";
    taskStore.completeTask(task);
  } else {
    task.status = "pending";
    taskStore.completeTask(task);
  }
};

const deleteTask = async (id: number) => {
  taskStore.deleteTask(id);
  openMenuTaskId.value = null; // Close menu after delete
};

// Menu functions
const toggleMenu = (taskId: number) => {
  openMenuTaskId.value = openMenuTaskId.value === taskId ? null : taskId;
};

const closeMenu = () => {
  openMenuTaskId.value = null;
};

// Click outside handler
const handleClickOutside = (event: Event) => {
  const target = event.target as HTMLElement;
  if (!target.closest(".task-menu-container")) {
    closeMenu();
  }
};

// Lifecycle hooks for click outside handling
onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

// Long press reordering (both mobile and desktop)
const startLongPress = (taskId: string) => {
  longPressTimer.value = setTimeout(() => {
    enterReorderMode(taskId);
  }, longPressDelay);
};

const cancelLongPress = () => {
  if (longPressTimer.value) {
    clearTimeout(longPressTimer.value);
    longPressTimer.value = null;
  }
};

const enterReorderMode = (taskId: string) => {
  isReorderMode.value = true;
  reorderTaskId.value = taskId;

  // Add haptic feedback if available
  if (navigator.vibrate) {
    navigator.vibrate(50);
  }
};

const exitReorderMode = () => {
  isReorderMode.value = false;
  reorderTaskId.value = null;
  cancelLongPress();
};

const moveTaskUp = async (taskId: string) => {
  const currentIndex = taskStore.tasks.findIndex(
    (task) => task.task_id === taskId,
  );
  if (currentIndex > 0) {
    const newTasks = [...taskStore.tasks];
    const temp = newTasks[currentIndex];
    newTasks[currentIndex] = newTasks[currentIndex - 1];
    newTasks[currentIndex - 1] = temp;
    taskStore.tasks = newTasks;
    await reorderTasks();

    // Haptic feedback
    if (navigator.vibrate) {
      navigator.vibrate(30);
    }
  }
};

const moveTaskDown = async (taskId: string) => {
  const currentIndex = taskStore.tasks.findIndex(
    (task) => task.task_id === taskId,
  );
  if (currentIndex < taskStore.tasks.length - 1) {
    const newTasks = [...taskStore.tasks];
    const temp = newTasks[currentIndex];
    newTasks[currentIndex] = newTasks[currentIndex + 1];
    newTasks[currentIndex + 1] = temp;
    taskStore.tasks = newTasks;
    await reorderTasks();

    // Haptic feedback
    if (navigator.vibrate) {
      navigator.vibrate(30);
    }
  }
};

// Touch gesture handling for swipe
const touchStartY = ref(0);
const touchEndY = ref(0);
const swipeThreshold = 50;

const handleTouchStart = (event: TouchEvent, taskId: string) => {
  touchStartY.value = event.touches[0].clientY;
  startLongPress(taskId);
};

const handleTouchEnd = (event: TouchEvent, taskId: string) => {
  cancelLongPress();
  touchEndY.value = event.changedTouches[0].clientY;

  if (isReorderMode.value && reorderTaskId.value === taskId) {
    const deltaY = touchEndY.value - touchStartY.value;

    if (Math.abs(deltaY) > swipeThreshold) {
      if (deltaY > 0) {
        moveTaskDown(taskId);
      } else {
        moveTaskUp(taskId);
      }
    }
  }
};

const handleTouchMove = (event: TouchEvent) => {
  // Cancel long press if user moves finger
  if (Math.abs(event.touches[0].clientY - touchStartY.value) > 10) {
    cancelLongPress();
  }
};

// Desktop mouse event handlers for long press
const handleMouseDown = (event: MouseEvent, taskId: string) => {
  // Only handle left mouse button
  if (event.button !== 0) return;

  startLongPress(taskId);
};

const handleMouseUp = () => {
  cancelLongPress();
};

const handleMouseMove = (event: MouseEvent) => {
  // Cancel long press if user moves mouse significantly
  if (
    longPressTimer.value &&
    (Math.abs(event.movementX) > 5 || Math.abs(event.movementY) > 5)
  ) {
    cancelLongPress();
  }
};

// Reordering tasks (updated for mobile)
const reorderTasks = async () => {
  const reorderedTasks = taskStore.tasks.map((task, index) => ({
    task_id: task.task_id,
    order: index + 1,
  }));
  await taskStore.reorderTask(reorderedTasks);
};
</script>

<template>
  <div class="px-6" :class="isReorderMode ? 'reorder-mode' : ''">
    <!-- Filters + Search -->
    <div
      class="shadow-sm p-3 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md rounded-lg mt-2 border border-surface-light/10 dark:border-surface-dark/10"
    >
      <div class="flex flex-col gap-3">
        <!-- Main Filters -->
        <div class="flex flex-wrap gap-2 items-center">
          <div class="flex gap-1">
            <button
              v-for="filterOption in ['all', 'linked', 'pending', 'complete']"
              :key="filterOption"
              :class="[
                'px-4 py-2 text-sm font-medium transition-all duration-300 rounded-md',
                filter === filterOption
                  ? 'border-b-2 border-primary text-primary dark:text-primary bg-primary/10'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800',
              ]"
              @click="setFilter(filterOption)"
            >
              {{
                filterOption === "linked"
                  ? "Linked"
                  : filterOption.charAt(0).toUpperCase() + filterOption.slice(1)
              }}
            </button>
          </div>
        </div>

        <!-- Priority Filters -->
        <div class="flex flex-wrap gap-2 items-center">
          <span class="text-sm text-text-light dark:text-text-dark font-medium"
            >Priority:</span
          >
          <div class="flex gap-1">
            <button
              v-for="priorityOption in [
                'priority-low',
                'priority-medium',
                'priority-high',
              ]"
              :key="priorityOption"
              :class="[
                'px-4 py-2 text-sm font-medium transition-all duration-300 rounded-md',
                filter === priorityOption
                  ? 'border-b-2 border-primary text-primary dark:text-primary bg-primary/10'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800',
              ]"
              @click="setFilter(priorityOption)"
            >
              {{
                priorityOption === "priority-low"
                  ? "Low"
                  : priorityOption === "priority-medium"
                    ? "Medium"
                    : "High"
              }}
            </button>
          </div>
          <select
            v-model="priorityFilter"
            class="px-2 py-1 text-sm border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark rounded-md focus:outline-none focus:ring-1 focus:ring-primary"
          >
            <option value="all">All</option>
            <option value="1">1 - Low</option>
            <option value="2">2 - Medium</option>
            <option value="3">3 - Medium</option>
            <option value="4">4 - High</option>
            <option value="5">5 - High</option>
          </select>
        </div>

        <!-- Search -->
        <input
          v-model="searchQuery"
          placeholder="Search tasks..."
          class="w-full shadow-sm bg-surface-light dark:bg-surface-dark px-3 py-2 rounded-md border border-surface-light dark:border-surface-dark text-text-light dark:text-text-dark placeholder:text-text-light/50 dark:placeholder:text-text-dark/50 focus:outline-none focus:ring-2 focus:ring-primary"
        />
      </div>
    </div>

    <div
      class="px-3 py-5 bg-surface-light/20 dark:bg-surface-dark/20 backdrop-blur-md shadow-sm mt-4 rounded-lg border border-surface-light/10 dark:border-surface-dark/10"
    >
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-text-light dark:text-text-dark">
          Tasks
        </h3>
      </div>
      <TaskFormTask />
      <p v-if="taskStore.loading" class="text-text-light dark:text-text-dark">
        Loading...
      </p>

      <template v-else>
        <!-- Reorder Mode Header -->
        <div
          v-if="isReorderMode"
          class="mb-4 p-3 bg-primary/10 dark:bg-primary/20 rounded-lg border border-primary/20"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <svg
                class="w-5 h-5 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4 8h16M4 16h16"
                ></path>
              </svg>
              <span class="text-sm font-medium text-primary">Reorder Mode</span>
            </div>
            <UiButton @click="exitReorderMode" variant="default" size="sm">
              Exit
            </UiButton>
          </div>
          <p class="text-xs text-primary/80 mt-1" v-if="isMobile">
            Swipe up/down to move tasks, or tap to exit
          </p>
          <p class="text-xs text-primary/80 mt-1" v-else>
            Drag and drop to reorder tasks, or click to exit
          </p>
        </div>

        <p
          v-if="taskStore.tasks.length === 0"
          class="text-text-light dark:text-text-dark"
        >
          No tasks found
        </p>

        <!-- Desktop: Draggable Component -->
        <draggable
          v-if="!isMobile"
          v-model="taskStore.tasks"
          item-key="task_id"
          @end="reorderTasks"
          :disabled="isFiltering || !isReorderMode"
          class="space-y-3"
          ghost-class="opacity-50"
          chosen-class="shadow-lg"
          drag-class="rotate-2"
        >
          <template #item="{ element: task }">
            <div v-if="matchFilter(task)" :data-task-id="task.task_id">
              <div
                class="bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
                :class="[
                  task.status === 'complete'
                    ? 'border-success'
                    : 'border-warning',
                  isReorderMode && reorderTaskId === task.task_id
                    ? 'ring-2 ring-primary shadow-lg scale-105'
                    : '',
                  isReorderMode ? 'cursor-grab' : 'cursor-pointer',
                ]"
                @mousedown="(e) => handleMouseDown(e, task.task_id)"
                @mouseup="handleMouseUp"
                @mousemove="handleMouseMove"
                @mouseleave="cancelLongPress"
                @click="isReorderMode && !isMobile ? exitReorderMode() : null"
              >
                <!-- Normal view -->
                <div
                  v-if="editingTaskId !== task.task_id"
                  class="flex justify-between items-start"
                >
                  <div @dblclick="startEdit(task)" class="flex-1">
                    <div class="flex items-center gap-2 mb-2">
                      <h3
                        class="text-lg font-semibold text-text-light dark:text-text-dark"
                      >
                        {{ task.title }}
                      </h3>
                    </div>
                    <p
                      class="text-sm text-text-light dark:text-text-dark/80 mb-2"
                    >
                      {{ task.description }}
                    </p>
                    <p
                      class="text-sm text-text-light dark:text-text-dark/60 mb-2"
                    >
                      {{ task.due_date ? dayjs(task.due_date).fromNow() : "" }}
                    </p>
                    <div class="flex items-center gap-2 mt-2">
                      <input
                        type="checkbox"
                        @click="completeTask(task)"
                        :checked="task.status === 'complete'"
                        class="accent-success w-4 h-4"
                      />
                      <span
                        class="text-sm font-medium text-text-light dark:text-text-dark capitalize"
                        >{{ task.status }}</span
                      >
                    </div>
                    <p
                      class="text-sm text-text-light dark:text-text-dark/60 mt-1"
                    >
                      Priority: {{ task.priority }}
                    </p>
                    <div v-if="task.time_estimate_minutes" class="flex items-center gap-1 mt-1">
                      <span class="hidden sm:inline">⏱️</span>
                      <span class="text-sm text-text-light dark:text-text-dark/60">
                        Est: {{ Math.floor(task.time_estimate_minutes / 60) }}h {{ task.time_estimate_minutes % 60 }}m
                      </span>
                    </div>

                    <!-- Subtasks -->
                    <div
                      v-if="task.subtasks && task.subtasks.length > 0"
                      class="mt-3"
                    >
                      <div
                        class="text-xs text-text-light dark:text-text-dark/60 mb-2"
                      >
                        Subtasks ({{
                          task.subtasks.filter((st) => st.status === "complete")
                            .length
                        }}/{{ task.subtasks.length }})
                      </div>
                      <div class="space-y-1">
                        <div
                          v-for="subtask in task.subtasks"
                          :key="subtask.task_id"
                          class="flex items-center text-xs"
                        >
                          <input
                            type="checkbox"
                            @click="completeTask(subtask)"
                            :checked="subtask.status === 'complete'"
                            class="h-3 w-3 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark mr-2 flex-shrink-0"
                          />
                          <span
                            class="text-text-light dark:text-text-dark/80 truncate"
                            :class="
                              subtask.status === 'complete'
                                ? 'line-through opacity-75'
                                : ''
                            "
                          >
                            {{ subtask.title }}
                          </span>
                        </div>
                      </div>
                    </div>

                    <!-- Reorder Mode Controls for Desktop -->
                    <div
                      v-if="
                        isReorderMode &&
                        reorderTaskId === task.task_id &&
                        !isMobile
                      "
                      class="flex items-center justify-between mt-3 p-2 bg-primary/5 rounded-lg"
                    >
                      <UiButton
                        @click="moveTaskUp(task.task_id)"
                        variant="default"
                        size="sm"
                        class="flex items-center gap-1"
                      >
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 15l7-7 7 7"
                          ></path>
                        </svg>
                        Move Up
                      </UiButton>
                      <UiButton
                        @click="moveTaskDown(task.task_id)"
                        variant="default"
                        size="sm"
                        class="flex items-center gap-1"
                      >
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 9l-7 7-7-7"
                          ></path>
                        </svg>
                        Move Down
                      </UiButton>
                    </div>
                  </div>
                  <!-- Three-dot menu button -->
                  <div class="relative ml-2 task-menu-container">
                    <button
                      @click="toggleMenu(task.task_id)"
                      class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors duration-200"
                      title="Task options"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                        ></path>
                      </svg>
                    </button>

                    <!-- Dropdown menu -->
                    <div
                      v-if="openMenuTaskId === task.task_id"
                      class="absolute right-0 mt-2 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-lg border border-surface-light/20 dark:border-surface-dark/20 z-10"
                    >
                      <div class="py-1">
                        <button
                          @click="startEdit(task)"
                          class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                        >
                          <svg
                            class="w-4 h-4 mr-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            ></path>
                          </svg>
                          Edit
                        </button>
                        <button
                          @click="moveTaskUp(task.task_id)"
                          class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                        >
                          <svg
                            class="w-4 h-4 mr-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M5 15l7-7 7 7"
                            ></path>
                          </svg>
                          Move Up
                        </button>
                        <button
                          @click="moveTaskDown(task.task_id)"
                          class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                        >
                          <svg
                            class="w-4 h-4 mr-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M19 9l-7 7-7-7"
                            ></path>
                          </svg>
                          Move Down
                        </button>
                        <div
                          class="border-t border-surface-light/20 dark:border-surface-dark/20"
                        ></div>
                        <button
                          @click="deleteTask(task.task_id)"
                          class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-200"
                        >
                          <svg
                            class="w-4 h-4 mr-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            ></path>
                          </svg>
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Edit mode -->
                <div v-else class="flex flex-col w-full space-y-3">
                  <input
                    v-model="editFormData.title"
                    placeholder="Task title"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <input
                    v-model="editFormData.description"
                    placeholder="Task description"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <input
                    type="datetime-local"
                    v-model="editFormData.due_date"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <input
                    type="number"
                    min="1"
                    max="5"
                    v-model="editFormData.priority"
                    placeholder="Priority (1-5)"
                    class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                  />
                  <div class="flex gap-2">
                    <UiButton
                      @click="saveEdit(task.task_id)"
                      variant="primary"
                      size="sm"
                    >
                      Save
                    </UiButton>
                    <UiButton @click="cancelEdit" variant="default" size="sm">
                      Cancel
                    </UiButton>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </draggable>

        <!-- Mobile: Gesture-Based Reordering -->
        <div v-else class="space-y-3">
          <div
            v-for="task in filteredTasks"
            :key="task.task_id"
            :data-task-id="task.task_id"
          >
            <div
              class="task-card bg-surface-light dark:bg-surface-dark shadow-md rounded-lg p-4 border-l-4 hover:shadow-lg transition-all duration-200"
              :class="[
                task.status === 'complete'
                  ? 'border-success'
                  : 'border-warning',
                isReorderMode && reorderTaskId === task.task_id
                  ? 'ring-2 ring-primary shadow-lg scale-105 reorder-active'
                  : '',
                isReorderMode
                  ? 'cursor-grab active:cursor-grabbing'
                  : 'cursor-pointer',
              ]"
              @touchstart="(e) => handleTouchStart(e, task.task_id)"
              @touchend="(e) => handleTouchEnd(e, task.task_id)"
              @touchmove="handleTouchMove"
              @click="isReorderMode ? exitReorderMode() : null"
            >
              <!-- Reorder Mode Controls -->
              <div
                v-if="isReorderMode && reorderTaskId === task.task_id"
                class="flex items-center justify-between mb-3 p-2 bg-primary/5 rounded-lg"
              >
                <UiButton
                  @click="moveTaskUp(task.task_id)"
                  variant="default"
                  size="sm"
                  class="flex items-center gap-1"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M5 15l7-7 7 7"
                    ></path>
                  </svg>
                  Move Up
                </UiButton>
                <UiButton
                  @click="moveTaskDown(task.task_id)"
                  variant="default"
                  size="sm"
                  class="flex items-center gap-1"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    ></path>
                  </svg>
                  Move Down
                </UiButton>
              </div>

              <!-- Normal view -->
              <div
                v-if="editingTaskId !== task.task_id"
                class="flex justify-between items-start"
              >
                <div @dblclick="startEdit(task)" class="flex-1">
                  <div class="flex items-center gap-2 mb-2">
                    <h3
                      class="text-lg font-semibold text-text-light dark:text-text-dark"
                    >
                      {{ task.title }}
                    </h3>
                  </div>
                  <p
                    class="text-sm text-text-light dark:text-text-dark/80 mb-2"
                  >
                    {{ task.description }}
                  </p>
                  <p
                    class="text-sm text-text-light dark:text-text-dark/60 mb-2"
                  >
                    {{ task.due_date ? dayjs(task.due_date).fromNow() : "" }}
                  </p>
                  <div class="flex items-center gap-2 mt-2">
                    <input
                      type="checkbox"
                      @click="completeTask(task)"
                      :checked="task.status === 'complete'"
                      class="accent-success w-4 h-4"
                    />
                    <span
                      class="text-sm font-medium text-text-light dark:text-text-dark capitalize"
                      >{{ task.status }}</span
                    >
                  </div>
                   <p
                     class="text-sm text-text-light dark:text-text-dark/60 mt-1"
                   >
                     Priority: {{ task.priority }}
                   </p>
                   <div v-if="task.time_estimate_minutes" class="flex items-center gap-1 mt-1">
                     <span class="hidden sm:inline">⏱️</span>
                     <span class="text-sm text-text-light dark:text-text-dark/60">
                       Est: {{ Math.floor(task.time_estimate_minutes / 60) }}h {{ task.time_estimate_minutes % 60 }}m
                     </span>
                   </div>

                   <!-- Subtasks -->
                  <div
                    v-if="task.subtasks && task.subtasks.length > 0"
                    class="mt-3"
                  >
                    <div
                      class="text-xs text-text-light dark:text-text-dark/60 mb-2"
                    >
                      Subtasks ({{
                        task.subtasks.filter((st) => st.status === "complete")
                          .length
                      }}/{{ task.subtasks.length }})
                    </div>
                    <div class="space-y-1">
                      <div
                        v-for="subtask in task.subtasks"
                        :key="subtask.task_id"
                        class="flex items-center text-xs"
                      >
                        <input
                          type="checkbox"
                          @click="completeTask(subtask)"
                          :checked="subtask.status === 'complete'"
                          class="h-3 w-3 text-success rounded focus:ring-success border-surface-light dark:border-surface-dark mr-2 flex-shrink-0"
                        />
                        <span
                          class="text-text-light dark:text-text-dark/80 truncate"
                          :class="
                            subtask.status === 'complete'
                              ? 'line-through opacity-75'
                              : ''
                          "
                        >
                          {{ subtask.title }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
                <!-- Three-dot menu button -->
                <div class="relative ml-2 task-menu-container">
                  <button
                    @click="toggleMenu(task.task_id)"
                    class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors duration-200"
                    title="Task options"
                  >
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
                      ></path>
                    </svg>
                  </button>

                  <!-- Dropdown menu -->
                  <div
                    v-if="openMenuTaskId === task.task_id"
                    class="absolute right-0 mt-2 w-48 bg-surface-light dark:bg-surface-dark rounded-md shadow-lg border border-surface-light/20 dark:border-surface-dark/20 z-10"
                  >
                    <div class="py-1">
                      <button
                        @click="startEdit(task)"
                        class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                      >
                        <svg
                          class="w-4 h-4 mr-3"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          ></path>
                        </svg>
                        Edit
                      </button>
                      <button
                        @click="moveTaskUp(task.task_id)"
                        class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                      >
                        <svg
                          class="w-4 h-4 mr-3"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 15l7-7 7 7"
                          ></path>
                        </svg>
                        Move Up
                      </button>
                      <button
                        @click="moveTaskDown(task.task_id)"
                        class="flex items-center w-full px-4 py-2 text-sm text-text-light dark:text-text-dark hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
                      >
                        <svg
                          class="w-4 h-4 mr-3"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 9l-7 7-7-7"
                          ></path>
                        </svg>
                        Move Down
                      </button>
                      <div
                        class="border-t border-surface-light/20 dark:border-surface-dark/20"
                      ></div>
                      <button
                        @click="deleteTask(task.task_id)"
                        class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-200"
                      >
                        <svg
                          class="w-4 h-4 mr-3"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          ></path>
                        </svg>
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Edit mode -->
              <div v-else class="flex flex-col w-full space-y-3">
                <input
                  v-model="editFormData.title"
                  placeholder="Task title"
                  class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                />
                <input
                  v-model="editFormData.description"
                  placeholder="Task description"
                  class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                />
                <input
                  type="datetime-local"
                  v-model="editFormData.due_date"
                  class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                />
                <input
                  type="number"
                  min="1"
                  max="5"
                  v-model="editFormData.priority"
                  placeholder="Priority (1-5)"
                  class="border border-surface-light dark:border-surface-dark bg-surface-light dark:bg-surface-dark text-text-light dark:text-text-dark px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-primary"
                />
                <div class="flex gap-2">
                  <UiButton
                    @click="saveEdit(task.task_id)"
                    variant="primary"
                    size="sm"
                  >
                    Save
                  </UiButton>
                  <UiButton @click="cancelEdit" variant="default" size="sm">
                    Cancel
                  </UiButton>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
