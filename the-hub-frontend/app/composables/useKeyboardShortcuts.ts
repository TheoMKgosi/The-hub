import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

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

export function useKeyboardShortcuts() {
  const router = useRouter()
  const shortcuts = ref<ShortcutAction[]>([])
  const commandPaletteOpen = ref(false)

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

      // General actions
      {
        key: 'n',
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
        key: 'Escape',
        action: () => {
          commandPaletteOpen.value = false
          const event = new CustomEvent('shortcut:escape')
          window.dispatchEvent(event)
        },
        description: 'Close modals / Cancel',
        category: 'Actions'
      }
    ]
  }

  const handleKeydown = (event: KeyboardEvent) => {
    // Don't trigger shortcuts when typing in input fields
    if (event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement ||
        event.target instanceof HTMLSelectElement ||
        (event.target as HTMLElement)?.contentEditable === 'true') {
      return
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

  onMounted(() => {
    defineShortcuts()
    window.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
  })

  return {
    shortcuts: readonly(shortcuts),
    commandPaletteOpen: readonly(commandPaletteOpen),
    addShortcut,
    removeShortcut,
    openCommandPalette: () => commandPaletteOpen.value = true,
    closeCommandPalette: () => commandPaletteOpen.value = false
  }
}