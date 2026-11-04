import { ref, readonly, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaskShortcuts } from './useTaskShortcuts'

export interface ShortcutAction {
  key: string
  ctrl?: boolean
  alt?: boolean
  shift?: boolean
  meta?: boolean
  action: () => void
  description: string
  category: string
}

// Shared module-level state so the layout, command palette and shortcut
// help modal all read and write the same registry.
const shortcuts = ref<ShortcutAction[]>([])
const commandPaletteOpen = ref(false)
const showHelp = ref(false)

let initialized = false

export function useKeyboardShortcuts() {
  const router = useRouter()
  const task = useTaskShortcuts()

  // Define all shortcuts
  const defineShortcuts = () => {
    shortcuts.value = [
      // Navigation shortcuts
      {
        key: '1',
        ctrl: true,
        action: () => router.push('/dashboard'),
        description: 'Go to Dashboard',
        category: 'Navigation'
      },
      {
        key: '2',
        ctrl: true,
        action: () => router.push('/plan'),
        description: 'Go to Plan',
        category: 'Navigation'
      },
      {
        key: '3',
        ctrl: true,
        action: () => router.push('/time'),
        description: 'Go to Time',
        category: 'Navigation'
      },
      {
        key: '4',
        ctrl: true,
        action: () => router.push('/learning'),
        description: 'Go to Learning',
        category: 'Navigation'
      },
      {
        key: '5',
        ctrl: true,
        action: () => router.push('/finance'),
        description: 'Go to Finance',
        category: 'Navigation'
      },
      {
        key: '6',
        ctrl: true,
        action: () => router.push('/settings'),
        description: 'Go to Settings',
        category: 'Navigation'
      },

      // Quick navigation with letters
      {
        key: 'g',
        ctrl: true,
        action: () => commandPaletteOpen.value = true,
        description: 'Open Command Palette',
        category: 'Navigation'
      },
      {
        key: 'd',
        ctrl: true,
        action: () => router.push('/dashboard'),
        description: 'Go to Dashboard',
        category: 'Navigation'
      },
      {
        key: 'p',
        ctrl: true,
        action: () => router.push('/plan'),
        description: 'Go to Plan',
        category: 'Navigation'
      },
      {
        key: 't',
        ctrl: true,
        action: () => router.push('/time'),
        description: 'Go to Time',
        category: 'Navigation'
      },
      {
        key: 'l',
        ctrl: true,
        action: () => router.push('/learning'),
        description: 'Go to Learning',
        category: 'Navigation'
      },
      {
        key: 'f',
        ctrl: true,
        action: () => router.push('/finance'),
        description: 'Go to Finance',
        category: 'Navigation'
      },
      {
        key: 's',
        ctrl: true,
        action: () => router.push('/settings'),
        description: 'Go to Settings',
        category: 'Navigation'
      },

      // Task management
      {
        key: 'j',
        action: () => task.moveSelection(1),
        description: 'Select next task',
        category: 'Tasks'
      },
      {
        key: 'k',
        action: () => task.moveSelection(-1),
        description: 'Select previous task',
        category: 'Tasks'
      },
      {
        key: 'ArrowDown',
        action: () => task.moveSelection(1),
        description: 'Select next task',
        category: 'Tasks'
      },
      {
        key: 'ArrowUp',
        action: () => task.moveSelection(-1),
        description: 'Select previous task',
        category: 'Tasks'
      },
      {
        key: ' ',
        action: () => task.toggleComplete(),
        description: 'Complete selected task',
        category: 'Tasks'
      },
      {
        key: 'e',
        action: () => task.editSelected(),
        description: 'Edit selected task',
        category: 'Tasks'
      },
      {
        key: 'n',
        action: () => task.newTask(),
        description: 'New task',
        category: 'Tasks'
      },
      {
        key: 'n',
        ctrl: true,
        action: () => task.newTask(),
        description: 'New task',
        category: 'Tasks'
      },
      {
        key: 'd',
        action: () => task.deleteTask(),
        description: 'Delete selected task (undo with u)',
        category: 'Tasks'
      },
      {
        key: 'u',
        action: () => task.undoDelete(),
        description: 'Undo last deleted task',
        category: 'Tasks'
      },

      // General actions
      {
        key: 'a',
        ctrl: true,
        action: () => {
          // This will be handled by individual components
          const event = new CustomEvent('shortcut:new-item')
          window.dispatchEvent(event)
        },
        description: 'Create New Item',
        category: 'Actions'
      },
      {
        key: 'k',
        ctrl: true,
        action: () => commandPaletteOpen.value = true,
        description: 'Open Command Palette',
        category: 'Actions'
      },
      {
        key: '?',
        shift: true,
        action: () => showHelp.value = !showHelp.value,
        description: 'Show keyboard shortcuts',
        category: 'Actions'
      },
      {
        key: 'Escape',
        action: () => {
          commandPaletteOpen.value = false
          showHelp.value = false
          const event = new CustomEvent('shortcut:escape')
          window.dispatchEvent(event)
        },
        description: 'Close modals / Cancel',
        category: 'Actions'
      }
    ]
  }

  const isTypingTarget = (target: EventTarget | null) => {
    return target instanceof HTMLInputElement ||
      target instanceof HTMLTextAreaElement ||
      target instanceof HTMLSelectElement ||
      (target as HTMLElement)?.contentEditable === 'true'
  }

  const handleKeydown = (event: KeyboardEvent) => {
    // Don't trigger shortcuts when typing in input fields
    if (isTypingTarget(event.target)) {
      return
    }

    // While the command palette or help modal is open, only Escape passes through
    if ((commandPaletteOpen.value || showHelp.value) && event.key !== 'Escape') {
      return
    }

    // Don't hijack Space/Enter for activation of focused interactive elements
    if (event.key === ' ') {
      const target = event.target as HTMLElement | null
      if (target?.closest?.('button, a, [role="button"], input, select')) {
        return
      }
    }

    const shortcut = shortcuts.value.find(s =>
      s.key.toLowerCase() === event.key.toLowerCase() &&
      !!s.ctrl === event.ctrlKey &&
      !!s.alt === event.altKey &&
      !!s.shift === event.shiftKey &&
      !!s.meta === event.metaKey
    )

    if (shortcut) {
      event.preventDefault()
      shortcut.action()
    }
  }

  // Define shortcuts immediately so they are available to callers,
  // and attach the global listener only once for the app lifetime.
  if (!initialized) {
    initialized = true
    defineShortcuts()
    if (import.meta.client) {
      onMounted(() => {
        window.addEventListener('keydown', handleKeydown)
      })
    }
  }

  const openHelp = () => {
    showHelp.value = true
  }

  const closeHelp = () => {
    showHelp.value = false
  }

  return {
    shortcuts: readonly(shortcuts),
    commandPaletteOpen: readonly(commandPaletteOpen),
    showHelp,
    addShortcut,
    removeShortcut,
    openCommandPalette: () => commandPaletteOpen.value = true,
    closeCommandPalette: () => commandPaletteOpen.value = false,
    openHelp,
    closeHelp,
  }
}

const addShortcut = (shortcut: ShortcutAction) => {
  shortcuts.value.push(shortcut)
}

const removeShortcut = (key: string, ctrl?: boolean, alt?: boolean, shift?: boolean, meta?: boolean) => {
  const index = shortcuts.value.findIndex(s =>
    s.key === key &&
    !!s.ctrl === !!ctrl &&
    !!s.alt === !!alt &&
    !!s.shift === !!shift &&
    !!s.meta === !!meta
  )
  if (index > -1) {
    shortcuts.value.splice(index, 1)
  }
}
