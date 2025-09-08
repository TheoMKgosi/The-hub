<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'

interface Command {
  id: string
  title: string
  description: string
  category: string
  action: () => void
  shortcut?: string
}

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const searchQuery = ref('')
const selectedIndex = ref(0)
const searchInput = ref<HTMLInputElement | null>(null)

// Define all available commands
const allCommands = ref<Command[]>([
  // Navigation commands
  {
    id: 'nav-dashboard',
    title: 'Go to Dashboard',
    description: 'Navigate to the main dashboard',
    category: 'Navigation',
    action: () => router.push('/dashboard'),
    shortcut: 'Ctrl+1 or Ctrl+D'
  },
  {
    id: 'nav-plan',
    title: 'Go to Plan',
    description: 'Navigate to task planning',
    category: 'Navigation',
    action: () => router.push('/plan'),
    shortcut: 'Ctrl+2 or Ctrl+P'
  },
  {
    id: 'nav-time',
    title: 'Go to Time',
    description: 'Navigate to time tracking',
    category: 'Navigation',
    action: () => router.push('/time'),
    shortcut: 'Ctrl+3 or Ctrl+T'
  },
  {
    id: 'nav-learning',
    title: 'Go to Learning',
    description: 'Navigate to learning section',
    category: 'Navigation',
    action: () => router.push('/learning'),
    shortcut: 'Ctrl+4 or Ctrl+L'
  },
  {
    id: 'nav-finance',
    title: 'Go to Finance',
    description: 'Navigate to financial management',
    category: 'Navigation',
    action: () => router.push('/finance'),
    shortcut: 'Ctrl+5 or Ctrl+F'
  },
  {
    id: 'nav-settings',
    title: 'Go to Settings',
    description: 'Navigate to settings',
    category: 'Navigation',
    action: () => router.push('/settings'),
    shortcut: 'Ctrl+6 or Ctrl+S'
  },

  // Action commands
  {
    id: 'action-new-task',
    title: 'Create New Task',
    description: 'Create a new task',
    category: 'Actions',
    action: () => {
      const event = new CustomEvent('command:new-task')
      window.dispatchEvent(event)
    },
    shortcut: 'Ctrl+N'
  },
  {
    id: 'action-new-goal',
    title: 'Create New Goal',
    description: 'Create a new goal',
    category: 'Actions',
    action: () => {
      const event = new CustomEvent('command:new-goal')
      window.dispatchEvent(event)
    },
    shortcut: 'Ctrl+N'
  },
  {
    id: 'action-new-budget',
    title: 'Create New Budget',
    description: 'Create a new budget',
    category: 'Actions',
    action: () => {
      const event = new CustomEvent('command:new-budget')
      window.dispatchEvent(event)
    },
    shortcut: 'Ctrl+N'
  }
])

// Filter commands based on search query
const filteredCommands = computed(() => {
  if (!searchQuery.value) return allCommands.value

  const query = searchQuery.value.toLowerCase()
  return allCommands.value.filter(command =>
    command.title.toLowerCase().includes(query) ||
    command.description.toLowerCase().includes(query) ||
    command.category.toLowerCase().includes(query)
  )
})

// Group commands by category
const groupedCommands = computed(() => {
  const groups: Record<string, Command[]> = {}
  filteredCommands.value.forEach(command => {
    if (!groups[command.category]) {
      groups[command.category] = []
    }
    groups[command.category].push(command)
  })
  return groups
})

// Handle keyboard navigation
const handleKeydown = (event: KeyboardEvent) => {
  if (!props.isOpen) return

  switch (event.key) {
    case 'ArrowDown':
      event.preventDefault()
      selectedIndex.value = Math.min(selectedIndex.value + 1, filteredCommands.value.length - 1)
      break
    case 'ArrowUp':
      event.preventDefault()
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
      break
    case 'Enter':
      event.preventDefault()
      if (filteredCommands.value[selectedIndex.value]) {
        executeCommand(filteredCommands.value[selectedIndex.value])
      }
      break
    case 'Escape':
      event.preventDefault()
      closePalette()
      break
  }
}

const executeCommand = (command: Command) => {
  command.action()
  closePalette()
}

const closePalette = () => {
  searchQuery.value = ''
  selectedIndex.value = 0
  emit('close')
}

const focusSearchInput = async () => {
  await nextTick()
  searchInput.value?.focus()
}

// Watch for open state changes
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    focusSearchInput()
  }
})

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-start justify-center pt-16 bg-black/50 backdrop-blur-sm"
        @click.self="closePalette"
      >
        <div class="w-full max-w-2xl mx-4">
          <!-- Search Input -->
          <div class="bg-surface-light dark:bg-surface-dark rounded-lg shadow-2xl border border-surface-light dark:border-surface-dark overflow-hidden">
            <div class="p-4 border-b border-surface-light dark:border-surface-dark">
              <div class="flex items-center gap-3">
                <div class="flex-shrink-0 w-5 h-5 text-text-light dark:text-text-dark/60">
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <input
                  ref="searchInput"
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search commands..."
                  class="flex-1 bg-transparent border-0 outline-none text-text-light dark:text-text-dark placeholder-text-light/60 dark:placeholder-text-dark/60"
                />
                <div class="flex-shrink-0 text-xs text-text-light/60 dark:text-text-dark/60">
                  <kbd class="px-1.5 py-0.5 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs">Esc</kbd>
                </div>
              </div>
            </div>

            <!-- Commands List -->
            <div class="max-h-96 overflow-y-auto">
              <div v-if="filteredCommands.length === 0" class="p-4 text-center text-text-light/60 dark:text-text-dark/60">
                No commands found
              </div>

              <div v-else>
                <div
                  v-for="(commands, category) in groupedCommands"
                  :key="category"
                  class="border-b border-surface-light/20 dark:border-surface-dark/20 last:border-b-0"
                >
                  <div class="px-4 py-2 text-xs font-medium text-text-light/60 dark:text-text-dark/60 uppercase tracking-wide">
                    {{ category }}
                  </div>

                  <div
                    v-for="(command, index) in commands"
                    :key="command.id"
                    :class="[
                      'px-4 py-3 cursor-pointer transition-colors flex items-center justify-between',
                      selectedIndex === filteredCommands.indexOf(command) ? 'bg-primary/10 dark:bg-primary/20' : 'hover:bg-surface-light/50 dark:hover:bg-surface-dark/50'
                    ]"
                    @click="executeCommand(command)"
                  >
                    <div class="flex-1">
                      <div class="font-medium text-text-light dark:text-text-dark">
                        {{ command.title }}
                      </div>
                      <div class="text-sm text-text-light/60 dark:text-text-dark/60">
                        {{ command.description }}
                      </div>
                    </div>

                    <div v-if="command.shortcut" class="flex-shrink-0 ml-4">
                      <kbd class="px-2 py-1 bg-surface-light/50 dark:bg-surface-dark/50 rounded text-xs text-text-light/60 dark:text-text-dark/60">
                        {{ command.shortcut }}
                      </kbd>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Footer -->
            <div class="px-4 py-3 border-t border-surface-light/20 dark:border-surface-dark/20 bg-surface-light/30 dark:bg-surface-dark/30">
              <div class="flex items-center justify-between text-xs text-text-light/60 dark:text-text-dark/60">
                <div class="flex items-center gap-4">
                  <span>↑↓ to navigate</span>
                  <span>↵ to select</span>
                </div>
                <div>
                  {{ filteredCommands.length }} commands
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>